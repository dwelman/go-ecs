package components

import (
	"go-ecs/ecs"

	"gonum.org/v1/gonum/spatial/r2"
	"gonum.org/v1/gonum/spatial/r3"
)

func Vector2(x float64, y float64) ecs.Component {
	return ecs.Component{
		Type: "Vector2",
		Data: r2.Vec{X: x, Y: y},
	}
}

func Vector3(x float64, y float64, z float64) ecs.Component {
	return ecs.Component{
		Type: "Vector3",
		Data: r3.Vec{X: x, Y: y, Z: z},
	}
}
