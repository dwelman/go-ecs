package ecs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestComponentString struct {
	content string
}

const TestComponentStringKey = "TestComponentString"

type TestComponentNumber struct {
	content int
}

const TestComponentNumberKey = "TestComponentNumber"

func TestManager_Component_CRUD(t *testing.T) {
	m := NewManager()
	t.Log("Add component to entity - succeeds")
	{
		e := m.CreateEntity()
		require.NoError(t, m.AddComponentToEntity(e, Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}}))
		t.Log("Get component of type from entity - succeeds")
		{
			components, err := m.GetComponentsOfEntity(e, TestComponentStringKey)
			require.NoError(t, err)
			require.NotNil(t, components)
			require.Len(t, *components, 1)
			require.Equal(t, "Hello", (*components)[0].Data.(TestComponentString).content)
		}
	}
}
func Test_CreateEntity(t *testing.T) {
	m := NewManager()
	t.Log("Creating an entity increments the ID counter")
	{
		for i := 0; i < 10; i++ {
			e := m.CreateEntity()
			require.Equal(t, Entity(i), e)
		}
	}
}
func Test_AddComponentToEntity(t *testing.T) {
	m := NewManager()
	t.Log("Add component to entity - succeeds")
	{
		e := m.CreateEntity()
		err := m.AddComponentToEntity(e, Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}})
		require.NoError(t, err)

		components, err := m.GetComponentsOfEntity(e, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, components)
		require.Len(t, *components, 1)
		require.Equal(t, "Hello", (*components)[0].Data.(TestComponentString).content)
	}

	t.Log("Add multiple components of the same type to entity - succeeds")
	{
		e := m.CreateEntity()
		err := m.AddComponentToEntity(e, Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}})
		require.NoError(t, err)
		err = m.AddComponentToEntity(e, Component{Type: TestComponentStringKey, Data: TestComponentString{content: "World"}})
		require.NoError(t, err)

		components, err := m.GetComponentsOfEntity(e, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, components)
		require.Len(t, *components, 2)
		require.Equal(t, "Hello", (*components)[0].Data.(TestComponentString).content)
		require.Equal(t, "World", (*components)[1].Data.(TestComponentString).content)
	}

	t.Log("Add component to multiple entities - succeeds")
	{
		e1 := m.CreateEntity()
		e2 := m.CreateEntity()
		err := m.AddComponentToEntity(e1, Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}})
		require.NoError(t, err)
		err = m.AddComponentToEntity(e2, Component{Type: TestComponentStringKey, Data: TestComponentString{content: "World"}})
		require.NoError(t, err)

		components1, err := m.GetComponentsOfEntity(e1, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, components1)
		require.Len(t, *components1, 1)
		require.Equal(t, "Hello", (*components1)[0].Data.(TestComponentString).content)

		components2, err := m.GetComponentsOfEntity(e2, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, components2)
		require.Len(t, *components2, 1)
		require.Equal(t, "World", (*components2)[0].Data.(TestComponentString).content)
	}

	t.Log("Add component to entity of different type - succeeds")
	{
		e := m.CreateEntity()
		err := m.AddComponentToEntity(e, Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 42}})
		require.NoError(t, err)
		err = m.AddComponentToEntity(e, Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}})
		require.NoError(t, err)

		componentsNumber, err := m.GetComponentsOfEntity(e, TestComponentNumberKey)
		require.NoError(t, err)
		require.NotNil(t, componentsNumber)
		require.Len(t, *componentsNumber, 1)
		require.Equal(t, 42, (*componentsNumber)[0].Data.(TestComponentNumber).content)

		componentsString, err := m.GetComponentsOfEntity(e, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, componentsString)
		require.Len(t, *componentsString, 1)
		require.Equal(t, "Hello", (*componentsString)[0].Data.(TestComponentString).content)
	}
}

