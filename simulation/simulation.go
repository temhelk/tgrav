package simulation

import (
	"gonum.org/v1/gonum/spatial/r2"
)

const G float64 = 6.674e-11

type Body struct {
	Mass float64
	Position r2.Vec
	Velocity r2.Vec

	appliedForce r2.Vec
}

type Simulation struct {
	TimeStep float64

	Bodies []Body
}

func NewSimulation() *Simulation {
	return &Simulation{
		TimeStep: 0.01,
	}
}

func (sim *Simulation) Step() {
	if sim.Bodies == nil {
		return
	}

	// Calculate forces for all bodies
	for bodyIndex := range sim.Bodies {
		body := &sim.Bodies[bodyIndex]

		body.appliedForce = r2.Vec{X: 0, Y: 0}

		for otherBodyIndex, otherBody := range sim.Bodies {
			if otherBodyIndex == bodyIndex {
				continue
			}

			bodyToOtherBody := r2.Sub(otherBody.Position, body.Position)
			forceAmplitude := G * (body.Mass * otherBody.Mass) / r2.Norm2(bodyToOtherBody)
			force := r2.Scale(forceAmplitude, r2.Unit(bodyToOtherBody))

			body.appliedForce = r2.Add(body.appliedForce, force)
		}
	}

	// Apply forces to all bodies
	for bodyIndex := range sim.Bodies {
		body := &sim.Bodies[bodyIndex]

		acceleration := r2.Scale(1 / body.Mass, body.appliedForce)

		body.Velocity = r2.Add(body.Velocity, r2.Scale(sim.TimeStep, acceleration))
		body.Position = r2.Add(body.Position, r2.Scale(sim.TimeStep, body.Velocity))
	}
}
