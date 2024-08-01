package monitors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

type Monitor struct {
	Name        string `json:"name"`
	Width       int64  `json:"width"`
	Height      int64  `json:"height"`
	RefreshRate int    `json:"refreshRate"`
	X           int64  `json:"x"`
	Y           int64  `json:"y"`
	Scale       int    `json:"scale"`
	State       string `json:"state"`
}

func (mon *Monitor) String() string {
	return fmt.Sprintf("Name: %s\nWidth: %d\nHeight: %d\nRefresh: %d\nX: %d\nY: %d\nScale: %d\nState: %s",
		mon.Name, mon.Width, mon.Height, mon.RefreshRate, mon.X, mon.Y, mon.Scale, mon.State)
}

func fmtInvalidTypeErr(str string, valType string, v any) error {
	return fmt.Errorf("key %s expect type %s, found %T", str, valType, v)
}

func (mon *Monitor) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(strings.NewReader(string(data)))
	tkn, err := dec.Token()
	if err == io.EOF {
		return errors.New("Reached EOF before any json tokens could be read")
	} else if err != nil {
		return err
	}

	delim, ok := tkn.(json.Delim)
	if !ok {
		return fmt.Errorf("first token was not a delimeter, instead was %T", tkn)
	} else if delim.String() != "{" {
		return fmt.Errorf("JSON token is not an begin-object, instead was %s", delim.String())
	}

	keysSet := map[string]bool{
		"name": false,
		"width": false,
		"height": false,
		"refreshRate": false,
		"x": false,
		"y": false,
		"scale": false,
		"state": false,
	}

	for dec.More() {
		tkn, err := dec.Token()
		if err != nil {
			return err
		}

		var key string
		switch tkn.(type) {
		case string:
			key = tkn.(string)
		case json.Delim:
			if delim := tkn.(json.Delim); delim.String() == "}" {
				break
			} else {
				return fmt.Errorf("Expected JSON key to close ")
			}
		}

		if !dec.More() {
			return fmt.Errorf("File terminates before key %s can be read", key)
		}

		switch key {
		case "name":
			tkn, err := dec.Token()
			if err != nil {
				return err
			}
			val, ok := tkn.(string)
			if !ok {
				return fmtInvalidTypeErr(key, "string", val)
			}
			mon.Name = val
			keysSet[key] = true
		case "width":
			tkn, err := dec.Token()
			if err != nil {
				return err
			}
			val, ok := tkn.(float64)
			if !ok {
				return fmtInvalidTypeErr(key, "int64", val)
			}
			mon.Width = int64(val)
			keysSet[key] = true
		case "height":
			tkn, err := dec.Token()
			if err != nil {
				return err
			}
			val, ok := tkn.(float64)
			if !ok {
				return fmtInvalidTypeErr(key, "int64", val)
			}
			mon.Height = int64(val)
			keysSet[key] = true
		case "refreshRate":
			tkn, err := dec.Token()
			if err != nil {
				return err
			}
			val, ok := tkn.(float64)
			if !ok {
				return fmtInvalidTypeErr(key, "int", val)
			}
			mon.RefreshRate = int(val)
			keysSet[key] = true
		case "x":
			tkn, err := dec.Token()
			if err != nil {
				return err
			}
			val, ok := tkn.(float64)
			if !ok {
				return fmtInvalidTypeErr(key, "int64", val)
			}
			mon.X = int64(val)
			keysSet[key] = true
		case "y":
			tkn, err := dec.Token()
			if err != nil {
				return err
			}
			val, ok := tkn.(float64)
			if !ok {
				return fmtInvalidTypeErr(key, "int64", val)
			}
			mon.Y = int64(val)
			keysSet[key] = true
		case "scale":
			tkn, err := dec.Token()
			if err != nil {
				return err
			}
			val, ok := tkn.(float64)
			if !ok {
				return fmtInvalidTypeErr(key, "int", val)
			}
			mon.Scale = int(val)
			keysSet[key] = true
		case "state":
			tkn, err := dec.Token()
			if err != nil {
				return err
			}
			val, ok := tkn.(string)
			if !ok {
				return fmtInvalidTypeErr(key, "string", val)
			}
			mon.State = val
			keysSet[key] = true
		default:
			dec.Token()
			continue
		}

	}

	for k, v := range keysSet {
		if !v {
			return fmt.Errorf("Key %s not set", k)
		}
	}
	return nil
}

type MonitorList []Monitor

func (monlist *MonitorList) FromJson(byt []byte) error {
	return json.Unmarshal(byt, monlist)
}

func ReadMonitorConfigFile(file string) ([]byte, error) {
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

func CompareMonitorLists(monl MonitorList, names []string) (bool, error) {
	cont := true
	for _, val := range monl {
		if !slices.Contains(names, val.Name) {
			cont = false
		}
	}
	return cont, nil
}
