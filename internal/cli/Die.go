package cli

import (
	"fmt"
	"os"
)

const (
	_ = iota
	ReadFileFailure
	MonitorConfigParseFailure
	MonitorStateFailure
	CommandExecutionError
	InfoRetrevalFailure
	ArgumentError
)

func Die(message string, exitCode int) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(exitCode)
}

func DieErr(message string, err error, exitCode int) {
	fmt.Fprintln(os.Stderr, message)
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(exitCode)
}

func DieIfErr(message string, err error, exitCode int) {
	if err != nil {
		DieErr(message, err, exitCode)
	}
	return
}
