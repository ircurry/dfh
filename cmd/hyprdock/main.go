package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/ircurry/dfh/internal/cli"
	"github.com/ircurry/dfh/internal/ipc"
	"github.com/ircurry/dfh/internal/monitors"
)

type cliArgs struct {
	profile          string
	enabledMonitors  bool
	disabledMonitors bool
	allMonitors      bool
}

func parseArgs(args []string) (cliArgs, error) {
	progargs := cliArgs{}
	flags := flag.NewFlagSet("hyprdock", flag.ContinueOnError)
	flags.BoolVar(&progargs.enabledMonitors, "e", false, "print the names of monitors to be enabled in profile")
	flags.BoolVar(&progargs.disabledMonitors, "d", false, "print the names of monitors to be disabled in profile")
	flags.BoolVar(&progargs.allMonitors, "a", false, "print the names of all monitors specified in profile")
	flags.BoolVar(&progargs.enabledMonitors, "enabled-monitors", false, "print the names of monitors to be enabled in profile")
	flags.BoolVar(&progargs.disabledMonitors, "disabled-monitors", false, "print the names of monitors to be disabled in profile")
	flags.BoolVar(&progargs.allMonitors, "all-monitors", false, "print the names of all monitors specified in profile")
	flags.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage: hyprdock [options] <profile>\n")
		flags.PrintDefaults()
	}
	err := flags.Parse(args)
	if err != nil {
		return progargs, err
	}
	switch len((flags.Args())) {
	case 0:
		fmt.Fprint(os.Stderr, "not enough arguments\n")
		flags.Usage()
		return progargs, nil
	case 1:
		progargs.profile = flags.Arg(0)
		return progargs, nil
	default:
		fmt.Fprint(os.Stderr, "more than one profile name supplied, using first name given\n")
		progargs.profile = flags.Arg(0)
		return progargs, nil
	}
}

func main() {
	progArgs, err := parseArgs(os.Args[1:])
	switch (err) {
	case nil:
		break
	case flag.ErrHelp:
		return
	default:
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
