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
		Name:        "eDP-1",
		Width:       2256,
		Height:      1504,
		RefreshRate: 60,
		X:           0,
		Y:           0,
		Scale:       2,
		State:       "dock",
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

func TestUnmarshalMissingKey(t *testing.T) {
	type missingKey struct {
		key  string
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
		Name:        "eDP-1",
		Width:       2256,
		Height:      1504,
		RefreshRate: 60,
		X:           0,
		Y:           0,
		Scale:       2,
		State:       "dock",
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
		data     []byte
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
		data    []byte
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

func TestUnmarshalWithWrongFieldTypes(t *testing.T) {
	type object struct {
		data       []byte
		key        string
		expectType string
		actualType string
	}
	fmtInvalidTypeErr := func(str string, valType string, v any) error {
		return fmt.Errorf("key %s expect type %s, found %s", str, valType, v)
	}

	var objects []object = []object{
		//name
		{
			[]byte(`{"name": 42}`),
			"name",
			"string",
			"float64",
		},
		{
			[]byte(`{"name": null}`),
			"name",
			"string",
			"<nil>",
		},
		{
			[]byte(`{"name": true}`),
			"name",
			"string",
			"bool",
		},
		{
			[]byte(`{"name": {}}`),
			"name",
			"string",
			"json.Delim",
		},
		//width
		{
			[]byte(`{"width": "asdf"}`),
			"width",
			"int64",
			"string",
		},
		{
			[]byte(`{"width": null}`),
			"width",
			"int64",
			"<nil>",
		},
		{
			[]byte(`{"width": true}`),
			"width",
			"int64",
			"bool",
		},
		{
			[]byte(`{"width": {}}`),
			"width",
			"int64",
			"json.Delim",
		},
		//height
		{
			[]byte(`{"height": "asdf"}`),
			"height",
			"int64",
			"string",
		},
		{
			[]byte(`{"height": null}`),
			"height",
			"int64",
			"<nil>",
		},
		{
			[]byte(`{"height": true}`),
			"height",
			"int64",
			"bool",
		},
		{
			[]byte(`{"height": {}}`),
			"height",
			"int64",
			"json.Delim",
		},
		//refreshRate
		{
			[]byte(`{"refreshRate": "asdf"}`),
			"refreshRate",
			"int",
			"string",
		},
		{
			[]byte(`{"refreshRate": null}`),
			"refreshRate",
			"int",
			"<nil>",
		},
		{
			[]byte(`{"refreshRate": true}`),
			"refreshRate",
			"int",
			"bool",
		},
		{
			[]byte(`{"refreshRate": {}}`),
			"refreshRate",
			"int",
			"json.Delim",
		},
		//x
		{
			[]byte(`{"x": "asdf"}`),
			"x",
			"int64",
			"string",
		},
		{
			[]byte(`{"x": null}`),
			"x",
			"int64",
			"<nil>",
		},
		{
			[]byte(`{"x": true}`),
			"x",
			"int64",
			"bool",
		},
		{
			[]byte(`{"x": {}}`),
			"x",
			"int64",
			"json.Delim",
		},
		//y
		{
			[]byte(`{"y": "asdf"}`),
			"y",
			"int64",
			"string",
		},
		{
			[]byte(`{"y": null}`),
			"y",
			"int64",
			"<nil>",
		},
		{
			[]byte(`{"y": true}`),
			"y",
			"int64",
			"bool",
		},
		{
			[]byte(`{"y": {}}`),
			"y",
			"int64",
			"json.Delim",
		},
		//scale
		{
			[]byte(`{"scale": "asdf"}`),
			"scale",
			"int",
			"string",
		},
		{
			[]byte(`{"scale": null}`),
			"scale",
			"int",
			"<nil>",
		},
		{
			[]byte(`{"scale": true}`),
			"scale",
			"int",
			"bool",
		},
		{
			[]byte(`{"scale": {}}`),
			"scale",
			"int",
			"json.Delim",
		},
		//state
		{
			[]byte(`{"state": 42}`),
			"state",
			"string",
			"float64",
		},
		{
			[]byte(`{"state": null}`),
			"state",
			"string",
			"<nil>",
		},
		{
			[]byte(`{"state": true}`),
			"state",
			"string",
			"bool",
		},
		{
			[]byte(`{"state": {}}`),
			"state",
			"string",
			"json.Delim",
		},
	}
	for _, obj := range objects {
		t.Log(string(obj.data))
		errWant := fmtInvalidTypeErr(obj.key, obj.expectType, obj.actualType)
		t.Log(errWant.Error())
		var mon Monitor
		err := mon.UnmarshalJSON(obj.data)

		if err.Error() != errWant.Error() {
			t.Errorf("Expected error \"%s\" but got \"%s\"", errWant.Error(), err.Error())
		}
	}
}

func TestUnmarshallMonitorList(t *testing.T) {
	data := []byte(`[{"name": "eDP-1","width": 2256,"height": 1504,"refreshRate": 60,"x": 0,"y": 0,"scale": 2,"state": "undock"},{"name": "DP-2","width": 1920,"height": 1080,"refreshRate": 60,"x": 0,"y": 0,"scale": 1,"state": "dock"}]`)
	expectedValue := []Monitor{
		{"eDP-1", 2256, 1504, 60, 0, 0, 2, "undock"},
		{"DP-2", 1920, 1080, 60, 0, 0, 1, "dock"},
	}
	var monl MonitorList
	err := monl.FromJson(data)
	if err != nil {
		t.Error(err.Error())
	}
	for i := range monl {
		if monl[i] != expectedValue[i] {
			t.Errorf("\nExpected value: %v\nGot value: %v\n", expectedValue, monl)
		}
	}
}

func TestCompareMonitorValues(t *testing.T) {
	monl := MonitorList{
		{"eDP-1", 2256, 1504, 60, 0, 0, 2, "undock"},
		{"DP-2", 1920, 1080, 60, 0, 0, 1, "dock"},
	}
	names := []string{"eDP-1", "DP-2"}

	b, err := CompareMonitorLists(monl, names)
	if err != nil {
		t.Error(err.Error())
	}
	if !b {
		t.Error("Expected all monitors to be present but found that some were missing")
	}
}
