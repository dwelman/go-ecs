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
			c, err := m.GetComponentOfEntity(e, TestComponentStringKey)
			require.NoError(t, err)
			require.NotNil(t, c)
			require.Equal(t, "Hello", (*c).Data.(TestComponentString).content)
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

func Test_DeleteEntity(t *testing.T) {
	t.Log("Delete non-existent entity - fails")
	{
		m := NewManager()
		require.ErrorIs(t, m.DeleteEntity(0), ErrEntityNotFound)
	}
	t.Log("Delete existing entity - succeeds")
	{
		m := Manager{
			components: map[string]map[Entity]*Component{
				TestComponentStringKey: {
					0: &Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}},
					1: &Component{Type: TestComponentStringKey, Data: TestComponentString{content: "World"}},
				},
			},
		}
		require.NoError(t, m.DeleteEntity(0))
		c, err := m.GetComponentOfEntity(0, TestComponentStringKey)
		require.ErrorIs(t, err, ErrComponentNotFound)
		require.Nil(t, c)

		c, err = m.GetComponentOfEntity(1, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, c)
		require.Equal(t, "World", (*c).Data.(TestComponentString).content)

		t.Log("Create new entity after deleting existing entity - reuses ID")
		{
			e := m.CreateEntity()
			require.Equal(t, Entity(0), e)
			require.Len(t, m.freeIDs, 0)
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

		c, err := m.GetComponentOfEntity(e, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, c)
		require.Equal(t, "Hello", (*c).Data.(TestComponentString).content)
	}

	t.Log("Add component to multiple entities - succeeds")
	{
		e1 := m.CreateEntity()
		e2 := m.CreateEntity()
		err := m.AddComponentToEntity(e1, Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}})
		require.NoError(t, err)
		err = m.AddComponentToEntity(e2, Component{Type: TestComponentStringKey, Data: TestComponentString{content: "World"}})
		require.NoError(t, err)

		c1, err := m.GetComponentOfEntity(e1, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, c1)
		require.Equal(t, "Hello", (*c1).Data.(TestComponentString).content)

		c2, err := m.GetComponentOfEntity(e2, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, c2)
		require.Equal(t, "World", (*c2).Data.(TestComponentString).content)
	}

	t.Log("Add component to entity of different type - succeeds")
	{
		e := m.CreateEntity()
		err := m.AddComponentToEntity(e, Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 42}})
		require.NoError(t, err)
		err = m.AddComponentToEntity(e, Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}})
		require.NoError(t, err)

		cNumber, err := m.GetComponentOfEntity(e, TestComponentNumberKey)
		require.NoError(t, err)
		require.NotNil(t, cNumber)
		require.Equal(t, 42, (*cNumber).Data.(TestComponentNumber).content)

		cString, err := m.GetComponentOfEntity(e, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, cString)
		require.Equal(t, "Hello", (*cString).Data.(TestComponentString).content)
	}
}

func Test_GetComponentsOfEntity(t *testing.T) {

	t.Log("Get component of non-existent type - fails")
	{
		m := NewManager()
		c, err := m.GetComponentOfEntity(0, "NonExistentType")
		require.ErrorIs(t, ErrComponentTypeNotFound, err)
		require.Nil(t, c)
	}

	t.Log("Get components of non-existent entity - fails")
	{
		m := Manager{
			components: map[string]map[Entity]*Component{
				TestComponentStringKey: {
					0: &Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}},
				},
			},
			nextID: 1,
		}
		c, err := m.GetComponentOfEntity(999, TestComponentStringKey)
		require.ErrorIs(t, ErrComponentNotFound, err)
		require.Nil(t, c)
	}

	t.Log("Get components of existing type and entity - succeeds")
	{
		m := Manager{
			components: map[string]map[Entity]*Component{
				TestComponentStringKey: {
					0: &Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}},
				},
			},
			nextID: 1,
		}
		c, err := m.GetComponentOfEntity(0, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, c)
		require.Equal(t, "Hello", (*c).Data.(TestComponentString).content)
	}

	t.Log("Get components of multiple types and entities - succeeds")
	{
		m := Manager{
			components: map[string]map[Entity]*Component{
				TestComponentStringKey: {
					0: &Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}},
				},
				TestComponentNumberKey: {
					0: &Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 42}},
				},
			},
			nextID: 1,
		}
		cString, err := m.GetComponentOfEntity(0, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, cString)
		require.Equal(t, "Hello", (*cString).Data.(TestComponentString).content)

		cNumber, err := m.GetComponentOfEntity(0, TestComponentNumberKey)
		require.NoError(t, err)
		require.NotNil(t, cNumber)
		require.Equal(t, 42, (*cNumber).Data.(TestComponentNumber).content)
	}

	t.Log("Get components of multiple types and entities - succeeds")
	{
		m := Manager{
			components: map[string]map[Entity]*Component{
				TestComponentStringKey: {
					0: &Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}},
					1: &Component{Type: TestComponentStringKey, Data: TestComponentString{content: "World"}},
				},
				TestComponentNumberKey: {
					0: &Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 42}},
					1: &Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 43}},
				},
			},
			nextID: 2,
		}

		cString0, err := m.GetComponentOfEntity(0, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, cString0)
		require.Equal(t, "Hello", (*cString0).Data.(TestComponentString).content)

		cString1, err := m.GetComponentOfEntity(1, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, cString1)
		require.Equal(t, "World", (*cString1).Data.(TestComponentString).content)

		cNumber0, err := m.GetComponentOfEntity(0, TestComponentNumberKey)
		require.NoError(t, err)
		require.NotNil(t, cNumber0)
		require.Equal(t, 42, (*cNumber0).Data.(TestComponentNumber).content)

		cNumber1, err := m.GetComponentOfEntity(1, TestComponentNumberKey)
		require.NoError(t, err)
		require.NotNil(t, cNumber1)
	}
}

