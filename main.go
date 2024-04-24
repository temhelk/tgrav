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
	const timeStep float64 = 0.0001

	sim := simulation.NewSimulation(timeStep)

	sim.Bodies = []simulation.Body{
		{
			Mass:     1e12,
			Position: r2.Vec{X: 0, Y: 0},
			Velocity: r2.Vec{X: 0, Y: 0},
		},
		{
			Mass:     1e11,
			Position: r2.Vec{X: 0, Y: 15},
			Velocity: r2.Vec{X: 2.3, Y: 0},
		},
		{
			Mass:     1e10,
			Position: r2.Vec{X: 0, Y: 13.5},
			Velocity: r2.Vec{X: 0, Y: 0},
		},
		{
			Mass:     1e9,
			Position: r2.Vec{X: 0, Y: 4},
			Velocity: r2.Vec{X: 4, Y: 0},
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

	screen.EnableMouse()

	renderForceField := false
	rend := renderer.NewRenderer()

	// Ratio between simulation time and real time
	var simulationSpeed float64 = 1
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
				r := event.Rune()

				if key == tcell.KeyCtrlC || key == tcell.KeyEscape {
					break outer
				}

				if r == '+' {
					simulationSpeed *= 2
				} else if r == '-' {
					simulationSpeed /= 2
				}

				if r == 'f' {
					renderForceField = !renderForceField
				}
			case *tcell.EventMouse:
				buttons := event.Buttons()

				if buttons&tcell.WheelDown != 0 {
					rend.WorldWidth *= 1.2
				} else if buttons&tcell.WheelUp != 0 {
					rend.WorldWidth /= 1.2
				}
			}
		}

		newFrameTime := time.Now()
		deltaTime := newFrameTime.Sub(lastFrameTime)
		lastFrameTime = newFrameTime

		rend.AddFrameMessage(fmt.Sprintf("Δt: %.2f", deltaTime.Seconds()*1000))
		rend.AddFrameMessage(fmt.Sprintf("Speed: %.2f", simulationSpeed))

		simulationTimeAvailable += deltaTime.Seconds() * simulationSpeed
		for simulationTimeAvailable >= sim.TimeStep {
			simulationTimeAvailable -= sim.TimeStep

			sim.Step()
		}
		rend.AddFrameMessage(fmt.Sprintf("Step: %d", sim.SimulationStep))

		// totalEnergy := sim.CalculateTotalEnergy()
		// rend.AddFrameMessage(fmt.Sprintf("Total energy: %.2e", totalEnergy))

		// @TODO: Don't recalculate it all the time?
		centerOfMass := sim.CalculateCenterOfMass()
		rend.Center = centerOfMass
		// rend.AddFrameMessage(fmt.Sprintf("Center: <%.3e, %.3e>", centerOfMass.X, centerOfMass.Y))

		screen.Clear()

		if renderForceField {
			rend.RenderForceField(screen, sim)
		}

		rend.Render(screen, sim)
		screen.Show()

		sleepFor := targetFrameTime - time.Now().Sub(lastFrameTime)
		time.Sleep(sleepFor)
	}

	screen.Fini()
}