func Test_GetComponentsOfEntity(t *testing.T) {

	t.Log("Get components of non-existent type - fails")
	{
		m := NewManager()
		components, err := m.GetComponentsOfEntity(0, "NonExistentType")
		require.ErrorIs(t, ErrComponentTypeNotFound, err)
		require.Nil(t, components)
	}

	t.Log("Get components of non-existent entity - fails")
	{
		m := Manager{
			components: map[string]map[Entity][]Component{
				TestComponentStringKey: {
					0: {Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}}},
				},
			},
			nextID: 1,
		}
		components, err := m.GetComponentsOfEntity(999, TestComponentStringKey)
		require.ErrorIs(t, ErrComponentNotFound, err)
		require.Nil(t, components)
	}

	t.Log("Get components of existing type and entity - succeeds")
	{
		m := Manager{
			components: map[string]map[Entity][]Component{
				TestComponentStringKey: {
					0: {Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}}},
				},
			},
			nextID: 1,
		}
		components, err := m.GetComponentsOfEntity(0, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, components)
		require.Len(t, *components, 1)
		require.Equal(t, "Hello", (*components)[0].Data.(TestComponentString).content)
	}

	t.Log("Get components of multiple types and entities - succeeds")
	{
		m := Manager{
			components: map[string]map[Entity][]Component{
				TestComponentStringKey: {
					0: {Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}}},
				},
				TestComponentNumberKey: {
					0: {Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 42}}},
				},
			},
			nextID: 1,
		}
		componentsString, err := m.GetComponentsOfEntity(0, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, componentsString)
		require.Len(t, *componentsString, 1)
		require.Equal(t, "Hello", (*componentsString)[0].Data.(TestComponentString).content)

		componentsNumber, err := m.GetComponentsOfEntity(0, TestComponentNumberKey)
		require.NoError(t, err)
		require.NotNil(t, componentsNumber)
		require.Len(t, *componentsNumber, 1)
		require.Equal(t, 42, (*componentsNumber)[0].Data.(TestComponentNumber).content)
	}

	t.Log("Get components of multiple types and entities - succeeds")
	{
		m := Manager{
			components: map[string]map[Entity][]Component{
				TestComponentStringKey: {
					0: {Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}}},
					1: {Component{Type: TestComponentStringKey, Data: TestComponentString{content: "World"}}},
				},
				TestComponentNumberKey: {
					0: {Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 42}}},
					1: {Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 43}}},
				},
			},
			nextID: 2,
		}

		componentsString0, err := m.GetComponentsOfEntity(0, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, componentsString0)
		require.Len(t, *componentsString0, 1)
		require.Equal(t, "Hello", (*componentsString0)[0].Data.(TestComponentString).content)

		componentsString1, err := m.GetComponentsOfEntity(1, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, componentsString1)
		require.Len(t, *componentsString1, 1)
		require.Equal(t, "World", (*componentsString1)[0].Data.(TestComponentString).content)

		componentsNumber0, err := m.GetComponentsOfEntity(0, TestComponentNumberKey)
		require.NoError(t, err)
		require.NotNil(t, componentsNumber0)
		require.Len(t, *componentsNumber0, 1)
		require.Equal(t, 42, (*componentsNumber0)[0].Data.(TestComponentNumber).content)

		componentsNumber1, err := m.GetComponentsOfEntity(1, TestComponentNumberKey)
		require.NoError(t, err)
		require.NotNil(t, componentsNumber1)
		require.Len(t, *componentsNumber1, 1)
	}
}

func Test_DeleteComponentsOfEntity(t *testing.T) {
	t.Log("Delete components of non-existent type - fails")
	{
		m := NewManager()
		require.ErrorIs(t, m.DeleteComponentsOfEntity(0, "NonExistentType"), ErrComponentTypeNotFound)
		components, err := m.GetComponentsOfEntity(0, "NonExistentType")
		require.ErrorIs(t, err, ErrComponentTypeNotFound)
		require.Nil(t, components)
	}

	t.Log("Delete components of non-existent entity - no-op")
	{
		m := Manager{
			components: map[string]map[Entity][]Component{
				TestComponentStringKey: {
					0: {Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}}},
				},
			},
			nextID: 1,
		}
		m.DeleteComponentsOfEntity(999, TestComponentStringKey)
		components, err := m.GetComponentsOfEntity(0, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, components)
		require.Len(t, *components, 1)
		require.Equal(t, "Hello", (*components)[0].Data.(TestComponentString).content)
	}

	t.Log("Delete components of existing type and entity - succeeds")
	{
		m := Manager{
			components: map[string]map[Entity][]Component{
				TestComponentStringKey: {
					0: {Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}}},
				},
			},
			nextID: 1,
		}

		require.NoError(t, m.DeleteComponentsOfEntity(0, TestComponentStringKey))
		components, err := m.GetComponentsOfEntity(0, TestComponentStringKey)
		require.Error(t, err)
		require.Nil(t, components)
	}

	t.Log("Delete components of multiple types and entities - succeeds")
	{

		m := Manager{
			components: map[string]map[Entity][]Component{
				TestComponentStringKey: {
					0: {Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}}},
					1: {Component{Type: TestComponentStringKey, Data: TestComponentString{content: "World"}}},
				},
				TestComponentNumberKey: {
					0: {Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 42}}},
					1: {Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 43}}},
				},
			},
			nextID: 2,
		}

		require.NoError(t, m.DeleteComponentsOfEntity(0, TestComponentStringKey))
		componentsString0, err := m.GetComponentsOfEntity(0, TestComponentStringKey)
		require.ErrorIs(t, err, ErrComponentNotFound)
		require.Nil(t, componentsString0)

		require.NoError(t, m.DeleteComponentsOfEntity(1, TestComponentNumberKey))
		componentsNumber1, err := m.GetComponentsOfEntity(1, TestComponentNumberKey)
		require.ErrorIs(t, err, ErrComponentNotFound)
		require.Nil(t, componentsNumber1)

		componentsString1, err := m.GetComponentsOfEntity(1, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, componentsString1)

		componentsNumber0, err := m.GetComponentsOfEntity(0, TestComponentNumberKey)
		require.NoError(t, err)
		require.NotNil(t, componentsNumber0)
	}
}
