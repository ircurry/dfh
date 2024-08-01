package monitors

import (
	"encoding/json"
	"errors"
	"fmt"
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
			t.Errorf("Struct read and struct tested against do not equal\n[[Read]]\n%s\n\n[[Base]]\n%s\n[[Read Json]]\n%s",
				mon.String(), base.String(), string(object))
		}
	}
}

//2:15:04
func TestUnmarshalMissingKey(t *testing.T) {
	type missingKey struct {
		key string
		json []byte
	}
	var objects []missingKey = []missingKey{
		{
			"name",
			[]byte(`{"width": 2256,"height": 1504,"refreshRate": 60,"x": 0,"y": 0,"scale": 2,"state": "dock"}`),
		},
		{
			"width",
			[]byte(`{"name": "eDP-1","height": 1504,"refreshRate": 60,"x": 0,"y": 0,"scale": 2,"state": "dock"}`),
		},
		{
			"height",
			[]byte(`{"name": "eDP-1","width": 2256,"refreshRate": 60,"x": 0,"y": 0,"scale": 2,"state": "dock"}`),
		},
		{
			"refreshRate",
			[]byte(`{"name": "eDP-1","width": 2256,"height": 1504,"x": 0,"y": 0,"scale": 2,"state": "dock"}`),
		},
		{
			"x",
			[]byte(`{"name": "eDP-1","width": 2256,"height": 1504,"refreshRate": 60,"y": 0,"scale": 2,"state": "dock"}`),
		},
		{
			"y",
			[]byte(`{"name": "eDP-1","width": 2256,"height": 1504,"refreshRate": 60,"x": 0,"scale": 2,"state": "dock"}`),
		},
		{
			"scale",
			[]byte(`{"name": "eDP-1","width": 2256,"height": 1504,"refreshRate": 60,"x": 0,"y": 0,"state": "dock"}`),
		},
		{
			"state",
			[]byte(`{"name": "eDP-1","width": 2256,"height": 1504,"refreshRate": 60,"x": 0,"y": 0,"scale": 2}`),
		},
	}
	for _, object := range objects {
		t.Log(string(object.json))
		mon := new(Monitor)
		err := json.Unmarshal(object.json, mon)
		errWant := fmt.Errorf("Key %s not set", object.key)
		if err.Error() != errWant.Error() {
			t.Errorf("Errors not equal\nWanted '%s'\nGot '%s'", errWant.Error(), err.Error())
		}
	}
}

func TestUnmarshalWithExtraFields(t *testing.T) {
	var objects [][]byte = [][]byte{
		[]byte(`{"name": "eDP-1","width": 2256,"height": 1504,"refreshRate": 60,"x": 0,"y": 0,"scale": 2,"state": "dock"}`),
		[]byte(`{"extra": "something","name": "eDP-1","width": 2256,"height": 1504,"refreshRate": 60,"x": 0,"y": 0,"scale": 2,"state": "dock"}`),
		[]byte(`{"name": "eDP-1","extra": "something","width": 2256,"height": 1504,"refreshRate": 60,"x": 0,"y": 0,"scale": 2,"state": "dock"}`),
		[]byte(`{"name": "eDP-1","width": 2256,"extra": "something","height": 1504,"refreshRate": 60,"x": 0,"y": 0,"scale": 2,"state": "dock"}`),
		[]byte(`{"name": "eDP-1","width": 2256,"height": 1504,"extra": "something","refreshRate": 60,"x": 0,"y": 0,"scale": 2,"state": "dock"}`),
		[]byte(`{"name": "eDP-1","width": 2256,"height": 1504,"refreshRate": 60,"extra": "something","x": 0,"y": 0,"scale": 2,"state": "dock"}`),
		[]byte(`{"name": "eDP-1","width": 2256,"height": 1504,"refreshRate": 60,"x": 0,"extra": "something","y": 0,"scale": 2,"state": "dock"}`),
		[]byte(`{"name": "eDP-1","width": 2256,"height": 1504,"refreshRate": 60,"x": 0,"y": 0,"extra": "something","scale": 2,"state": "dock"}`),
		[]byte(`{"name": "eDP-1","width": 2256,"height": 1504,"refreshRate": 60,"x": 0,"y": 0,"scale": 2,"extra": "something","state": "dock"}`),
		[]byte(`{"name": "eDP-1","width": 2256,"height": 1504,"refreshRate": 60,"x": 0,"y": 0,"scale": 2,"state": "dock","extra": "something"}`),
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
			t.Errorf("Struct read and struct tested against do not equal\n[[Read]]\n%s\n\n[[Base]]\n%s\n[[Read Json]]\n%s",
				mon.String(), base.String(), string(object))
		}
	}
}

func TestUnmarshalWithEmpty(t *testing.T) {
	var mon Monitor
	dat := []byte("")
	err := mon.UnmarshalJSON(dat)
	errWant := errors.New("Reached EOF before any json tokens could be read")
	if err == nil {
		t.Error("Expected function to error but did not")
	}
	if err.Error() != errWant.Error() {
		t.Errorf("Expected error \"%s\" but got \"%s\"", errWant.Error(), err.Error())
	}
}

func TestUnmarshalWithInitialDelim(t *testing.T) {
	type object struct {
		data []byte
		dataType string
	}
	var objects []object = []object{
		{
			[]byte(`"asdf"`),
			"string",
		},
		{
			[]byte(`1`),
			"float64",
		},
		{
			[]byte(`4.25`),
			"float64",
		},
		{
			[]byte(`true`),
			"bool",
		},
		{
			[]byte(`null`),
			"<nil>",
		},
	}
	for _, obj := range objects {
		t.Log(string(obj.data))
		t.Log(string(obj.dataType))
		var mon Monitor
		err := mon.UnmarshalJSON(obj.data)
		errWant := fmt.Errorf("first token was not a delimeter, instead was %s", obj.dataType)
		
		if err.Error() != errWant.Error() {
			t.Errorf("Expected error \"%s\" but got \"%s\"", errWant.Error(), err.Error())
		}
	}
}

func TestUnmarshalWithNonOpenBracketDelim(t *testing.T) {
	type object struct {
		data []byte
		errWant error
	}
	var objects []object = []object{
		{
			[]byte(`[]`),
			fmt.Errorf("JSON token is not an begin-object, instead was ["),
		},
	}
	for _, obj := range objects {
		t.Log(string(obj.data))
		t.Log(string(obj.errWant.Error()))
		var mon Monitor
		err := mon.UnmarshalJSON(obj.data)

		if err.Error() != obj.errWant.Error() {
			t.Errorf("Expected error \"%s\" but got \"%s\"", obj.errWant.Error(), err.Error())
		}
	}
}