func Test_DeleteComponentsOfEntity(t *testing.T) {
	t.Log("Delete components of non-existent type - fails")
	{
		m := NewManager()
		require.ErrorIs(t, m.DeleteComponentOfEntity(0, "NonExistentType"), ErrComponentTypeNotFound)
		c, err := m.GetComponentOfEntity(0, "NonExistentType")
		require.ErrorIs(t, err, ErrComponentTypeNotFound)
		require.Nil(t, c)
	}

	t.Log("Delete components of non-existent entity - no-op")
	{
		m := Manager{
			components: map[string]map[Entity]*Component{
				TestComponentStringKey: {
					0: &Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}},
				},
			},
			nextID: 1,
		}
		m.DeleteComponentOfEntity(999, TestComponentStringKey)
		c, err := m.GetComponentOfEntity(0, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, c)
		require.Equal(t, "Hello", (*c).Data.(TestComponentString).content)
	}

	t.Log("Delete components of existing type and entity - succeeds")
	{
		m := Manager{
			components: map[string]map[Entity]*Component{
				TestComponentStringKey: {
					0: &Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}},
				},
			},
			nextID: 1,
		}

		require.NoError(t, m.DeleteComponentOfEntity(0, TestComponentStringKey))
		c, err := m.GetComponentOfEntity(0, TestComponentStringKey)
		require.Error(t, err)
		require.Nil(t, c)
	}

	t.Log("Delete components of multiple types and entities - succeeds")
	{

		m := Manager{
			components: map[string]map[Entity]*Component{
				TestComponentStringKey: {
					0: &Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}},
					1: &Component{Type: TestComponentStringKey, Data: TestComponentString{content: "World"}},
				},
				TestComponentNumberKey: {
					0: &Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 42}},
					1: &Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 43}},
				},
			},
			nextID: 2,
		}

		require.NoError(t, m.DeleteComponentOfEntity(0, TestComponentStringKey))
		cString0, err := m.GetComponentOfEntity(0, TestComponentStringKey)
		require.ErrorIs(t, err, ErrComponentNotFound)
		require.Nil(t, cString0)

		require.NoError(t, m.DeleteComponentOfEntity(1, TestComponentNumberKey))
		cNumber1, err := m.GetComponentOfEntity(1, TestComponentNumberKey)
		require.ErrorIs(t, err, ErrComponentNotFound)
		require.Nil(t, cNumber1)

		cString1, err := m.GetComponentOfEntity(1, TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, cString1)

		cNumber0, err := m.GetComponentOfEntity(0, TestComponentNumberKey)
		require.NoError(t, err)
		require.NotNil(t, cNumber0)
	}
}

