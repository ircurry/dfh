package monitors

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalProperFormat(t *testing.T) {
	var objects [][]byte = [][]byte{
		[]byte(`{"name": "eDP-1","width": 2256,"height": 1504,"refreshRate": 60,"x": 0,"y": 0,"scale": 2,"state": "dock"}`),
		[]byte(`{"width": 2256,"height": 1504,"refreshRate": 60,"x": 0,"y": 0,"scale": 2,"state": "dock","name": "eDP-1"}`),
		[]byte(`{"height": 1504,"refreshRate": 60,"x": 0,"y": 0,"scale": 2,"state": "dock","name": "eDP-1","width": 2256}`),
		[]byte(`{"refreshRate": 60,"x": 0,"y": 0,"scale": 2,"state": "dock","name": "eDP-1","width": 2256,"height": 1504}`),
		[]byte(`{"x": 0,"y": 0,"scale": 2,"state": "dock","name": "eDP-1","width": 2256,"height": 1504,"refreshRate": 60}`),
		[]byte(`{"y": 0,"scale": 2,"state": "dock","name": "eDP-1","width": 2256,"height": 1504,"refreshRate": 60,"x": 0}`),
		[]byte(`{"scale": 2,"state": "dock","name": "eDP-1","width": 2256,"height": 1504,"refreshRate": 60,"x": 0,"y": 0}`),
		[]byte(`{"state": "dock","name": "eDP-1","width": 2256,"height": 1504,"refreshRate": 60,"x": 0,"y": 0,"scale": 2}`),
		[]byte(
`{
    "name": "eDP-1",
    "width": 2256,
    "height": 1504,
    "refreshRate": 60,
    "x": 0,
    "y": 0,
    "scale": 2,
    "state": "dock"
}`),
	}
	base := Monitor{
		Name: "eDP-1",
		Width: 2256,
		Height: 1504,
		RefreshRate: 60,
		X: 0,
		Y: 0,
		Scale: 2,
		State: "dock",
	}
	for _, object := range objects {
		t.Log(string(object))
		mon := new(Monitor)
		err := json.Unmarshal(object, mon)
		if err != nil {
			t.Error(err.Error())
		}
		if !(base == *mon) {
			t.Errorf("Struct read and struct tested against do not equal\nRead:\n%s\n\nBase:\n%s\nRead Json: %s",
				mon.String(), base.String(), string(object))
		}
	}
}
