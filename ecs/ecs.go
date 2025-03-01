package ecs

import (
	"errors"
)

var ErrComponentTypeNotFound = errors.New("manager does not have components of this type")
var ErrComponentNotFound = errors.New("component not found")
var ErrEntityNotFound = errors.New("entity not found")
var ErrComponentDataMismatch = errors.New("component data type mismatch")

// Entity acts as a container of components, it is just an ID
type Entity uint32

type Component struct {
	Type string
	Data any
}

func GetDataAsType[T any](c *Component) (T, error) {
	data, ok := c.Data.(T)
	if !ok {
		var zero T
		return zero, ErrComponentDataMismatch
	}
	return data, nil
}

// Manager is a generic type that manages components
type Manager struct {
	// components mapped by type -> [entity -> component]
	components map[string]map[Entity]*Component
	nextID     Entity
	// freeIDs is a list of IDs that have been deleted and can be reused
	freeIDs []Entity
}

func NewManager() *Manager {
	return &Manager{
		components: make(map[string]map[Entity]*Component),
		nextID:     0,
		freeIDs:    make([]Entity, 0),
	}
}

/** Entity management */

// CreateEntity increments the entity ID counter and returns the next ID in the sequence
func (m *Manager) CreateEntity() Entity {
	if len(m.freeIDs) > 0 {
		id := m.freeIDs[0]
		m.freeIDs = m.freeIDs[1:]
		return id
	}
	m.nextID++
	return m.nextID - 1
}

// DeleteEntity deletes the entity and all its components, and stores the ID in the freeIDs list
func (m *Manager) DeleteEntity(entity Entity) error {
	deleted := false
	for _, components := range m.components {
		if _, ok := components[entity]; !ok {
			continue
		}
		delete(components, entity)
		deleted = true
	}

	if !deleted {
		return ErrEntityNotFound
	}
	m.freeIDs = append(m.freeIDs, entity)
	return nil
}

/** Component management **/

// AddComponentToEntity adds a component to an entity
func (m *Manager) AddComponentToEntity(entity Entity, component Component) error {
	if _, ok := m.components[component.Type]; !ok {
		m.components[component.Type] = make(map[Entity]*Component)
	}
	m.components[component.Type][entity] = &component
	return nil
}

// GetComponentOfEntity returns the component of the given type on the given entity
func (m *Manager) GetComponentOfEntity(entity Entity, componentType string) (*Component, error) {
	if _, ok := m.components[componentType]; !ok {
		return nil, ErrComponentTypeNotFound
	}
	c, ok := m.components[componentType][entity]
	if !ok {
		return nil, ErrComponentNotFound
	}
	return c, nil
}

func (m *Manager) DeleteComponentOfEntity(entity Entity, componentType string) error {
	if _, ok := m.components[componentType]; !ok {
		return ErrComponentTypeNotFound
	}
	if _, ok := m.components[componentType][entity]; !ok {
		return ErrComponentNotFound
	}
	delete(m.components[componentType], entity)
	return nil
}

// GetEntitiesWithComponents returns entities and components where the entity has all types of components
func (m *Manager) GetEntitiesWithComponents(types []string) (map[Entity][]*Component, error) {
	result := make(map[Entity][]*Component)

	for _, t := range types {
		if _, ok := m.components[t]; !ok {
			return nil, ErrComponentTypeNotFound
		}

		for entity, component := range m.components[t] {
			if existingComponents, ok := result[entity]; ok {
				result[entity] = append(existingComponents, component)
			} else {
				result[entity] = []*Component{component}
			}
		}
	}

	// Filter out entities that do not have all the required types of components
	for entity, components := range result {
		componentTypes := make(map[string]bool)
		for _, component := range components {
			componentTypes[component.Type] = true
		}

		for _, t := range types {
			if !componentTypes[t] {
				delete(result, entity)
				break
			}
		}
	}

	return result, nil
}

// GetComponentData returns the data of a component of the given type from a list of components
func GetComponentData[T any](components []*Component, componentType string) (*T, error) {
	for _, component := range components {
		if component.Type == componentType {
			data, ok := component.Data.(T)
			if !ok {
				return nil, ErrComponentDataMismatch
			}
			return &data, nil
		}
	}
	return nil, ErrComponentNotFound
}
