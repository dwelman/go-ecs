package systems

import (
	"testing"

	"go-ecs/ecs"

	"github.com/stretchr/testify/require"
	"gonum.org/v1/gonum/spatial/r2"
)

func Test_MovementSystem(t *testing.T) {
	m := ecs.NewManager()

	// Create test entities with Vector2 and Velocity2D components
	entity := m.CreateEntity()
	m.AddComponentToEntity(entity, ecs.Component{
		Type: "Vector2",
		Data: r2.Vec{X: 0, Y: 0},
	})
	m.AddComponentToEntity(entity, ecs.Component{
		Type: "Velocity2D",
		Data: r2.Vec{X: 1, Y: 1},
	})

	// Run the MovementSystem
	deltaT := 1.0
	require.NoError(t, MovementSystem(m, deltaT))

	// Check the updated position
	vector, err := m.GetComponentOfEntity(entity, "Vector2")
	require.NoError(t, err)
	expectedVector := r2.Vec{X: 1, Y: 1}
	actualVector, err := ecs.GetDataAsType[r2.Vec](vector)
	require.NoError(t, err)
	require.Equal(t, expectedVector, actualVector)
}

func Test_MovementSystem_NoVelocity(t *testing.T) {
	m := ecs.NewManager()

	// Create test entities with Vector2 and Velocity2D components
	entity := m.CreateEntity()
	m.AddComponentToEntity(entity, ecs.Component{
		Type: "Vector2",
		Data: r2.Vec{X: 0, Y: 0},
	})

	// Run the MovementSystem
	deltaT := 1.0
	err := MovementSystem(m, deltaT)
	require.ErrorIs(t, err, ecs.ErrComponentTypeNotFound)
}

func Test_MovementSystem_MoveInSequence(t *testing.T) {
	type testSuite struct {
		name                 string
		startPositions       []r2.Vec
		velocityChanges      []r2.Vec
		expectedEndPositions []r2.Vec
	}

	tests := []testSuite{
		{
			name: "Straight line",
			startPositions: []r2.Vec{
				{X: 0, Y: 0},
				{X: -1, Y: 1},
			},
			velocityChanges: []r2.Vec{
				{X: 1, Y: 0},
				{X: 1, Y: 0},
				{X: 1, Y: 0},
				{X: 1, Y: 0},
				{X: 1, Y: 0},
			},
			expectedEndPositions: []r2.Vec{
				{X: 5, Y: 0},
				{X: 4, Y: 1},
			},
		},
		{
			name: "Diagonal",
			startPositions: []r2.Vec{
				{X: 0, Y: 0},
				{X: 1, Y: 1},
			},
			velocityChanges: []r2.Vec{
				{X: 1, Y: 1},
				{X: 1, Y: 1},
				{X: 1, Y: 1},
				{X: 1, Y: 1},
				{X: 1, Y: 1},
			},
			expectedEndPositions: []r2.Vec{
				{X: 5, Y: 5},
				{X: 6, Y: 6},
			},
		},
		{
			name: "Zigzag",
			startPositions: []r2.Vec{
				{X: 0, Y: 0},
				{X: -1, Y: 1},
			},
			velocityChanges: []r2.Vec{
				{X: 1, Y: 0},
				{X: 0, Y: 1},
				{X: 1, Y: 0},
				{X: 0, Y: 1},
				{X: 1, Y: 0},
			},
			expectedEndPositions: []r2.Vec{
				{X: 3, Y: 2},
				{X: 2, Y: 3},
			},
		},
		{
			name: "Negative",
			startPositions: []r2.Vec{
				{X: 0, Y: 0},
				{X: -1, Y: 0},
			},
			velocityChanges: []r2.Vec{
				{X: -1, Y: 0},
				{X: -1, Y: 0},
				{X: -1, Y: 0},
				{X: -1, Y: 0},
				{X: -1, Y: 0},
			},
			expectedEndPositions: []r2.Vec{
				{X: -5, Y: 0},
				{X: -6, Y: 0},
			},
		},
		{
			name: "Mixed",
			startPositions: []r2.Vec{
				{X: 0, Y: 0},
				{X: -1, Y: 2},
			},
			velocityChanges: []r2.Vec{
				{X: 1, Y: 0},
				{X: 0, Y: 1},
				{X: -1, Y: 0},
				{X: 0, Y: -1},
				{X: 1, Y: 0},
			},
			expectedEndPositions: []r2.Vec{
				{X: 1, Y: 0},
				{X: 0, Y: 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := ecs.NewManager()

			entities := make([]ecs.Entity, len(tt.startPositions))
			for i, startPosition := range tt.startPositions {
				entity := m.CreateEntity()
				m.AddComponentToEntity(entity, ecs.Component{
					Type: "Vector2",
					Data: startPosition,
				})
				entities[i] = entity
			}

			for _, velocity := range tt.velocityChanges {
				for _, entity := range entities {
					m.AddComponentToEntity(entity, ecs.Component{
						Type: "Velocity2D",
						Data: velocity,
					})
				}

				deltaT := 1.0
				require.NoError(t, MovementSystem(m, deltaT))
			}

			for i, entity := range entities {
				vector, err := m.GetComponentOfEntity(entity, "Vector2")
				require.NoError(t, err)

				expectedVector := tt.expectedEndPositions[i]
				actualVector, err := ecs.GetDataAsType[r2.Vec](vector)
				require.NoError(t, err)
				require.Equal(t, expectedVector, actualVector)
			}
		})
	}
}
