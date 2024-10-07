package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ircurry/dfh/internal/ipc"
)

func main() {
	lsnCmd := flag.NewFlagSet("spawn", flag.ExitOnError)
	// TODO: Document Negative numbers being the same as 0
	lsnNum := lsnCmd.Int("n", 0, "a zero or positive number for how many events to listen to and print")

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "spwn":
			ipc.HyprMessage("dispatch exec " + os.Args[2])
			return
		case "lsn":
			ipc.HyprPrintEvents(*lsnNum)
			return
		}
	}

	fmt.Fprintln(os.Stderr, "Error: no valid sub command given as first argument")
	return
}
