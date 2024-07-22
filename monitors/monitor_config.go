package monitors

import (
	"encoding/json"
	"errors"
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

func (mon *Monitor) CheckStateField() error {
	if mon.State == nil {
		return fmt.Errorf("cannont use monitor cofiguration, does not contain a State")
	} else {
		return nil
	}
}

func (mon *Monitor) CheckStringFields() error {
	switch v := true; v {
	case (mon.Name == nil):
		return errors.New("cannot format string, Name is nil.")
	case (mon.Width == nil):
		return errors.New("cannot format string, Width is nil.")
	case (mon.Height == nil):
		return errors.New("cannot format string, Height is nil.")
	case (mon.RefreshRate == nil):
		return errors.New("cannot format string, Refresh Rate is nil.")
	case (mon.X == nil):
		return errors.New("cannot format string, X Position is nil.")
	case (mon.Y == nil):
		return errors.New("cannot format string, Y Position is nil.")
	case (mon.Scale == nil):
		return errors.New("cannot format string, Scale is nil.")
	default:
		return nil
	}
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
