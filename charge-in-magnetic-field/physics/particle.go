package physics

import "charge-in-magnetic-field/vector"

type Particle struct {
	Position vector.Vector `json:"Position"`
	Velocity vector.Vector `json:"Velocity"`
	Mass     float64       `json:"Mass"`
	Charge   float64       `json:"Charge"`
}

//Step uses numerical analysis to calculate particle's next position and velocity
func (p *Particle) Step(config Config) {
	force := vector.Add(config.MF.Force(*p), config.Fric.Force(p.Velocity))
	acceleration := vector.Scalar(force, 1/p.Mass)

	p.Position.Add(vector.Scalar(p.Velocity, config.TimeDelta.Seconds()))
	p.Velocity.Add(vector.Scalar(acceleration, config.TimeDelta.Seconds()))
}
