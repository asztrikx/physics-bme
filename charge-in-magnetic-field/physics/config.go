package physics

import (
	"encoding/json"
	"errors"
	"time"
)

type Config struct {
	MF   MagneticField `json:"MagneticField"`
	Part Particle      `json:"Particle"`
	Fric Friction      `json:"Friction"`
	//https://github.com/golang/go/issues/10275
	TimeDeltaJSON Duration `json:"TimeDelta"`
	TimeDelta     time.Duration
	VelocityBound float64 `json:"VelocityBound"`
	//specifies draw.buffer size
	RAM int `json:"RAM"`
}

type Duration time.Duration

//UnmarshalJSON parses string from json to time.Duration
func (d *Duration) UnmarshalJSON(b []byte) error {
	var text string
	if err := json.Unmarshal(b, &text); err != nil {
		return err
	}

	t, err := time.ParseDuration(text)
	if err != nil {
		return errors.New(`a duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h"`)
	}
	*d = Duration(t)

	return nil
}
