package physics

import (
	"charge-in-magnetic-field/vector"
	"encoding/json"
	"errors"
	"strings"
)

type Direction int

const (
	Up Direction = iota
	Down
)

type MagneticField struct {
	Width  uint      `json:"Width"`
	Height uint      `json:"Height"`
	Dir    Direction `json:"Direction"`
	Value  float64   `json:"Value"`
}

func (mf *MagneticField) Force(p Particle) vector.Vector {
	force := p.Velocity
	force.SetLength(p.Charge * p.Velocity.Length() * mf.Value)

	//determine rotation direction
	var positive bool
	if mf.Dir == Up {
		positive = false
	} else {
		positive = true
	}
	if p.Charge < 0 {
		positive = !positive
	}

	if positive {
		force.RotateP90()
	} else {
		force.RotateN90()
	}

	return force
}

//UnmarshalJSON parses string from json to Direction
func (d *Direction) UnmarshalJSON(b []byte) error {
	var text string
	if err := json.Unmarshal(b, &text); err != nil {
		return err
	}
	text = strings.ToLower(text)
	if text == "up" {
		*d = Up
	} else if text == "down" {
		*d = Down
	} else {
		return errors.New("Direction isn't up nor down")
	}
	return nil
}
