package main

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

// sleepWithContext sleeps for the specified duration or until the context is done
func sleepWithContext(ctx context.Context, d time.Duration) {
	timer := time.NewTimer(d)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
	case <-timer.C:
	}
}

// executeCommand executes a command and returns the output as a string
func executeCommand(command string, args ...string) (string, error) {
	// split the command into the program and its arguments
	parts := strings.Split(command, " ")
	args = append(parts[1:], args...)
	debugln("Running command:", parts[0], strings.Join(args, " "))
	cmd := exec.Command(parts[0], args...)

	output, err := cmd.CombinedOutput()
	return string(output), err
}

// quote quotes a string with double quotes
func quote(str string) string {
	return "\"" + str + "\""
}

func infoln(args ...string) {
	println("[INFO]", strings.Join(args, " "))
}

func debugln(args ...string) {
	println("[DEBUG]", strings.Join(args, " "))
}

func errorln(args ...string) {
	println("[ERROR]", strings.Join(args, " "))
}
