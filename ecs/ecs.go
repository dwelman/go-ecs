package ecs

import (
	"errors"
)

var ErrNotComponent = errors.New("not a component")

// Entity acts as a container of components, it is just an ID
type Entity uint32

type Component struct {
	Type string
	Data any
}

// Manager is a generic type that manages components
type Manager struct {
	// components mapped by type -> [entity -> component]
	components map[string]map[Entity][]Component
	nextID     Entity
}

func NewManager() *Manager {
	return &Manager{
		components: make(map[string]map[Entity][]Component),
		nextID:     0,
	}
}

/** Entity management */

// CreateEntity increments the entity ID counter and returns the next ID in the sequence
func (m *Manager) CreateEntity() Entity {
	m.nextID++
	return m.nextID - 1
}

/** Component management **/

// AddComponent adds a component to an entity
func (m *Manager) AddComponent(entity Entity, component Component) error {
	if _, ok := m.components[component.Type]; !ok {
		m.components[component.Type] = make(map[Entity][]Component)
	}
	m.components[component.Type][entity] = append(m.components[component.Type][entity], component)
	return nil
}

// GetComponent returns the components of the given type on the given entity
func (m *Manager) GetComponent(entity Entity, componentType string) (*[]Component, bool) {
	if _, ok := m.components[componentType]; !ok {
		return nil, false
	}
	c, ok := m.components[componentType][entity]
	return &c, ok
}

// DeleteComponent deletes the component key for the given entity
func (m *Manager) DeleteComponent(entity Entity, componentType string) {
	if _, ok := m.components[componentType]; !ok {
		return
	}
	delete(m.components[componentType], entity)
}
