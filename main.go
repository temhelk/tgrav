package main

import (
	"log"
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
			Mass:     1e13,
			Position: r2.Vec{X: 10, Y: 0},
			Velocity: r2.Vec{X: 0, Y: -3},
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

		sim.Step()

		screen.Clear()
		rend.AddFrameMessage("Hello")
		rend.Render(screen, defaultStyle, sim)
		screen.Show()

		time.Sleep(time.Second / 100)
	}

	screen.Fini()
}
