package systems

import (
	"go-ecs/ecs"

	"gonum.org/v1/gonum/spatial/r2"
)

// TODO: Better error handling, don't crash loop for error
func MovementSystem(m *ecs.Manager, deltaT float64) error {
	entities, err := m.GetEntitiesWithComponents([]string{"Vector2", "Velocity2D"})

	if err != nil {
		return err
	}

	for e, c := range entities {
		vector, err := ecs.GetComponentData[r2.Vec](c, "Vector2")
		if err != nil {
			return err
		}

		velocity, err := ecs.GetComponentData[r2.Vec](c, "Velocity2D")
		if err != nil {
			return err
		}

		vector.X += velocity.X * deltaT
		vector.Y += velocity.Y * deltaT

		m.AddComponentToEntity(e, ecs.Component{
			Type: "Vector2",
			Data: *vector,
		})
	}

	return nil
}
