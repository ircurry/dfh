package ipc

import (
	"fmt"
	"testing"

	"github.com/ircurry/dfh/monitors"
)

func TestHyprMonString(t *testing.T) {
	type object struct {
		mon monitors.Monitor
		testState string
		expectedString string
		expectedBool bool
	}
	objects := []object{
		{
			monitors.Monitor{
				Name: "eDP-1",
				Width: 2256,
				Height: 1504,
				RefreshRate: 60,
				X: 0,
				Y: 0,
				Scale: 2,
				State: "undocked",
			},
			"undocked",
			"eDP-1,2256x1504@60,0x0,2",
			true,
		},
		{
			monitors.Monitor{
				Name: "eDP-1",
				Width: 2256,
				Height: 1504,
				RefreshRate: 60,
				X: 0,
				Y: 0,
				Scale: 2,
				State: "undocked",
			},
			"docked",
			"eDP-1,disable",
			false,
		},
	}

	for _, obj := range objects {
		monStr, b := hyprMonString(obj.mon, obj.testState)
		if monStr != obj.expectedString || b != obj.expectedBool {
			t.Errorf("\nExpected: %s\nGot: %s\n", obj.expectedString, monStr)
		}
	}
}

func TestStateStrings(t *testing.T) {
	type object struct {
		name string
		monl monitors.MonitorList
		testState string
		expectedStrings []string
		expectedErr error
	}

	objects := []object{
		{
			"enable eDP-1 and disable DP-2",
			monitors.MonitorList{
				monitors.Monitor{
					Name: "eDP-1",
					Width: 2256,
					Height: 1504,
					RefreshRate: 60,
					X: 0,
					Y: 0,
					Scale: 2,
					State: "undocked",
				},
				monitors.Monitor{
					Name: "DP-2",
					Width: 1920,
					Height: 1080,
					RefreshRate: 60,
					X: 0,
					Y: 0,
					Scale: 1,
					State: "docked",
				},
			},
			"undocked",
			[]string{"eDP-1,2256x1504@60,0x0,2","DP-2,disable"},
			nil,
		},
		{
			"disable eDP-1 and enable DP-2",
			monitors.MonitorList{
				monitors.Monitor{
					Name: "eDP-1",
					Width: 2256,
					Height: 1504,
					RefreshRate: 60,
					X: 0,
					Y: 0,
					Scale: 2,
					State: "undocked",
				},
				monitors.Monitor{
					Name: "DP-2",
					Width: 1920,
					Height: 1080,
					RefreshRate: 60,
					X: 0,
					Y: 0,
					Scale: 1,
					State: "docked",
				},
			},
			"fake",
			nil,
			fmt.Errorf("this monitor configuration does not contain state: fake"),
		},
	}
	for _, obj := range objects {
		t.Run(obj.name, func(t *testing.T) {
			strList, err := StateStrings(obj.monl, obj.testState)
			if obj.expectedErr != nil {
				if err != nil {
					if err.Error() != obj.expectedErr.Error() {
						t.Errorf("\nExpected: %s\nGot: %s\n", obj.expectedErr.Error(), err.Error())
					}
				} else {
					t.Error("Expected function to error but did not")
				}
			} else {
				if err != nil {
					t.Error(err.Error())
				}
				for i := range obj.expectedStrings {
					if obj.expectedStrings[i] != strList[i] {
						t.Errorf("\nExpected: %v\nGot: %v\n", obj.expectedStrings, strList)
					}
				}
			}
		})
	}
}
