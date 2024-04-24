package renderer

import (
	"cmp"
	"math"
	"slices"

	"github.com/temhelk/tgrav/simulation"

	"github.com/gdamore/tcell/v2"
	"gonum.org/v1/gonum/spatial/r2"
)

type Renderer struct {
	Center     r2.Vec
	WorldWidth float64

	frameMessage string
}

func NewRenderer() *Renderer {
	return &Renderer{
		WorldWidth: 100,
	}
}

func (rend *Renderer) Render(screen tcell.Screen, sim *simulation.Simulation) {
	defaultStyle := tcell.StyleDefault
	width, height := screen.Size()

	worldHeight := (float64(height) / float64(width)) * rend.WorldWidth

	scaleX := float64(width) / rend.WorldWidth

	// Scale y by 0.497 to adjust for non square character terminal (adjusted for my font)
	scaleY := float64(height) / worldHeight * 0.497

	for _, body := range sim.Bodies {
		x := (body.Position.X-rend.Center.X)*scaleX + (float64(width) / 2)
		y := (body.Position.Y-rend.Center.Y)*scaleY + (float64(height) / 2)

		xDecimal := x - math.Floor(x)
		yDecimal := y - math.Floor(y)

		xInt := int(x)
		yInt := height - int(y) - 1

		if xInt >= 0 && xInt < width && yInt >= 0 && yInt < height {
			xPart := clamp(int(xDecimal*2), 0, 1)
			yPart := 3 - clamp(int(yDecimal*4), 0, 3)

			partNumber := yPart + xPart*4

			existingSymbol, _, style, _ := screen.GetContent(xInt, yInt)

			newDotSymbol := makeBraille(partNumber)

			if existingSymbol == ' ' {
				screen.SetContent(xInt, yInt, newDotSymbol, nil, style)
			} else {
				combinedSymbol := combineBraille(existingSymbol, newDotSymbol)
				screen.SetContent(xInt, yInt, combinedSymbol, nil, style)
			}
		}
	}

	rend.writeString(screen, 0, height-1, defaultStyle, rend.frameMessage)
	rend.frameMessage = ""
}

func (rend *Renderer) RenderForceField(screen tcell.Screen, sim *simulation.Simulation) {
	width, height := screen.Size()

	forceValues := make([]float64, width*height)

	for y := range height {
		for x := range width {
			offsets := [4]r2.Vec{
				{X: 0, Y: 0},
				{X: 0, Y: 1},
				{X: 1, Y: 0},
				{X: 1, Y: 1},
			}

			var cornerValues [len(offsets)]float64

			var accelerationSum float64
			for index, offset := range offsets {
				cellPos := r2.Add(r2.Vec{X: float64(x), Y: float64(y)}, offset)
				worldPos := rend.CellToWorld(screen, cellPos)
				acceleration := r2.Norm(sim.CalculateAccelerationAt(worldPos))

				cornerValues[index] = acceleration
				accelerationSum += acceleration
			}

			forceValues[y*width+x] =
				(accelerationSum - slices.Max(cornerValues[:])) /
					float64(len(offsets)-1)
		}
	}

	maxMass := -math.MaxFloat64
	for _, body := range sim.Bodies {
		maxMass = math.Max(maxMass, body.Mass)
	}

	accelerationMax := simulation.G * maxMass / math.Pow(rend.WorldWidth/100, 2)
	accelerationMin := simulation.G * maxMass / math.Pow(rend.WorldWidth*1.75, 2)

	defaultStyle := tcell.StyleDefault.Foreground(tcell.ColorBlack)

	for y := range height {
		for x := range width {
			value := forceValues[y*width+x]

			relativeValue :=
				(-math.Log(accelerationMin) + math.Log(value)) /
					math.Log(accelerationMax/accelerationMin)

			color := colorMap(relativeValue)
			colorStyle := defaultStyle.Background(color)
			screen.SetContent(x, height-y-1, ' ', nil, colorStyle)
		}
	}
}

func (rend *Renderer) AddFrameMessage(message string) {
	if rend.frameMessage != "" {
		rend.frameMessage += " | "
	}

	rend.frameMessage += message
}

func colorMap(t float64) tcell.Color {
	return tcell.NewRGBColor(
		clamp(int32(2*t*256), 0, 255),
		clamp(int32(2*(1-t)*256), 0, 255),
		0,
	)
}

func (rend *Renderer) writeString(screen tcell.Screen, x, y int, style tcell.Style, str string) {
	width, _ := screen.Size()

	for index, r := range []rune(str) {
		if x+index >= width {
			return
		}

		screen.SetContent(x+index, y, r, nil, style)
	}
}

// @TODO: combine that with computations we do in Render()?
// And maybe use matrix multiplications for that?
func (rend *Renderer) CellToWorld(screen tcell.Screen, cell r2.Vec) r2.Vec {
	width, height := screen.Size()

	worldHeight := (float64(height) / float64(width)) * rend.WorldWidth

	scaleX := float64(width) / rend.WorldWidth

	// Scale y by 0.497 to adjust for non square character terminal (adjusted for my font)
	scaleY := float64(height) / worldHeight * 0.497

	worldX := (cell.X-(float64(width)/2))/scaleX + rend.Center.X
	worldY := (cell.Y-(float64(height)/2))/scaleY + rend.Center.Y

	return r2.Vec{X: worldX, Y: worldY}
}

func makeBraille(partNumber int) rune {
	unicodeOffset := 0
	switch partNumber {
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

func combineBraille(lhs, rhs rune) rune {
	lhsOffset := int(lhs) - 0x2800
	rhsOffset := int(rhs) - 0x2800

	return rune(0x2800 + lhsOffset | rhsOffset)
}

func clamp[T cmp.Ordered](n, a, b T) T {
	if n <= a {
		return a
	} else if n >= b {
		return b
	} else {
		return n
	}
}
