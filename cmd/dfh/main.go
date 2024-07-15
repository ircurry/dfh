package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

type Monitor struct {
	Name string `json:"name"`
	Width int64 `json:"width"`
	Height int64 `json:"height"`
	RefreshRate int `json:"refreshRate"`
	X int64 `json:"x"`
	Y int64 `json:"y"`
	Scale int `json:"scale"`
	State string `json:"state"`
};

func (mon *Monitor) String() string {
	hyprstring := fmt.Sprintf("%s,%dx%d,%dx%d@%d,%d",
		mon.Name, mon.Height, mon.Width, mon.X, mon.Y, mon.RefreshRate, mon.Scale)
	return hyprstring
}

func (mon *Monitor) DisableString() string {
	hyprstring := fmt.Sprintf("%s,disable", mon.Name)
	return hyprstring
}


func listenEvents(con *net.Conn, lim int, ch chan string) {
    b := make([]byte, 1)
    str := ""

    for i := 0; (i < lim) || (lim <= 0); i++{
		for {
			(*con).Read(b)
			if string(b[:]) == "\n" {
				break
			}
			str = str + string(b[:])
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
		fmt.Println(msg)
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
	// TODO: Document Negative numbers being the same as 0
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
