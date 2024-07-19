package main

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

func (mon *Monitor) String() string {
	val, err := json.Marshal(mon)
	if err != nil {
		panic(err)
	}
	return string(val)
}

func (mon *Monitor) CheckStateField() error {
	if mon.State == nil {
		return fmt.Errorf("cannont use monitor cofiguration, does not contain a State")
	} else {
		return nil
	}
}

func (mon *Monitor) CheckStringFields() error {
	switch v:= true; v {
	case (mon.Name == nil):
		return fmt.Errorf("cannot format string, Name is nil.")
	case (mon.Width == nil):
		return fmt.Errorf("cannot format string, Width is nil.")
	case (mon.Height == nil):
		return fmt.Errorf("cannot format string, Height is nil.")
	case (mon.RefreshRate == nil):
		return fmt.Errorf("cannot format string, Refresh Rate is nil.")
	case (mon.X == nil):
		return fmt.Errorf("cannot format string, X Position is nil.")
	case (mon.Y == nil):
		return fmt.Errorf("cannot format string, Y Position is nil.")
	case (mon.Scale == nil):
		return fmt.Errorf("cannot format string, Scale is nil.")
	default:
		return nil
	}
}

func (mon *Monitor) EnableString() (string, error) {
	var name, resolution, position, scale string
	if err := mon.CheckStringFields(); err != nil {
		return "", err
	} else {
		name = *mon.Name
		resolution = fmt.Sprintf("%dx%d@%d", *mon.Width, *mon.Height, *mon.RefreshRate)
		position = fmt.Sprintf("%dx%d", *mon.X, *mon.Y)
		scale = fmt.Sprintf("%d", *mon.Scale)
	}
	
	hyprstring := fmt.Sprintf("%s,%s,%s,%s", name, resolution, position, scale)
	return hyprstring, nil
}

func (mon *Monitor) DisableString() (string, error) {
	if err := mon.CheckStringFields(); err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%s,disable", *mon.Name), nil
	}
}

type MonitorList []Monitor

func (monlist *MonitorList) fromJson(byt []byte) error {
	return json.Unmarshal(byt, monlist)
}

func (monlist *MonitorList) stateStrings(state string) ([]string, error) {
	strlist := make([]string, 0)
	enablelist := make([]string, 0)
	disablelist := make([]string, 0)
	var stateErr error = nil
	containsState := false
	for _, mon := range *monlist {
		if err := mon.CheckStateField(); err != nil {
			return nil, err
		}
		if *mon.State == state {
			str, err := mon.EnableString()
			if err != nil {
				return nil, err
			}
			enablelist = append(enablelist, str)
			containsState = true
		} else {
			str, err := mon.DisableString()
			if err != nil {
				return nil, err
			}
			disablelist = append(disablelist, str)
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

func compareMonitorLists (monl MonitorList, names []string) (bool, error) {
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
