package main

import "os"

func main() {
	err := run(os.Args[1:])
	if err != nil {
		println("[ERROR]", err.Error())
		os.Exit(1)
	}
}
