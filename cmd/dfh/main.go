package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

const (
	_ = iota
	ReadFileFailure
	MonitorConfigParseFailure
	MonitorStateFailure
	CommandExecutionError
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

func main() {
	lsnCmd := flag.NewFlagSet("spawn", flag.ExitOnError)
	// TODO: Document Negative numbers being the same as 0
	lsnNum := lsnCmd.Int("n", 0, "a zero or positive number for how many events to listen to and print")
	monsCmd := flag.NewFlagSet("monitors", flag.ExitOnError)
	monsNum := monsCmd.String("f", "", "the json file to read monitor definitions from")

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
			contents, err := readMonitorConfigFile(*monsNum)
			if err != nil {
				DieErr("Unable to read file.", err, ReadFileFailure)
			}
			var monl MonitorList
			err = monl.fromJson(contents)
			if err != nil {
				DieErr(
					"Something went wrong parsing config file",
					err, MonitorConfigParseFailure)
			}
			stateStrings, err := monl.stateStrings(state)
			if err != nil {
				DieErr("Error creating hyprland monitor settings", err, MonitorStateFailure)
			}
			for _, str := range stateStrings {
				// TODO: make this work with just IPC
				cmd := exec.Command("hyprctl", "keyword", "monitor", str)
				if output, err := cmd.Output(); err != nil {
					DieErr("something went wrong executing hyprctl", err, CommandExecutionError)
				} else {
					fmt.Print(string(output))
				}
			}
			return
		case "dir":
			fmt.Printf("Hyprland dir: %s\n", getHyprDir())
			return
		case "spwn":
			hyprCtlSock := getHyprCtlSocket()
			hyprMessage(hyprCtlSock, "dispatch exec "+os.Args[2])
			return
		case "lsn":
			eventSock := getEventSocket()
			lsnCmd.Parse(os.Args[2:])
			ch := make(chan string)
			go listenEvents(&eventSock, *lsnNum, ch)
			printEvents(ch)
			return
		}
	}

	fmt.Fprintln(os.Stderr, "Error: no valid sub command given as first argument")
	return
}
