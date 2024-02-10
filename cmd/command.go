package main

import (
	"flag"
	"os"
	"time"
)

var (
	version = "main"
	commit  = "?"
	date    = ""
)

type Options struct {
	name            string
	identifyCommand string
	wifiCommand     string
	timeout         time.Duration
	debug           bool
}

func run(args []string) error {
	options, err := parseArguments(args)
	if err != nil {
		return err
	}

	return startAdvertising(options)
}

func parseArguments(args []string) (*Options, error) {
	hostame, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	var name string
	var identifyCommand string
	var wifiCommand string
	var timeout string
	var debug bool

	flag.StringVar(&name, "name", hostame, "")
	flag.StringVar(&name, "n", hostame, "")
	flag.StringVar(&identifyCommand, "identify-command", "", "")
	flag.StringVar(&identifyCommand, "i", "", "")
	flag.StringVar(&wifiCommand, "wifi-command", "", "")
	flag.StringVar(&wifiCommand, "w", "", "")
	flag.StringVar(&timeout, "timeout", "2m", "")
	flag.StringVar(&timeout, "t", "2m", "")
	flag.BoolVar(&debug, "debug", false, "")
	flag.BoolVar(&debug, "d", false, "")

	flag.Usage = func() {
		println("improv - A tool for advertising wifi settings over bluetooth")
		println(version + " (" + commit + ") " + date)
		println()
		println("Usage: improv [options]")
		println()
		println("Options:")
		println("  -n, --name <name>	The name of the bluetooth device. (default is the hostname)")
		println("  -i, --identify-command <command>	The command to run when identifying the device")
		println("  -w, --wifi-command <command>	The command to run when setting the wifi settings. (default is to use nmcli)")
		println("  -t, --timeout <duration>	The number of minutes to advertise the device for. (default is 2m. 0 means advertise forever)")
		println("  -d, --debug		Enable debug logging")
		println("  -h, --help		Show this help message")
	}

	flag.Parse()
	parsedTimeout, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, err
	}
	return &Options{
		name:            name,
		identifyCommand: identifyCommand,
		wifiCommand:     wifiCommand,
		timeout:         parsedTimeout,
		debug:           debug,
	}, nil
}
