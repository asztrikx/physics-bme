package main

import (
	"charge-in-magnetic-field/draw"
	"charge-in-magnetic-field/physics"
	"charge-in-magnetic-field/vector"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

var config physics.Config
var bottomleft = vector.Create(0, 0)
var topright vector.Vector

//topright exclusive
func inside(pos vector.Vector) bool {
	return pos.X >= bottomleft.X &&
		pos.Y >= bottomleft.Y &&
		pos.X < topright.X &&
		pos.Y < topright.Y
}

func numerical() {
	//step particle with the current values for timeDelta time
	//then draw it
	part := config.Part
	for part.Velocity.Length() > config.VelocityBound && inside(part.Position) {
		draw.Draw(part.Position)

		part.Step(config)
	}
}

func velocityLength(t time.Duration) float64 {
	return math.Exp(-config.Fric.Coefficient/config.Part.Mass*t.Seconds()) * config.Part.Velocity.Length()
}

func radiusLength(t time.Duration) float64 {
	return config.Part.Mass / (config.Part.Charge * config.MF.Value) * velocityLength(t)
}

func radius(t time.Duration) vector.Vector {
	F_field := config.MF.Force(config.Part)
	rLength := config.Part.Mass * velocityLength(t) / (config.Part.Charge * config.MF.Value)
	r := vector.Scalar(F_field, rLength/F_field.Length())
	r.Negate()
	return r
}

func center() vector.Vector {
	r0 := radius(0)
	//radius vector goes from its origo
	r0.Negate()
	O := vector.Add(config.Part.Position, r0)
	return O
}

func alpha0() float64 {
	O := center()
	r0 := radiusLength(0)
	dx := config.Part.Position.X - O.X
	dy := config.Part.Position.Y - O.Y
	dxNorm := dx / r0
	dyNorm := dy / r0

	//float error correction
	if dxNorm < -1 {
		dxNorm = -1
	} else if dxNorm > 1 {
		dxNorm = 1
	}
	if dyNorm < -1 {
		dyNorm = -1
	} else if dyNorm > 1 {
		dyNorm = 1
	}

	alpha0cos := math.Acos(dxNorm)
	alpha0sin := math.Asin(dyNorm)
	if alpha0sin == 0 {
		return alpha0cos
	}
	if alpha0sin < 0 {
		return -alpha0cos
	}
	return alpha0cos
}

type Direction int

const (
	Left Direction = iota
	Right
	Straight
)

func direction(base, from, to vector.Vector) Direction {
	from.Sub(base)
	to.Sub(base)
	area := from.Y*to.X - from.X*to.Y
	if area < 0 {
		return Left
	}
	if area > 0 {
		return Right
	}
	return Straight
}

func omega() float64 {
	return config.Part.Charge * config.MF.Value / config.Part.Mass
}

func position(t time.Duration, O vector.Vector, w, a0 float64, pm Direction) vector.Vector {
	//parametric equatation
	da := w * t.Seconds()
	if pm == Right {
		da = -da
	}
	posX := O.X + radiusLength(t)*math.Cos(a0+da)
	posY := O.Y + radiusLength(t)*math.Sin(a0+da)
	return vector.Create(posX, posY)
}

func mathematical() {
	//constant values
	O := center()
	a0 := alpha0()
	pm := direction(vector.Create(0, 0), radius(0), config.Part.Velocity)
	if pm == Straight {
		panic("vec is not perpendicular to mf force")
	}
	w := omega()
	tEnd := -config.Part.Mass / config.Fric.Coefficient * math.Log(config.VelocityBound/config.Part.Velocity.Length())

	//calculate position for each time to draw it for comparision
	t := time.Duration(0)
	pos := position(t, O, w, a0, pm)
	//always recalulate t as timeDelta is small
	for i := 0; t.Seconds() < tEnd && inside(pos); i++ {
		draw.Draw(pos)

		t = time.Duration(i+1) * config.TimeDelta
		pos = position(t, O, w, a0, pm)
	}

	//Geogebra format for equatation
	dir := "+"
	if pm == Right {
		dir = "-"
	}
	fmt.Print("Curve((")
	fmt.Printf(
		"%f + %f/(%f*%f) * e^(-%f/%f * t) * %f * cos(%f %s %f * t)",
		O.X,
		config.Part.Mass, config.Part.Charge, config.MF.Value,
		config.Fric.Coefficient, config.Part.Mass, config.Part.Velocity.Length(),
		a0, dir, w,
	)
	fmt.Print(",")
	fmt.Printf(
		"%f + %f/(%f*%f) * e^(-%f/%f * t) * %f * sin(%f %s %f * t)",
		O.Y,
		config.Part.Mass, config.Part.Charge, config.MF.Value,
		config.Fric.Coefficient, config.Part.Mass, config.Part.Velocity.Length(),
		a0, dir, w,
	)
	fmt.Print("),")
	fmt.Printf("t,0,%f", t.Seconds())
	fmt.Println(")")
}

func main() {
	//read json
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("can not open")
	}
	if json.NewDecoder(file).Decode(&config) != nil {
		log.Fatalln("can not parse")
	}
	file.Close()
	config.TimeDelta = time.Duration(config.TimeDeltaJSON)

	//check if vec pos on mf
	topright = vector.Create(float64(config.MF.Width), float64(config.MF.Height))
	if !inside(config.Part.Position) {
		panic("starting position is not inside magnetic field")
	}

	//draw
	draw.Start(config)

	draw.SetRGB(1, 0, 0)
	numerical()

	draw.SetRGB(0, 1, 0)
	mathematical()

	draw.Save()
}
