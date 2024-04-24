package simulation

import (
	"gonum.org/v1/gonum/spatial/r2"
)

const G float64 = 6.674e-11

type Body struct {
	Mass     float64

	Acceleration r2.Vec
	Position r2.Vec
	Velocity r2.Vec
}

type Simulation struct {
	TimeStep       float64
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

		body.Acceleration = r2.Vec{X: 0, Y: 0}

		for otherBodyIndex, otherBody := range sim.Bodies {
			if otherBodyIndex == bodyIndex {
				continue
			}

			bodyToOtherBody := r2.Sub(otherBody.Position, body.Position)
			accelerationAmplitude := G * otherBody.Mass / r2.Norm2(bodyToOtherBody)
			acceleration := r2.Scale(accelerationAmplitude, r2.Unit(bodyToOtherBody))

			body.Acceleration = r2.Add(body.Acceleration, acceleration)
		}
	}

	// Apply acceleration to all bodies
	for bodyIndex := range sim.Bodies {
		body := &sim.Bodies[bodyIndex]

		body.Velocity = r2.Add(body.Velocity, r2.Scale(sim.TimeStep, body.Acceleration))
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
		centerOfMass = r2.Add(centerOfMass, r2.Scale(body.Mass/totalMass, body.Position))
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

func (sim *Simulation) CalculateAccelerationAt(pos r2.Vec) r2.Vec {
	var totalAcceleration r2.Vec

	for _, body := range sim.Bodies {
		posToBody := r2.Sub(body.Position, pos)
		accelerationAmplitude := G * body.Mass / r2.Norm2(posToBody)
		acceleration := r2.Scale(accelerationAmplitude, r2.Unit(posToBody))

		totalAcceleration = r2.Add(totalAcceleration, acceleration)
	}

	return totalAcceleration
}
