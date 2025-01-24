package systems

import (
	ecs "go-ecs/ecs/components"
	"math"
)

func Distance1D(x1 float64, x2 float64) float64 {
	return math.Pow(x2-x1, 2)
}

func Distance2D(v1 ecs.Vector2, v2 ecs.Vector2) float64 {
	return math.Sqrt(Distance1D(v1.X, v2.X) + Distance1D(v1.Y, v2.Y))
}

func Distance3D(v1 ecs.Vector3, v2 ecs.Vector3) float64 {
	return math.Sqrt(Distance1D(v1.X, v2.X) + Distance1D(v1.Y, v2.Y) + Distance1D(v1.Z, v2.Z))
}

func Distance4D(v1 ecs.Vector4, v2 ecs.Vector4) float64 {
	return math.Sqrt(Distance1D(v1.X, v2.X) + Distance1D(v1.Y, v2.Y) + Distance1D(v1.Z, v2.Z) + Distance1D(v1.W, v2.W))
}

func Distance(v1 interface{}, v2 interface{}) float64 {
	switch v1.(type) {
	case ecs.Vector2:
		return Distance2D(v1.(ecs.Vector2), v2.(ecs.Vector2))
	case ecs.Vector3:
		return Distance3D(v1.(ecs.Vector3), v2.(ecs.Vector3))
	case ecs.Vector4:
		return Distance4D(v1.(ecs.Vector4), v2.(ecs.Vector4))
	}
	return 0
}
