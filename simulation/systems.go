package simulation

import "gonum.org/v1/gonum/spatial/r2"

var FourBodySystem = [...]Body{
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

var ThreeBodyUnstableSystem = [...]Body{
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
