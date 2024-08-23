package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ircurry/dfh/internal/ipc"
	"github.com/ircurry/dfh/internal/monitors"
)

const (
	_ = iota
	ReadFileFailure
	MonitorConfigParseFailure
	MonitorStateFailure
	CommandExecutionError
	InfoRetrevalFailure
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

func main() {
	lsnCmd := flag.NewFlagSet("spawn", flag.ExitOnError)
	// TODO: Document Negative numbers being the same as 0
	lsnNum := lsnCmd.Int("n", 0, "a zero or positive number for how many events to listen to and print")
	monsCmd := flag.NewFlagSet("monitors", flag.ExitOnError)
	monsFile := monsCmd.String("f", "", "the json file to read monitor definitions from")
	monsAllowUnconnected := monsCmd.Bool("a", false, "allow unconnected monitors to be configured")

	if len(os.Args) > 1 {
		switch os.Args[1] {
		// TODO: make an IPC monitor parser to tell what monitors are attached
		case "mons":
			// hyprCtlSock := getHyprCtlSocket()
			monsCmd.Parse(os.Args[2:])
			state := monsCmd.Arg(0)
			if state == "" {
				Die("No state given", MonitorStateFailure)
			}
			contents, err := monitors.ReadMonitorConfigFile(*monsFile)
			DieIfErr("Unable to read file.", err, ReadFileFailure)
			var monl monitors.MonitorList
			err = monl.FromJson(contents)
			DieIfErr("Something went wrong parsing config file",
				err, MonitorConfigParseFailure)
			stateStrings, err := ipc.StateStrings(monl, state)
			DieIfErr("Error creating hyprland monitor settings", err, MonitorStateFailure)

			wlrdata, err := ipc.WlrRandrJson()
			DieIfErr("Something went wrong requesting monitor information", err, CommandExecutionError)
			monitorNames, err := ipc.WlrRandrGetMonitors(wlrdata)
			DieIfErr("Could not get monitor names from program", err, InfoRetrevalFailure)
			allMonsPresent, err := monitors.CompareMonitorLists(monl, monitorNames)
			DieIfErr("Monitor name not found", err, MonitorConfigParseFailure)

			if !(allMonsPresent || *monsAllowUnconnected) {
				monsCmd.Usage()
				Die("Not all monitors in config are present. Consider using the flag -a",
					InfoRetrevalFailure)
			}

			for _, str := range stateStrings {
				fmt.Println(str)
				// TODO: make this work with just IPC
				if output, err := ipc.HyprctlExecCommand("keyword", "monitor", str); err != nil {
					DieErr("something went wrong executing hyprctl", err, CommandExecutionError)
				} else {
					fmt.Print(string(output))
				}
			}
			return
		case "spwn":
			ipc.HyprMessage("dispatch exec " + os.Args[2])
			return
		case "lsn":
			ipc.HyprPrintEvents(*lsnNum)
			return
		case "test":
			data, err := ipc.WlrRandrJson()
			DieIfErr("Something Went Wrong", err, 25)
			names, err := ipc.WlrRandrGetMonitors(data)
			DieIfErr("Something Went Wrong", err, 26)
			for _, name := range names {
				fmt.Println(name)
			}

			return
		}
	}

	fmt.Fprintln(os.Stderr, "Error: no valid sub command given as first argument")
	return
}
