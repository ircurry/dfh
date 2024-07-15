package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

type monitor struct {
	name string
	width int64
	height int64
	refreshRate int
	x int64
	y int64
	scale int
	state string
};

func (mon *monitor) String() string {
	hyprstring := fmt.Sprintf("%s,%dx%d,%dx%d@%d,%d",
		mon.name, mon.height, mon.width, mon.x, mon.y, mon.refreshRate, mon.scale)
	return hyprstring
}

func (mon *monitor) DisableString() string {
	hyprstring := fmt.Sprintf("%s,disable", mon.name)
	return hyprstring
}


func listenEvents(con *net.Conn, lim int, ch chan string) {
    b := make([]byte, 1)
    str := ""

    for i := 0; (i < lim) || (lim <= 0); i++{
		for {
			(*con).Read(b)
			str = str + string(b[:])
			if string(b[:]) == "\n" {
				break
			}
		}
		ch <- str
		str = ""
    }

    close(ch)
    return
}

func printEvents(ch chan string) {
    for {
		msg, ok := <-ch
		if !ok {
			break
		}
		fmt.Print(msg)
    }
    return
}

func hyprMessage(con net.Conn, msg string) {
    _, err := con.Write([]byte(msg + "\n"))
    if err != nil {
		panic(err)
    }
}

func getHyprDir() string {
    var xdgRunDir, his string
    var ok bool
    xdgRunDir, ok = os.LookupEnv("XDG_RUNTIME_DIR")
    if !ok {
		panic("XDG_RUNTIME_DIR not found")
    }
    
    his, ok = os.LookupEnv("HYPRLAND_INSTANCE_SIGNATURE")
    if !ok {
		panic("HYPRLAND_INSTANCE_SIGNATURE not found")
    }
    
    return xdgRunDir + "/hypr/" + his
}

func getHyprCtlSocket() (net.Conn) {
    hyprDir := getHyprDir()
    sock1, err := net.Dial("unix", hyprDir + "/.socket.sock")
    if err != nil {
		panic(err)
    }
    return sock1
}

func getEventSocket() (net.Conn) {
    hyprDir := getHyprDir()
    sock2, err := net.Dial("unix", hyprDir + "/.socket2.sock")
    if err != nil {
		panic(err)
    }
    return sock2
}

func main() {
	lsnCmd := flag.NewFlagSet("spawn", flag.ExitOnError)
	lsnNum := lsnCmd.Int("n", 0, "a zero or positive number for how many events to listen to and print")

	if (len(os.Args) > 1) {
		switch os.Args[1] {
		case "dir":
			fmt.Printf("Hyprland dir: %s\n", getHyprDir())
			return
		case "spwn":
			hyprCtlSock := getHyprCtlSocket()
			hyprMessage(hyprCtlSock, "dispatch exec " + os.Args[2])
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
