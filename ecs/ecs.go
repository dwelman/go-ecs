package ecs

// Entity acts as a container of components, it is just an ID
type Entity uint32

type Component struct {
	Type string
}

// TODO: Naive implementation, this will be upgraded as needed
type Manager struct {
	// components mapped by type -> [entity -> component]
	components map[string]map[Entity]Component
	nextID     Entity
}

func NewManager() *Manager {
	return &Manager{
		components: make(map[string]map[Entity]Component),
		nextID:     0,
	}
}

func (m *Manager) CreateEntity() Entity {
	m.nextID++
	return m.nextID - 1
}

/** Component management **/

// AddComponent adds a component to an entity
func (m *Manager) AddComponent(entity Entity, component Component) {
	if _, ok := m.components[component.Type]; !ok {
		m.components[component.Type] = make(map[Entity]Component)
	}
	m.components[component.Type][entity] = component
}

// GetComponent returns a component of an entity
func (m *Manager) GetComponent(entity Entity, componentType string) (Component, bool) {
	if _, ok := m.components[componentType]; !ok {
		return Component{}, false
	}
	component, ok := m.components[componentType][entity]
	return component, ok
}

// DeleteComponent deletes a component from an entity
func (m *Manager) DeleteComponent(entity Entity, componentType string) {
	if _, ok := m.components[componentType]; !ok {
		return
	}
	delete(m.components[componentType], entity)
}

// TODO: DeleteEntity
