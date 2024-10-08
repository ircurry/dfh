package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ircurry/dfh/internal/cli"
	"github.com/ircurry/dfh/internal/ipc"
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
	monitors []string
}

func newCliArgs() cliArgs {
	return cliArgs{
		monitors: make([]string, 0),
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
				eeMessage: "monchk [-h | --help] [-m <monitor>]\n",
			}
		case "-m":
			c.monitors = append(c.monitors ,args[i+1])
			i++
		default:
			extraArgs = append(extraArgs, args[i])
		}
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
		switch err.(type) {
		case earlyExit:
			fmt.Print(err.Error())
		default:
			cli.Die(err.Error(), cli.ArgumentError)
		}
	}
	output, err := ipc.HyprctlExecCommand("monitors", "-j")
	if err != nil {
		errMsg := fmt.Sprintf("There was an \033[1;31merror\033[0m getting the list attached monitors.\n%s", err.Error())
		cli.Die(errMsg, cli.CommandExecutionError)
	}

	monitors := make([]map[string]interface{}, 0)
	json.Unmarshal(output, &monitors)
	if len(monitors) == 0 {
		cli.Die("Failed to get monitor information.", cli.MonitorStateFailure)
	}

	foundMonitor := make([]bool, len(progArgs.monitors))
	for i, mon := range progArgs.monitors {
		for _, monJson := range monitors {
			monName, ok := monJson["name"]
			if !ok {
				msg := fmt.Sprintf("There was a monitor without a name.")
				fmt.Fprint(os.Stderr, msg)
				continue
			}
			if monName == mon {
				foundMonitor[i] = true
			}
		}
	}

	for _, val := range foundMonitor {
		if !val {
			os.Exit(127)
		}
	}
	
}
