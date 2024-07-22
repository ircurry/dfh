package ipc

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ircurry/dfh/monitors"
)

func WlrRandrExecCommand(args ...string) ([]byte, error) {
	cmd := exec.Command("wlr-randr", args...)
	data, err := cmd.Output()
	return data, err
}

func monOptions(mon monitors.Monitor, state string) ([]string) {
	wlrRandrStrs := []string{ "--output", mon.Name }
	if mon.State == state {
		wlrRandrStrs = append(wlrRandrStrs, "--on")
		wlrRandrStrs = append(wlrRandrStrs, "--mode")
		wlrRandrStrs = append(wlrRandrStrs, fmt.Sprintf("%dx%d@%d",
			mon.Width, mon.Height, mon.RefreshRate))
		wlrRandrStrs = append(wlrRandrStrs, "--pos", fmt.Sprintf("%d,%d", mon.X, mon.Y))
		wlrRandrStrs = append(wlrRandrStrs, "--scale", fmt.Sprintf("%d", mon.Scale))

	} else {
		wlrRandrStrs = append(wlrRandrStrs, "--off")
	}
	return wlrRandrStrs
}

func WlrRandrJson() (output []byte, err error) {
	cmd := exec.Command("wlr-randr", "--json")
	output, err = cmd.Output()
	return
}

func WlrRandrGetMonitors(data []byte) ([]string, error) {
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
