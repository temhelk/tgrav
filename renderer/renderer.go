package renderer

import (
	"math"

	"github.com/temhelk/tgrav/simulation"

	"github.com/gdamore/tcell/v2"
)

type Renderer struct {
	frameMessage string
}

func NewRenderer() *Renderer {
	return &Renderer{
		frameMessage: "",
	}
}

func (rend *Renderer) Render(screen tcell.Screen, style tcell.Style, sim *simulation.Simulation) {
	width, height := screen.Size()

	for _, body := range sim.Bodies {
		x := body.Position.X + (float64(width) / 2)
		y := body.Position.Y + (float64(height) / 2)

		xDecimal := x - math.Floor(x)
		yDecimal := y - math.Floor(y)

		xInt := int(x)
		yInt := int(y)

		if xInt >= 0 && xInt < width && yInt >= 0 && yInt < height {
			xPart := clampInt(int(xDecimal * 2), 0, 1)
			yPart := clampInt(int(yDecimal * 4), 0, 3)

			partNumber := yPart + xPart * 4

			symbol := makeBraille(partNumber)

			screen.SetContent(xInt, yInt, symbol, nil, style)
		}
	}

	rend.writeString(screen, 0, height - 1, style, rend.frameMessage)
	rend.frameMessage = ""
}

func (rend *Renderer) AddFrameMessage(message string) {
	rend.frameMessage += message
}

func (rend *Renderer) writeString(screen tcell.Screen, x, y int, style tcell.Style, str string) {
	width, _ := screen.Size()

	for index, r := range str {
		if (x + index >= width) {
			return
		}

		screen.SetContent(x + index, y, r, nil, style)
	}
}

func makeBraille(partNumber int) rune {
	unicodeOffset := 0
	switch (partNumber) {
	case 0:
		unicodeOffset = 0x1
	case 1:
		unicodeOffset = 0x2
	case 2:
		unicodeOffset = 0x4
	case 3:
		unicodeOffset = 0x40
	case 4:
		unicodeOffset = 0x8
	case 5:
		unicodeOffset = 0x10
	case 6:
		unicodeOffset = 0x20
	case 7:
		unicodeOffset = 0x80
	}

	return rune(0x2800 + unicodeOffset)
}

func clampInt(n, a, b int) int {
	if n <= a {
		return a
	} else if n >= b {
		return b
	} else {
		return n
	}
}