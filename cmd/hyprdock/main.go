package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ircurry/dfh/internal/cli"
	"github.com/ircurry/dfh/internal/ipc"
	"github.com/ircurry/dfh/internal/monitors"
)

type earlyExit struct {
	eeType    string
	eeMessage string
}

func (e earlyExit) Error() string {
	return e.eeMessage
}

type cliParseError struct {
	eMessage string
}

func (e cliParseError) Error() string {
	return e.eMessage
}

func newCliParseError(msg string) cliParseError {
	return cliParseError{msg}
}

type cliArgs struct {
	profile string
}

func newCliArgs() cliArgs {
	return cliArgs{
		profile: "",
	}
}

func (c *cliArgs) parseArgs(args []string) error {
	extraArgs := make([]string, 0)
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-h":
			fallthrough
		case "--help":
			return earlyExit{
				eeType:    "help",
				eeMessage: "hyprdock <profile>\n",
			}
		default:
			if c.profile == "" {
				c.profile = args[i]
			} else {
				extraArgs = append(extraArgs, args[i])
			}
		}
	}

	if c.profile == "" {
		return newCliParseError("No profile has been specified\n")
	}
	if len(extraArgs) != 0 {
		extraArgsString := ""
		for _, v := range extraArgs {
			extraArgsString += fmt.Sprintf("  Unknown Argument: '%s'\n", v)
		}
		return newCliParseError(fmt.Sprintf("An issue occured while parsing arguments:\n" + extraArgsString))
	}
	return nil
}

func main() {
	progArgs := newCliArgs()
	err := progArgs.parseArgs(os.Args[1:])
	if err != nil {
		cli.Die(err.Error(), cli.CommandParseFailure)
	}

	configDir := ""
	configDir, err = os.UserConfigDir()
	if err != nil {
		configDir, err = os.UserHomeDir()
		if err != nil {
			panic("Could not locate user home dir")
		}
		configDir += "/.config/nocturne"
	}
	configDir += "/nocturne"

	configFileConents, err := os.ReadFile(configDir + "/monitors.json")
	if err != nil {
		errMsg := "Something when wrong reading monitor config file.\n" + err.Error()
		cli.Die(errMsg, cli.ReadFileFailure)
	}

	var monitorConfig []monitors.Profile
	json.Unmarshal(configFileConents, &monitorConfig)

	prflFound := false
	for _, prfl := range monitorConfig {
		if prfl.Name != progArgs.profile {
			continue
		}
		prflFound = true
		fmt.Printf("Configuring monitors according to profile \033[1;32m%s\033[0m.\n", prfl.Name)

		hyprMonitorStrings := ipc.MonitorProfileToHyprlandString(prfl)
		for _, hyprstr := range hyprMonitorStrings {
			fmt.Printf("  Monitor '\033[1;33m%s\033[0m'\n", hyprstr)
			output, err := ipc.HyprctlExecCommand("keyword", "monitor", hyprstr)
			if err != nil {
				errMsg := fmt.Sprintf("Something went wrong configuring monitor '\033[1;33m%s\033[0m'\n%s", hyprstr, string(output))
				cli.Die(errMsg, cli.CommandExecutionError)
			}
		}
		break
	}
	if !prflFound {
		cli.Die(fmt.Sprintf("Could not find profile '\033[1;31m%s\033[0m'", progArgs.profile), cli.ArgumentError)
	}
}