func Test_GetEntitiesWithComponents(t *testing.T) {
	t.Log("Get components of non-existent type - fails")
	{
		m := NewManager()
		c, err := m.GetEntitiesWithComponents([]string{"NonExistentType"})
		require.ErrorIs(t, err, ErrComponentTypeNotFound)
		require.Nil(t, c)
	}
	t.Log("Get components of multiple types - succeeds")
	{
		m := Manager{
			components: map[string]map[Entity]*Component{
				TestComponentStringKey: {
					0: &Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}},
					1: &Component{Type: TestComponentStringKey, Data: TestComponentString{content: "World"}},
					3: &Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Foo"}},
				},
				TestComponentNumberKey: {
					0: &Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 42}},
					1: &Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 43}},
					2: &Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 44}},
				},
			},
		}

		ec, err := m.GetEntitiesWithComponents([]string{TestComponentStringKey, TestComponentNumberKey})
		require.NoError(t, err)
		require.NotNil(t, ec)
		require.Len(t, ec, 2)
		require.Len(t, ec[0], 2)
		require.Len(t, ec[1], 2)
		_, ok := ec[2]
		require.False(t, ok)

		require.Equal(t, "Hello", ec[0][0].Data.(TestComponentString).content)
		require.Equal(t, 42, ec[0][1].Data.(TestComponentNumber).content)
		require.Equal(t, "World", ec[1][0].Data.(TestComponentString).content)
		require.Equal(t, 43, ec[1][1].Data.(TestComponentNumber).content)
	}

}
func Test_GetComponentData(t *testing.T) {
	t.Log("Get component data of existing type - succeeds")
	{
		m := Manager{
			components: map[string]map[Entity]*Component{
				TestComponentStringKey: {
					0: &Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}},
					1: &Component{Type: TestComponentStringKey, Data: TestComponentString{content: "World"}},
				},
				TestComponentNumberKey: {
					0: &Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 42}},
					1: &Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 43}},
				},
			},
			nextID: 1,
		}

		entitiesWithComponents, err := m.GetEntitiesWithComponents([]string{TestComponentStringKey, TestComponentNumberKey})
		require.NoError(t, err)
		require.NotNil(t, entitiesWithComponents)

		data, err := GetComponentData[TestComponentString](entitiesWithComponents[0], TestComponentStringKey)
		require.NoError(t, err)
		require.NotNil(t, data)
		require.Equal(t, "Hello", data.content)

		numberData, err := GetComponentData[TestComponentNumber](entitiesWithComponents[0], TestComponentNumberKey)
		require.NoError(t, err)
		require.NotNil(t, numberData)
		require.Equal(t, 42, numberData.content)
	}

	t.Log("Get component data of non-existent type - fails")
	{
		m := Manager{
			components: map[string]map[Entity]*Component{
				TestComponentStringKey: {
					0: &Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}},
					1: &Component{Type: TestComponentStringKey, Data: TestComponentString{content: "World"}},
				},
				TestComponentNumberKey: {
					0: &Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 42}},
					1: &Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 43}},
				},
			},
			nextID: 1,
		}

		entitiesWithComponents, err := m.GetEntitiesWithComponents([]string{TestComponentStringKey, TestComponentNumberKey})
		require.NoError(t, err)
		require.NotNil(t, entitiesWithComponents)

		data, err := GetComponentData[TestComponentString](entitiesWithComponents[0], "NonExistentType")
		require.ErrorIs(t, err, ErrComponentNotFound)
		require.Nil(t, data)
	}

	t.Log("Get component data with type mismatch - fails")
	{
		m := Manager{
			components: map[string]map[Entity]*Component{
				TestComponentStringKey: {
					0: &Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}},
					1: &Component{Type: TestComponentStringKey, Data: TestComponentString{content: "World"}},
				},
				TestComponentNumberKey: {
					0: &Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 42}},
					1: &Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 43}},
				},
			},
			nextID: 1,
		}

		entitiesWithComponents, err := m.GetEntitiesWithComponents([]string{TestComponentStringKey, TestComponentNumberKey})
		require.NoError(t, err)
		require.NotNil(t, entitiesWithComponents)

		data, err := GetComponentData[TestComponentString](entitiesWithComponents[0], TestComponentNumberKey)
		require.ErrorIs(t, err, ErrComponentDataMismatch)
		require.Nil(t, data)
	}
}
func Test_GetDataAsType(t *testing.T) {
	t.Log("Get data as correct type - succeeds")
	{
		component := Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}}
		data, err := GetDataAsType[TestComponentString](&component)
		require.NoError(t, err)
		require.Equal(t, "Hello", data.content)
	}

	t.Log("Get data as incorrect type - fails")
	{
		component := Component{Type: TestComponentStringKey, Data: TestComponentString{content: "Hello"}}
		data, err := GetDataAsType[TestComponentNumber](&component)
		require.ErrorIs(t, err, ErrComponentDataMismatch)
		require.Equal(t, TestComponentNumber{}, data)
	}

	t.Log("Get data as correct type with number - succeeds")
	{
		component := Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 42}}
		data, err := GetDataAsType[TestComponentNumber](&component)
		require.NoError(t, err)
		require.Equal(t, 42, data.content)
	}

	t.Log("Get data as incorrect type with number - fails")
	{
		component := Component{Type: TestComponentNumberKey, Data: TestComponentNumber{content: 42}}
		data, err := GetDataAsType[TestComponentString](&component)
		require.ErrorIs(t, err, ErrComponentDataMismatch)
		require.Equal(t, TestComponentString{}, data)
	}

	t.Log("Get data as pointer type - succeeds")
	{
		component := Component{Type: TestComponentStringKey, Data: &TestComponentString{content: "Hello"}}
		data, err := GetDataAsType[*TestComponentString](&component)
		require.NoError(t, err)
		require.Equal(t, "Hello", data.content)
	}
}
