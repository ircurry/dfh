package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
)

const (
	_ = iota
	ReadFileFailure
	MonitorConfigParseFailure
	MonitorStateFailure
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
	val, err := json.Marshal(mon)
	if err != nil {
		panic(err)
	}
	return string(val)
}

func (mon *Monitor) EnableString() string {
	hyprstring := fmt.Sprintf("%s,%dx%d,%dx%d@%d,%d",
		mon.Name, mon.Width, mon.Height, mon.X, mon.Y, mon.RefreshRate, mon.Scale)
	return hyprstring
}

func (mon *Monitor) DisableString() string {
	hyprstring := fmt.Sprintf("%s,disable", mon.Name)
	return hyprstring
}

type MonitorList []Monitor

func (monlist *MonitorList) fromJson(byt []byte) error {
	return json.Unmarshal(byt, monlist);
}

func (monlist *MonitorList) stateStrings(state string) ([]string, error) {
	strlist := make([]string, 0)
	enablelist := make([]string, 0)
	disablelist := make([]string, 0)
	var stateErr error = nil
	containsState := false
	for _, mon := range *monlist {
		if mon.State == state {
			enablelist = append(enablelist, mon.EnableString())
			containsState = true
		} else {
			disablelist = append(disablelist, mon.DisableString())
		}
	}

	for _, str := range enablelist {
		strlist = append(strlist, str)
	}
	for _, str := range disablelist {
		strlist = append(strlist, str)
	}

	if !containsState {
		stateErr = fmt.Errorf("this monitor configuration does not contain state: %s", state)
	}

	return strlist, stateErr
}

func readMonitorConfigFile(file string) ([]byte, error) {
	switch file {
	case "":
		configDir, err := os.UserConfigDir()
		if err != nil {
			panic(err)
		}
		fileContents, errFile := os.ReadFile(configDir + "/dfh/monitors.json")
		return fileContents, errFile
	default:
		fileContents, err := os.ReadFile(file)
		return fileContents, err
	}
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
	monsCmd := flag.NewFlagSet("monitors", flag.ExitOnError)
	monsNum := monsCmd.String("f", "", "the json file to read monitor definitions from")

	if (len(os.Args) > 1) {
		switch os.Args[1] {
		case "mons":
			// hyprCtlSock := getHyprCtlSocket()
			monsCmd.Parse(os.Args[2:])
			state := monsCmd.Arg(0)
			if state == "" {
				fmt.Fprintln(os.Stderr, "No state given")
				os.Exit(MonitorStateFailure)
			}
			contents, err := readMonitorConfigFile(*monsNum)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Unable to read file.")
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(ReadFileFailure)
			}
			var monl MonitorList
			err = monl.fromJson(contents)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Something went wrong parsing config file")
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(MonitorConfigParseFailure)
			}
			stateStrings, err := monl.stateStrings(state)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error creating hyprland monitor settings")
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(MonitorStateFailure)
			}
			for _, str := range stateStrings {
				// TODO: make this work with just IPC
				cmd := exec.Command("hyprctl", "keyword", "monitor", str)
				if output, err := cmd.Output(); err != nil {
					panic(err)
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
