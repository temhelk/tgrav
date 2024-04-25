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

	sim.Bodies = simulation.FourBodySystem[:]
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

	screenDragging := false
	var previousMouseX, previousMouseY int
	var worldOffset r2.Vec

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
				x, y := event.Position()

				if buttons&tcell.WheelDown != 0 {
					rend.WorldWidth *= 1.2
				} else if buttons&tcell.WheelUp != 0 {
					rend.WorldWidth /= 1.2
				}

				if screenDragging {
					offsetX, offsetY := x - previousMouseX, y - previousMouseY

					screenOffset := r2.Vec{X: float64(-offsetX), Y: float64(offsetY)}

					worldOffset = r2.Add(worldOffset, rend.CellDirToWorld(screen, screenOffset))

					previousMouseX = x
					previousMouseY = y
				}

				if !screenDragging && (buttons&tcell.Button1 != 0) {
					screenDragging = true
					previousMouseX, previousMouseY = x, y
				} else if screenDragging && (buttons&tcell.Button1 == 0) {
					screenDragging = false
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

		rend.Center = r2.Add(centerOfMass, worldOffset)
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
