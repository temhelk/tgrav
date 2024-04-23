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
	SimulationStep uint64

	Bodies []Body
}

func NewSimulation(timeStep float64) *Simulation {
	return &Simulation{
		TimeStep: timeStep,
	}
}

func (sim *Simulation) Step() {
	sim.SimulationStep += 1

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

func (sim *Simulation) CalculateCenterOfMass() r2.Vec {
	var totalMass float64

	for _, body := range sim.Bodies {
		totalMass += body.Mass
	}

	var centerOfMass r2.Vec

	for _, body := range sim.Bodies {
		centerOfMass = r2.Add(centerOfMass, r2.Scale(body.Mass / totalMass, body.Position))
	}

	return centerOfMass
}

func (sim *Simulation) CalculateTotalEnergy() float64 {
	potentialEnergy := 0.0
	for body1Index, body1 := range sim.Bodies {
		for body2Index, body2 := range sim.Bodies {
			if body1Index == body2Index {
				continue
			}

			body2ToBody1 := r2.Sub(body1.Position, body2.Position)
			potentialEnergy += -G * body1.Mass * body2.Mass / r2.Norm(body2ToBody1)
		}
	}
	potentialEnergy /= 2

	kineticEnergy := 0.0
	for _, body := range sim.Bodies {
		kineticEnergy += body.Mass * r2.Norm2(body.Velocity) / 2
	}

	return potentialEnergy + kineticEnergy
}
