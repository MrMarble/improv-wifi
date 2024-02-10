// Package improv implements the [Improv] protocol for IoT devices.
//
// The Improv protocol is a simple protocol for configuring IoT devices over Bluetooth Low Energy (BLE).
// Allowing for the configuration of Wi-Fi settings and device identification over BLE.
//
// [Improv]: https://www.improv-wifi.com/
package improv

// UUIDs for the Improv service.
var (
	SERVICE_UUID      = [16]byte{0x00, 0x46, 0x77, 0x68, 0x62, 0x28, 0x22, 0x72, 0x46, 0x63, 0x27, 0x74, 0x78, 0x26, 0x80, 0x00} // Advertisement UUID
	STATUS_UUID       = [16]byte{0x00, 0x46, 0x77, 0x68, 0x62, 0x28, 0x22, 0x72, 0x46, 0x63, 0x27, 0x74, 0x78, 0x26, 0x80, 0x01} // Current state of the device
	ERROR_UUID        = [16]byte{0x00, 0x46, 0x77, 0x68, 0x62, 0x28, 0x22, 0x72, 0x46, 0x63, 0x27, 0x74, 0x78, 0x26, 0x80, 0x02} // Current error state of the device
	RPC_COMMAND_UUID  = [16]byte{0x00, 0x46, 0x77, 0x68, 0x62, 0x28, 0x22, 0x72, 0x46, 0x63, 0x27, 0x74, 0x78, 0x26, 0x80, 0x03} // Command received from the client
	RPC_RESULT_UUID   = [16]byte{0x00, 0x46, 0x77, 0x68, 0x62, 0x28, 0x22, 0x72, 0x46, 0x63, 0x27, 0x74, 0x78, 0x26, 0x80, 0x04} // Result of the command received from the client
	CAPABILITIES_UUID = [16]byte{0x00, 0x46, 0x77, 0x68, 0x62, 0x28, 0x22, 0x72, 0x46, 0x63, 0x27, 0x74, 0x78, 0x26, 0x80, 0x05} // Capabilities of the device (Identify device)

)

type (
	improvState   byte
	improvCommand byte
	improvError   byte
)

// State constants
const (
	STATE_STOPPED                improvState = iota
	STATE_AWAITING_AUTHORIZATION             // Awaiting authorization via physical interaction.
	STATE_AUTHORIZED                         // Ready to accept credentials.
	STATE_PROVISIONING                       // Credentials received, attempt to connect.
	STATE_PROVISIONED                        // Connection successful.

)

// Command constants
const (
	COMMAND_UNKNOWN improvCommand = iota
	COMMAND_WIFI_SETTINGS
	COMMAND_IDENTIFY
)

// Error constants
const (
	ERROR_NONE              improvError = iota // This shows there is no current error state.
	ERROR_INVALID_RPC                          // RPC packet was malformed/invalid.
	ERROR_UNKNOWN_RPC                          // The command sent is unknown.
	ERROR_UNABLE_TO_CONNECT                    // The credentials have been received and an attempt to connect to the network has failed.
	ERROR_NOT_AUTHORIZED                       // Credentials were sent via RPC but the Improv service is not authorized.
	ERROR_UNKNOWN           = 0xFF
)

func (s improvCommand) String() string {
	switch s {
	case COMMAND_UNKNOWN:
		return "COMMAND_UNKNOWN"
	case COMMAND_WIFI_SETTINGS:
		return "COMMAND_WIFI_SETTINGS"
	case COMMAND_IDENTIFY:
		return "COMMAND_IDENTIFY"
	}
	return "UNKNOWN"
}

// ParseImprovData parses the data received from the Improv service and returns the command and arguments.
func ParseImprovData(data []byte) (improvCommand, []string) {
	cmd := improvCommand(data[0])

	switch cmd {
	case COMMAND_WIFI_SETTINGS:
		ssidLength := int(data[2])
		ssidStart := 3
		ssidEnd := ssidStart + ssidLength

		passLength := int(data[ssidEnd])
		passStart := ssidEnd + 1
		passEnd := passStart + passLength

		ssid := string(data[ssidStart:ssidEnd])
		password := string(data[passStart:passEnd])

		return cmd, []string{ssid, password}
	case COMMAND_IDENTIFY:
		return cmd, nil
	}

	return COMMAND_UNKNOWN, nil
}

// BuildImprovResponse builds the Improv response packet.
func BuildImprovResponse(cmd improvCommand, args []string) []byte {
	output := []byte{0x00, byte(cmd)}
	length := 0
	for _, arg := range args {
		len := len(arg)
		length += len + 1
		output = append(output, byte(len))
		output = append(output, []byte(arg)...)
	}
	output[0] = byte(length)
	return output
}
