package ipc

import (
	"fmt"
	"net"
	"os"
	"os/exec"

	"github.com/ircurry/dfh/monitors"
)

func hyprMonString(mon monitors.Monitor, state string) (string, bool, error) {
	if err := mon.CheckStringFields(); err != nil {
		return "", false, err
	}
	if err := mon.CheckStringFields(); err != nil {
		return "", false, err
	}
	if *mon.State == state {
		var name, resolution, position, scale string

		name = *mon.Name
		resolution = fmt.Sprintf("%dx%d@%d", *mon.Width, *mon.Height, *mon.RefreshRate)
		position = fmt.Sprintf("%dx%d", *mon.X, *mon.Y)
		scale = fmt.Sprintf("%d", *mon.Scale)

		hyprstring := fmt.Sprintf("%s,%s,%s,%s", name, resolution, position, scale)
		return hyprstring, true, nil
	} else {
		return fmt.Sprintf("%s,disable", *mon.Name), false, nil
	}

}

func StateStrings(monlist monitors.MonitorList, state string) ([]string, error) {
	strlist := make([]string, 0)
	var stateErr error = nil
	containsState := false
	for _, mon := range monlist {
		if str, isState, err := hyprMonString(mon, state); err != nil {
			return nil, err
		} else {
			strlist = append(strlist, str)
			if isState {
				containsState = true
			}
		}
	}

	if !containsState {
		stateErr = fmt.Errorf("this monitor configuration does not contain state: %s", state)
	}

	return strlist, stateErr
}

func listenEvents(con net.Conn, lim int, ch chan string) {
	b := make([]byte, 1)
	str := ""

	for i := 0; (i < lim) || (lim <= 0); i++ {
		for {
			con.Read(b)
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

func HyprPrintEvents(num int) {
	ch := make(chan string)
	go listenEvents(getEventSocket(), num, ch)
	printEvents(ch)
	return
}

func sendMessage(con net.Conn, msg string) {
	_, err := con.Write([]byte(msg + "\n"))
	if err != nil {
		panic(err)
	}
	return
}

func HyprMessage(msg string) {
	sendMessage(getHyprCtlSocket(), msg)
	return
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

func getHyprCtlSocket() net.Conn {
	hyprDir := getHyprDir()
	sock1, err := net.Dial("unix", hyprDir+"/.socket.sock")
	if err != nil {
		panic(err)
	}
	return sock1
}

func getEventSocket() net.Conn {
	hyprDir := getHyprDir()
	sock2, err := net.Dial("unix", hyprDir+"/.socket2.sock")
	if err != nil {
		panic(err)
	}
	return sock2
}

func HyprctlExecCommand(args ...string) (output []byte, err error) {
	cmd := exec.Command("hyprctl", args...)
	output, err = cmd.Output()
	return
}
