# Improv

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mrmarble/improv-wifi)
[![Go Report Card](https://goreportcard.com/badge/github.com/mrmarble/improv-wifi)](https://goreportcard.com/report/github.com/mrmarble/improv-wifi)

Simple [improv wifi](https://www.improv-wifi.com/) implementation in Go

```
Usage: improv [options]

Options:
  -n, --name <name>     The name of the bluetooth device. (default is the hostname)
  -i, --identify-command <command>      The command to run when identifying the device
  -w, --wifi-command <command>  The command to run when setting the wifi settings. (default is to use nmcli)
  -t, --timeout <duration>      The number of minutes to advertise the device for. (default is 2m. 0 means advertise forever)
  -d, --debug           Enable debug logging
  -h, --help            Show this help message
```
