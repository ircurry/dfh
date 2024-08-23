package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ircurry/dfh/internal/cli"
	"github.com/ircurry/dfh/internal/ipc"
	"github.com/ircurry/dfh/internal/monitors"
)

func main() {
	monsCmd := flag.NewFlagSet("monitors", flag.ExitOnError)
	monsFile := monsCmd.String("f", "", "the json file to read monitor definitions from")
	monsAllowUnconnected := monsCmd.Bool("a", false, "allow unconnected monitors to be configured")


	if len(os.Args) <= 1 {
		cli.Die("Not enough arguments", cli.ArgumentError)
	}
	monsCmd.Parse(os.Args[2:])
	state := monsCmd.Arg(0)
	if state == "" {
		cli.Die("No state given", cli.MonitorStateFailure)
	}

	contents, err := monitors.ReadMonitorConfigFile(*monsFile)
	cli.DieIfErr("Unable to read file.", err, cli.ReadFileFailure)

	var monl monitors.MonitorList
	err = monl.FromJson(contents)
	cli.DieIfErr("Something went wrong parsing config file",
		err, cli.MonitorConfigParseFailure)
	stateStrings, err := ipc.StateStrings(monl, state)
	cli.DieIfErr("Error creating hyprland monitor settings", err, cli.MonitorStateFailure)

	wlrdata, err := ipc.WlrRandrJson()
	cli.DieIfErr("Something went wrong requesting monitor information", err, cli.CommandExecutionError)
	monitorNames, err := ipc.WlrRandrGetMonitors(wlrdata)
	cli.DieIfErr("Could not get monitor names from program", err, cli.InfoRetrevalFailure)
	allMonsPresent, err := monitors.CompareMonitorLists(monl, monitorNames)
	cli.DieIfErr("Monitor name not found", err, cli.MonitorConfigParseFailure)

	if !(allMonsPresent || *monsAllowUnconnected) {
		monsCmd.Usage()
		cli.Die("Not all monitors in config are present. Consider using the flag -a",
			cli.InfoRetrevalFailure)
	}

	for _, str := range stateStrings {
		fmt.Println(str)
		// TODO: make this work with just IPC
		if output, err := ipc.HyprctlExecCommand("keyword", "monitor", str); err != nil {
			cli.DieErr("something went wrong executing hyprctl", err, cli.CommandExecutionError)
		} else {
			fmt.Print(string(output))
		}
	}
	return
}
