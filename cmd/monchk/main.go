package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ircurry/dfh/internal/cli"
	"github.com/ircurry/dfh/internal/ipc"
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

const helpText = `monchk - check if monitors are connected
Usage:
  [options...] [monitors...]

Options:
  -h, --help        display help information
      --            read names from stdin, with each line being the name
`

type cliArgs struct {
	monitors []string
	help     bool
	stdin    bool
}

func newCliArgs() cliArgs {
	return cliArgs{
		monitors: make([]string, 0),
	}
}

func (c *cliArgs) parseArgs(args []string) error {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-h", "--help":
			c.help = true
		case "--":
			c.stdin = true
		default:
			c.monitors = append(c.monitors, args[i])
		}
	}

	if c.stdin {
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		mons := strings.Split(string(bytes), "\n")
		for _, v := range mons {
			if v != "" {
				c.monitors = append(c.monitors, v)
			}
		}
	}

	return nil
}

func main() {
	progArgs := newCliArgs()
	err := progArgs.parseArgs(os.Args[1:])
	if err != nil {
		cli.Die(err.Error(), cli.ArgumentError)
	}

	if progArgs.help {
		fmt.Print(helpText)
		os.Exit(0)
	} else if (len(os.Args) < 2) && !progArgs.stdin {
		fmt.Print(helpText)
		os.Exit(cli.CommandParseFailure)
	}

	output, err := ipc.HyprctlExecCommand("monitors", "-j")
	if err != nil {
		errMsg := fmt.Sprintf("There was an \033[1;31merror\033[0m getting the list attached monitors.\n%s", err.Error())
		cli.Die(errMsg, cli.CommandExecutionError)
	}
	monitors := make([]map[string]interface{}, 0)
	err = json.Unmarshal(output, &monitors)
	if err != nil {
		errMsg := "\033[31mFailed to parse JSON file.\033[0m\n" + err.Error()
		cli.Die(errMsg, cli.MonitorConfigParseFailure)
	}
	if len(monitors) == 0 {
		cli.Die("Failed to get monitor information or information is blank.", cli.MonitorStateFailure)
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

	for i, val := range foundMonitor {
		if !val {
			msg := fmt.Sprintf("Could not find monitor '\033[1;33m%s\033[0m'.\n", progArgs.monitors[i])
			fmt.Fprint(os.Stderr, msg)
			os.Exit(127)
		} else {
			msg := fmt.Sprintf("Found monitor '\033[1;33m%s\033[0m'.\n", progArgs.monitors[i])
			fmt.Fprint(os.Stderr, msg)
		}
	}

}
