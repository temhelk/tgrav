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
	sim := simulation.NewSimulation()

	sim.Bodies = []simulation.Body{
		{
			Mass:     1e13,
			Position: r2.Vec{X: -10, Y: 0},
			Velocity: r2.Vec{X: 0, Y: 3},
		},
		{
			Mass:     1e12,
			Position: r2.Vec{X: 10, Y: 0},
			Velocity: r2.Vec{X: 0, Y: -3},
		},
		{
			Mass:     1e12,
			Position: r2.Vec{X: 0, Y: 5},
			Velocity: r2.Vec{X: 2, Y: -4},
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

		sim.Step()

		totalEnergy := sim.CalculateTotalEnergy()
		rend.AddFrameMessage(fmt.Sprintf("Total energy: %e", totalEnergy))

		// @TODO: Don't recalculate it all the time?
		centerOfMass := sim.CalculateCenterOfMass()
		rend.SetCenter(centerOfMass)
		rend.AddFrameMessage(fmt.Sprintf("Center: %+v", centerOfMass))

		screen.Clear()
		rend.Render(screen, defaultStyle, sim)
		screen.Show()

		sleepFor := targetFrameTime - time.Now().Sub(lastFrameTime)
		time.Sleep(sleepFor)
	}

	screen.Fini()
}
