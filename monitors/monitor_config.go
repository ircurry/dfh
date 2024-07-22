package monitors

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
)

type Monitor struct {
	Name        *string `json:"name"`
	Width       *int64  `json:"width"`
	Height      *int64  `json:"height"`
	RefreshRate *int    `json:"refreshRate"`
	X           *int64  `json:"x"`
	Y           *int64  `json:"y"`
	Scale       *int    `json:"scale"`
	State       *string `json:"state"`
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
		if val.Name == nil {
			return false, fmt.Errorf("monitor has Name nil")
		}
		if !slices.Contains(names, *val.Name) {
			cont = false
		}
	}
	return cont, nil
}
