package vector

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"strings"
)

type Vector struct {
	X, Y float64
}

func Create(x, y float64) Vector {
	return Vector{
		X: x,
		Y: y,
	}
}
func (v *Vector) Negate() *Vector {
	v.X = -v.X
	v.Y = -v.Y
	return v
}
func (v *Vector) Scalar(value float64) *Vector {
	v.X *= value
	v.Y *= value
	return v
}
func Scalar(v Vector, value float64) Vector {
	return *v.Scalar(value)
}
func (v *Vector) RotateP90() {
	v.X, v.Y = -v.Y, v.X
}
func (v *Vector) RotateN90() {
	v.X, v.Y = v.Y, -v.X
}
func (v *Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}
func (v *Vector) SetLength(length float64) *Vector {
	return v.Scalar(length / v.Length())
}
func (v *Vector) Add(v2 Vector) *Vector {
	v.X += v2.X
	v.Y += v2.Y
	return v
}
func Add(v1 Vector, v2 Vector) Vector {
	return *v1.Add(v2)
}
func (v *Vector) Sub(v2 Vector) *Vector {
	v.X -= v2.X
	v.Y -= v2.Y
	return v
}
func Sub(v1 Vector, v2 Vector) Vector {
	return *v1.Sub(v2)
}

//UnmarshalJSON parses string from json to vector
func (v *Vector) UnmarshalJSON(b []byte) error {
	var text string
	if err := json.Unmarshal(b, &text); err != nil {
		return err
	}
	cooS := strings.Split(text, " ")
	if len(cooS) != 2 {
		return errors.New("coordinate should have two values")
	}

	var err error
	if v.X, err = strconv.ParseFloat(cooS[0], 64); err != nil {
		return errors.New("could not parse X")
	}
	if v.Y, err = strconv.ParseFloat(cooS[1], 64); err != nil {
		return errors.New("could not parse Y")
	}
	return nil
}
