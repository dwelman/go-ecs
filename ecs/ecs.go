package ecs

import (
	"errors"
)

var ErrComponentTypeNotFound = errors.New("manager does not have components of this type")
var ErrComponentNotFound = errors.New("component not found")
var ErrEntityNotFound = errors.New("entity not found")

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
func (m *Manager) AddComponentToEntity(entity Entity, component Component) error {
	if _, ok := m.components[component.Type]; !ok {
		m.components[component.Type] = make(map[Entity][]Component)
	}
	m.components[component.Type][entity] = append(m.components[component.Type][entity], component)
	return nil
}

// GetComponentsOfEntity returns the components of the given type on the given entity
func (m *Manager) GetComponentsOfEntity(entity Entity, componentType string) (*[]Component, error) {
	if _, ok := m.components[componentType]; !ok {
		return nil, ErrComponentTypeNotFound
	}
	c, ok := m.components[componentType][entity]
	if !ok {
		return nil, ErrComponentNotFound
	}
	return &c, nil
}

// DeleteComponent deletes the component key for the given entity
func (m *Manager) DeleteComponentsOfEntity(entity Entity, componentType string) error {
	if _, ok := m.components[componentType]; !ok {
		return ErrComponentTypeNotFound
	}
	delete(m.components[componentType], entity)
	return nil
}
