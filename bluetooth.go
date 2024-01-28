package main

import (
	"context"
	"errors"

	"github.com/mrmarble/improv/improv"
	"tinygo.org/x/bluetooth"
)

func startAdvertising(opts *Options) error {
	adapter := bluetooth.DefaultAdapter
	err := adapter.Enable()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	err = setCapabilities(adapter, opts, cancel)
	if err != nil {
		return errors.New("Failed to set capabilities: " + err.Error())
	}
	adv, err := configureAdvertisement(adapter, opts)
	if err != nil {
		return err
	}

	addr, err := adapter.Address()
	if err != nil {
		return errors.New("Failed to get address: " + err.Error())
	}

	err = adv.Start()
	if err != nil {
		return errors.New("Failed to start advertisement: " + err.Error())
	}

	if opts.timeout > 0 {
		infoln("Advertising", opts.name, "for", opts.timeout.String(), "on address:", addr.String())
		sleepWithContext(ctx, opts.timeout)
	} else {
		infoln("Advertising", opts.name, "until interrupted on address:", addr.String())
		select {}
	}
	return nil
}

func configureAdvertisement(adapter *bluetooth.Adapter, opts *Options) (*bluetooth.Advertisement, error) {
	adv := adapter.DefaultAdvertisement()
	defer adv.Stop()

	err := adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: opts.name,
		ServiceUUIDs: []bluetooth.UUID{
			bluetooth.NewUUID(improv.SERVICE_UUID),
		},
	})
	if err != nil {
		return nil, errors.New("Failed to configure advertisement: " + err.Error())
	}
	return adv, nil
}

func setCapabilities(adapter *bluetooth.Adapter, opts *Options, cancel context.CancelFunc) error {
	protocol := improv.New()

	var statusHandler bluetooth.Characteristic
	var rpcResultHandler bluetooth.Characteristic
	var errorHandler bluetooth.Characteristic

	identfyCapability := []byte{0x0}
	if opts.identifyCommand != "" {
		identfyCapability = []byte{0x1}
	}
	return adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.NewUUID(improv.SERVICE_UUID),
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &statusHandler,
				UUID:   bluetooth.NewUUID(improv.STATUS_UUID),
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicNotifyPermission,
				Value:  []byte{byte(improv.STATE_AUTHORIZED)},
			},
			{
				Handle: &errorHandler,
				UUID:   bluetooth.NewUUID(improv.ERROR_UUID),
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicNotifyPermission,
				Value:  []byte{byte(improv.ERROR_NONE)},
			},
			{
				UUID:  bluetooth.NewUUID(improv.RPC_COMMAND_UUID),
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicWriteWithoutResponsePermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					cmd, args := protocol.ParseImprovData(value)
					switch cmd {
					case improv.COMMAND_WIFI_SETTINGS:
						configureWIFI(args, errorHandler, statusHandler, opts, rpcResultHandler, protocol, cancel)
					case improv.COMMAND_IDENTIFY:
						identifyDevice(opts, errorHandler)
					default:
						infoln("Unknown command:", cmd.String())
						errorHandler.Write([]byte{byte(improv.ERROR_UNKNOWN_RPC)})
					}
				},
			},
			{
				Handle: &rpcResultHandler,
				UUID:   bluetooth.NewUUID(improv.RPC_RESULT_UUID),
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				UUID:  bluetooth.NewUUID(improv.CAPABILITIES_UUID),
				Flags: bluetooth.CharacteristicReadPermission,
				Value: identfyCapability,
			},
		},
	})
}

func identifyDevice(opts *Options, errorHandler bluetooth.Characteristic) {
	debugln("Running identify command:", opts.identifyCommand)
	output, err := executeCommand(opts.identifyCommand)
	if err != nil {
		errorln("Running identify command:", err.Error())
		errorHandler.Write([]byte{byte(improv.ERROR_UNKNOWN)})
	}
	if output != "" {
		debugln("Identify command output:", output)
	}
}

func configureWIFI(args []string, errorHandler bluetooth.Characteristic, statusHandler bluetooth.Characteristic, opts *Options, rpcResultHandler bluetooth.Characteristic, protocol *improv.Improv, cancel context.CancelFunc) {
	if args == nil || len(args) != 2 {
		errorln("Invalid wifi settings command")
		errorHandler.Write([]byte{byte(improv.ERROR_INVALID_RPC)})
	}
	infoln("Provisioning wifi")
	statusHandler.Write([]byte{byte(improv.STATE_PROVISIONING)})
	if opts.wifiCommand != "" {
		debugln("Running wifi command")
		output, err := executeCommand(opts.wifiCommand, args...)
		if err != nil {
			errorln("Got error", quote(err.Error()), "running wifi command, command output:", output)
			errorHandler.Write([]byte{byte(improv.ERROR_UNABLE_TO_CONNECT)})
			statusHandler.Write([]byte{byte(improv.STATE_AUTHORIZED)})
			errorHandler.Write([]byte{byte(improv.ERROR_NONE)})
			return
		}
		if output != "" {
			infoln("Wifi command output:", output)
			url := findUrl(output)
			if url != "" {
				rpcResultHandler.Write(protocol.BuildImprovResponse(improv.COMMAND_WIFI_SETTINGS, []string{url}))
			}
		}
	} else {
		debugln("Running nmcli")
		output, err := executeCommand("nmcli", "device", "wifi", "rescan")
		if err != nil {
			errorln("Got error", quote(err.Error()), "running nmcli, command output:", output)
			errorHandler.Write([]byte{byte(improv.ERROR_UNABLE_TO_CONNECT)})
			statusHandler.Write([]byte{byte(improv.STATE_AUTHORIZED)})
			errorHandler.Write([]byte{byte(improv.ERROR_NONE)})
			return
		}
		executeCommand("nmcli", "device", "wifi", "delete", args[0])
		output, err = executeCommand("nmcli", "device", "wifi", "connect", quote(args[0]), "password", quote(args[1]))
		if err != nil {
			errorln("Got error", quote(err.Error()), "running nmcli, command output:", output)
			errorHandler.Write([]byte{byte(improv.ERROR_UNABLE_TO_CONNECT)})
			statusHandler.Write([]byte{byte(improv.STATE_AUTHORIZED)})
			errorHandler.Write([]byte{byte(improv.ERROR_NONE)})
			return
		}
	}
	statusHandler.Write([]byte{byte(improv.STATE_PROVISIONED)})
	infoln("Wifi provisioned")
	cancel()
}
