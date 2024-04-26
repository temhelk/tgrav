package simulation

import (
	"math"

	"gonum.org/v1/gonum/spatial/r2"
)

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

var GravitySlingshot = [...]Body{
	{
		Mass:     1e13,
		Position: r2.Vec{X: 0, Y: 0},
		Velocity: r2.Vec{X: 0, Y: 0},
	},
	{
		Mass:     1e12,
		Position: r2.Vec{X: 15, Y: 0},
		Velocity: r2.Vec{X: 0, Y: -7},
	},
	{
		Mass:     1e11,
		Position: r2.Vec{X: -15, Y: 0},
		Velocity: r2.Vec{X: 0, Y: 5},
	},
}

var LagrangeL4L5 = [...]Body{
	{
		Mass:     1e12,
		Position: r2.Vec{X: 0, Y: 0},
		Velocity: r2.Vec{X: 0, Y: 0},
	},
	{
		Mass:     1e10,
		Position: r2.Vec{X: 20, Y: 0},
		Velocity: r2.Vec{X: 0, Y: -1.836},
	},
	{
		Mass:     1,
		Position: r2.Vec{X: 10, Y: 17.32},
		Velocity: r2.Vec{X: 0.965*1.836*math.Cos(30.0*math.Pi/180), Y: -1.836*math.Sin(30*math.Pi/180)},
	},
}

var EarthMoon = [...]Body{
	{
		Mass:     5.972e24,
		Position: r2.Vec{X: 0, Y: 0},
		Velocity: r2.Vec{X: 0, Y: 0},
	},
	{
		Mass:     7.347e22,
		Position: r2.Vec{X: 384.4e6, Y: 0},
		Velocity: r2.Vec{X: 0, Y: -1022.0},
	},
	{
		Mass:     1,
		Position: r2.Vec{X: 1.922e8, Y: 3.329e8},
		Velocity: r2.Vec{X: 1022.0*math.Cos(30.0*math.Pi/180)*0.98, Y: -1022.0*math.Sin(30*math.Pi/180)},
	},
}
