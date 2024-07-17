package main

import (
	"fmt"
	"encoding/json"
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
