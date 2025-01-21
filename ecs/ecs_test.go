package ecs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestComponent struct {
	Component
	message string
}

func TestManager_Entity_CRUD(t *testing.T) {
	m := NewManager()
	t.Log("Creating an entity increments the ID counter")
	{
		for i := 0; i < 10; i++ {
			e := m.CreateEntity()
			require.Equal(t, Entity(i), e)
		}
	}
}

func TestManager_Component_CRUD(t *testing.T) {
	m := NewManager()
	t.Log("Add component to entity - succeeds")
	{
		e := m.CreateEntity()
		m.AddComponent(e, TestComponent{
			Component: Component{
				Type: "TestComponent",
			},
			message: "foo",
		})
		t.Log("Get component of type from entity - succeeds")
		{
			components, ok := m.GetComponent(e, "TestComponent")
			require.True(t, ok)
			require.NotNil(t, components)
		}
	}
}
