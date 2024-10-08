package monitors

import ()

type Resolution struct {
	Width       int64 `json:"width"`
	Height      int64 `json:"height"`
	RefreshRate uint8 `json:"refresh_rate"`
}

type Position struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

type Monitor struct {
	Name    *string     `json:"name"`
	Res     *Resolution `json:"resolution"`
	Pos     *Position   `json:"position"`
	Scale   *float32    `json:"scale"`
	Enabled bool        `json:"enabled"`
}

type Profile struct {
	Name     string    `json:"name"`
	Monitors []Monitor `json:"monitors"`
}
