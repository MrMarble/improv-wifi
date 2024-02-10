package improv_test

import (
	"testing"

	"github.com/mrmarble/improv-wifi"
)

func TestParseCommand(t *testing.T) {
	t.Parallel()
	t.Run("WifiSettings", func(t *testing.T) {
		d := []byte{(improv.COMMAND_WIFI_SETTINGS), 0x00}
		ssid := "mySsid"
		pass := "secretPass"
		d = append(d, byte(len(ssid)))
		d = append(d, []byte(ssid)...)
		d = append(d, byte(len(pass)))
		d = append(d, []byte(pass)...)
		d[1] = byte(len(d) - 2)

		cmd, args := improv.ParseImprovData(d)
		if cmd != improv.COMMAND_WIFI_SETTINGS {
			t.Errorf("Expected COMMAND_WIFI_SETTINGS, got %v", cmd)
		}
		if len(args) != 2 {
			t.Errorf("Expected 2 arguments, got %v", len(args))
		}
		if args[0] != ssid {
			t.Errorf("Expected %s, got %v", ssid, args[0])
		}
		if args[1] != pass {
			t.Errorf("Expected %s, got %v", pass, args[1])
		}
	})

	t.Run("Identify", func(t *testing.T) {
		d := []byte{(improv.COMMAND_IDENTIFY)}
		cmd, args := improv.ParseImprovData(d)
		if cmd != improv.COMMAND_IDENTIFY {
			t.Errorf("Expected COMMAND_IDENTIFY, got %v", cmd)
		}
		if args != nil {
			t.Errorf("Expected nil, got %v", args)
		}
	})
}

func TestBuildResponse(t *testing.T) {
	t.Parallel()
	t.Run("WifiSettings", func(t *testing.T) {
		url := "http://0.0.0.0:8080/setup"
		expected := []byte{0x00, (improv.COMMAND_WIFI_SETTINGS), byte(len(url))}
		expected = append(expected, []byte(url)...)
		expected[0] = byte(len(expected) - 2)

		got := improv.BuildImprovResponse(improv.COMMAND_WIFI_SETTINGS, []string{url})
		if len(got) != len(expected) {
			t.Errorf("Expected %v, got %v", expected, got)
		}
		for i := range got {
			if got[i] != expected[i] {
				t.Errorf("Expected %v, got %v", expected, got)
			}
		}
	})
}
