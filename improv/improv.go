package improv

var (
	SERVICE_UUID      = [16]byte{0x00, 0x46, 0x77, 0x68, 0x62, 0x28, 0x22, 0x72, 0x46, 0x63, 0x27, 0x74, 0x78, 0x26, 0x80, 0x00}
	STATUS_UUID       = [16]byte{0x00, 0x46, 0x77, 0x68, 0x62, 0x28, 0x22, 0x72, 0x46, 0x63, 0x27, 0x74, 0x78, 0x26, 0x80, 0x01}
	ERROR_UUID        = [16]byte{0x00, 0x46, 0x77, 0x68, 0x62, 0x28, 0x22, 0x72, 0x46, 0x63, 0x27, 0x74, 0x78, 0x26, 0x80, 0x02}
	RPC_COMMAND_UUID  = [16]byte{0x00, 0x46, 0x77, 0x68, 0x62, 0x28, 0x22, 0x72, 0x46, 0x63, 0x27, 0x74, 0x78, 0x26, 0x80, 0x03}
	RPC_RESULT_UUID   = [16]byte{0x00, 0x46, 0x77, 0x68, 0x62, 0x28, 0x22, 0x72, 0x46, 0x63, 0x27, 0x74, 0x78, 0x26, 0x80, 0x04}
	CAPABILITIES_UUID = [16]byte{0x00, 0x46, 0x77, 0x68, 0x62, 0x28, 0x22, 0x72, 0x46, 0x63, 0x27, 0x74, 0x78, 0x26, 0x80, 0x05}
)

type improvState byte

const (
	STATE_STOPPED improvState = iota
	STATE_AWAITING_AUTHORIZATION
	STATE_AUTHORIZED
	STATE_PROVISIONING
	STATE_PROVISIONED
)

type improvCommand byte

const (
	COMMAND_UNKNOWN improvCommand = iota
	COMMAND_WIFI_SETTINGS
	COMMAND_IDENTIFY
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

type improvError byte

const (
	ERROR_NONE improvError = iota
	ERROR_INVALID_RPC
	ERROR_UNKNOWN_RPC
	ERROR_UNABLE_TO_CONNECT
	ERROR_NOT_AUTHORIZED
	ERROR_UNKNOWN
)

type Improv struct {
	state improvState
}

func (i *Improv) ParseImprovData(data []byte) (improvCommand, []string) {
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

func (i *Improv) BuildImprovResponse(cmd improvCommand, args []string) []byte {
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

func New() *Improv {
	return &Improv{
		state: STATE_AUTHORIZED,
	}
}
