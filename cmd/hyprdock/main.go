package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ircurry/dfh/internal/cli"
	"github.com/ircurry/dfh/internal/ipc"
	"github.com/ircurry/dfh/internal/monitors"
)

type cliParseError struct {
	eMessage string
}

func (e cliParseError) Error() string {
	return e.eMessage
}

func newCliParseError(msg string) cliParseError {
	return cliParseError{msg}
}

const helpText = `hyprdock - cli tool to easily configure tools
Usage:
  [options...] [monitors...]

Options:
  -h, --help                    display help information
  -e, --enabled-monitors        print the names of monitors to be enabled in profile
  -d, --disabled-monitors       print the names of monitors to be disabled in profile
  -a, --all-monitors            print the names of all monitors specified in profile
`

type cliArgs struct {
	profile          string
	help             bool
	enabledMonitors  bool
	disabledMonitors bool
	allMonitors      bool
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
		case "-h", "--help":
			c.help = true
			return nil
		case "-e", "--enabled-monitors":
			c.enabledMonitors = true
		case "-d", "--disabled-monitors":
			c.disabledMonitors = true
		case "-a", "--all-monitors":
			c.allMonitors = true
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

	if progArgs.help {
		fmt.Print(helpText)
		os.Exit(0)
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
profileLoop:
	for _, prfl := range monitorConfig {
		if prfl.Name != progArgs.profile {
			continue
		}

		prflFound = true
		switch {
		case progArgs.enabledMonitors:
			enabledMons := make([]string, 0)
			for _, mon := range prfl.Monitors {
				if mon.Enabled && mon.Name != nil {
					enabledMons = append(enabledMons, *mon.Name)
				}
			}
			for i := 0; i < len(enabledMons); i++ {
				fmt.Printf("%s\n", enabledMons[i])
			}
			break profileLoop
		case progArgs.disabledMonitors:
			disabledMons := make([]string, 0)
			for _, mon := range prfl.Monitors {
				if !mon.Enabled && mon.Name != nil {
					disabledMons = append(disabledMons, *mon.Name)
				}
			}
			for i := 0; i < len(disabledMons); i++ {
				fmt.Printf("%s\n", disabledMons[i])
			}
			break profileLoop
		case progArgs.allMonitors:
			for i := 0; i < len(prfl.Monitors); i++ {
				monName := ""
				if prfl.Monitors[i].Name != nil {
					monName = fmt.Sprintf("%s\n", *prfl.Monitors[i].Name)
				}
				fmt.Print(monName)
			}
			break profileLoop
		}

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
