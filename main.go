package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/temhelk/tgrav/renderer"
	"github.com/temhelk/tgrav/simulation"

	"github.com/gdamore/tcell/v2"
	"gonum.org/v1/gonum/spatial/r2"
)

func main() {
	const timeStep float64 = 0.0001;

	sim := simulation.NewSimulation(timeStep)

	sim.Bodies = []simulation.Body{
		{
			Mass:     1e12,
			Position: r2.Vec{X: 0, Y: 0},
			Velocity: r2.Vec{X: 0, Y: 0},
		},
		{
			Mass:     1e12,
			Position: r2.Vec{X: 15, Y: 0},
			Velocity: r2.Vec{X: 0, Y: 3},
		},
		{
			Mass:     1e12,
			Position: r2.Vec{X: 0, Y: 11},
			Velocity: r2.Vec{X: -2, Y: 0},
		},
	}

	screen, err := tcell.NewScreen()

	if err != nil {
		log.Panicf("%+v", err)
	}

	if err := screen.Init(); err != nil {
		log.Panicf("%+v", err)
	}

	defaultStyle := tcell.StyleDefault
	screen.SetStyle(defaultStyle)

	rend := renderer.NewRenderer()

	// Ratio between simulation time and real time
	var simulationSpeed float64 = 3
	var simulationTimeAvailable float64

	targetFrameTime := time.Duration(math.Floor(1.0 / 60 * float64(time.Second)))
	lastFrameTime := time.Now()

outer:
	for {
		for screen.HasPendingEvent() {
			event := screen.PollEvent()

			switch event := event.(type) {
			case *tcell.EventKey:
				key := event.Key()

				if key == tcell.KeyCtrlC || key == tcell.KeyEscape {
					break outer
				}
			}
		}

		newFrameTime := time.Now()
		deltaTime := newFrameTime.Sub(lastFrameTime)
		lastFrameTime = newFrameTime

		rend.AddFrameMessage(fmt.Sprintf("dt: %.2f", deltaTime.Seconds() * 1000))

		simulationTimeAvailable += deltaTime.Seconds() * simulationSpeed
		for simulationTimeAvailable >= sim.TimeStep {
			simulationTimeAvailable -= sim.TimeStep

			sim.Step()
		}
		rend.AddFrameMessage(fmt.Sprintf("Step: %d", sim.SimulationStep))

		totalEnergy := sim.CalculateTotalEnergy()
		rend.AddFrameMessage(fmt.Sprintf("Total energy: %.2e", totalEnergy))

		// @TODO: Don't recalculate it all the time?
		centerOfMass := sim.CalculateCenterOfMass()
		rend.Center = centerOfMass
		// rend.AddFrameMessage(fmt.Sprintf("Center: <%.3e, %.3e>", centerOfMass.X, centerOfMass.Y))

		screen.Clear()
		rend.Render(screen, defaultStyle, sim)
		screen.Show()

		sleepFor := targetFrameTime - time.Now().Sub(lastFrameTime)
		time.Sleep(sleepFor)
	}

	screen.Fini()
}
