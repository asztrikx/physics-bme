package physics

import "charge-in-magnetic-field/vector"

type Friction struct {
	//Coefficient should be positive here
	Coefficient float64 `json:"Coefficient"`
}

func (f *Friction) Force(velocity vector.Vector) vector.Vector {
	return *velocity.Scalar(f.Coefficient).Negate()
}
