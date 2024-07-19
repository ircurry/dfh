package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

func listenEvents(con *net.Conn, lim int, ch chan string) {
	b := make([]byte, 1)
	str := ""

	for i := 0; (i < lim) || (lim <= 0); i++ {
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

func runHyprctl(args ...string) (output []byte, err error) {
	cmd := exec.Command("hyprctl", args...)
	output, err = cmd.Output()
	return
}

func wlrRandrJson() (output []byte, err error) {
	cmd := exec.Command("wlr-randr", "--json")
	output, err = cmd.Output()
	return
}

func wlrRandrGetMonitors(data []byte) ([]string, error) {
	dec := json.NewDecoder(strings.NewReader(string(data)))
	isArray := false
	tkn, err := dec.Token()
	if err != nil {
		return nil, err
	}
	delim, ok := tkn.(json.Delim)
	if !ok {
		return nil, fmt.Errorf("Not a Set or Array")
	}
	switch delim {
	case '[':
		isArray = true
	case '{':
		isArray = false
	default:
		return nil, fmt.Errorf("Unexpected JSON delimiter")
	}

	if isArray {
		var dat []map[string]interface{}
		err := json.Unmarshal(data, &dat)
		if err != nil {
			return nil, err
		}

		strs := make([]string, 0)
		for _, val := range dat {
			name, ok := val["name"].(string)
			if !ok {
				return nil, fmt.Errorf("unable to get the name of one or more monitors\n%s",
					string(data))
			}
			strs = append(strs, name)
		}
		return strs, nil
	} else {
		var dat map[string]interface{}
		err := json.Unmarshal(data, &dat)
		if err != nil {
			return nil, err
		}

		name, ok := dat["name"].(string)
		if !ok {
			return nil, fmt.Errorf("unable to get the name of one or more monitors\n%s",
				string(data))
		}
		str := make([]string, 0)
		return append(str, name), nil
	}
}
