package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ircurry/dfh/internal/cli"
	"github.com/ircurry/dfh/internal/ipc"
)

var (
	topLevelUsage = `Usage: dfh [options] <subcommand> [subcommand option...]
Subcommands:
  spwn          Spawn a program through hyprctl keyword exec
  lsn           Listen to hyprland events
`
	spwnUsage = "Usage: dfh spwn <command>\n"
	lsnUsage  = "Usage: dfh lsn [options]\n"
)

func main() {
	exitCode := 0
	defer func() { os.Exit(exitCode) }()
	spwnCmd := flag.NewFlagSet("spwn", flag.ExitOnError)
	spwnCmd.Usage = func() {
		fmt.Print(spwnUsage)
		fmt.Fprint(os.Stderr, cli.Usage(spwnCmd))
	}
	lsnCmd := flag.NewFlagSet("lsn", flag.ExitOnError)
	lsnNum := lsnCmd.Uint("n", 0,
		`a zero or positive number for how many events to listen to and print
zero means output events until interupted which is the default behavior`)
	lsnCmd.Usage = func() {
		fmt.Print(lsnUsage)
		fmt.Fprint(os.Stderr, cli.Usage(lsnCmd))
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "spwn":
			spwnCmd.Parse(os.Args[2:])
			argsLen := len(spwnCmd.Args())
			switch argsLen {
			case 1:
				ipc.HyprMessage("dispatch exec " + spwnCmd.Arg(0))
			case 0:
				fmt.Fprintf(os.Stderr, "no command supplied\n")
				spwnCmd.Usage()
				exitCode = 1
			default:
				fmt.Fprintf(os.Stderr, "more than one command supplied, only the first will be used\n")
				ipc.HyprMessage("dispatch exec " + spwnCmd.Arg(0))
			}
			return
		case "lsn":
			lsnCmd.Parse(os.Args[2:])
			ipc.HyprPrintEvents(int(*lsnNum))
			return
		case "-h", "--help":
			fmt.Print(topLevelUsage)
			return
		}
	}

	fmt.Fprintln(os.Stderr, "Error: no valid sub command given as first argument")
	fmt.Print(topLevelUsage)
	exitCode = 1
	return
}
