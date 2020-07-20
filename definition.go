package peas

import (
	core "github.com/procyon-projects/procyon-core"
	"sync"
)

type PeaDefinition interface {
	GetName() string
	GetPeaType() *core.Type
	GetScope() string
}

type SimplePeaDefinition struct {
	name  string
	typ   *core.Type
	scope string
}

func NewSimplePeaDefinition(name string, typ *core.Type, scope string) SimplePeaDefinition {
	return SimplePeaDefinition{
		name,
		typ,
		scope,
	}
}

func (def SimplePeaDefinition) GetName() string {
	return def.name
}

func (def SimplePeaDefinition) GetPeaType() *core.Type {
	return def.typ
}

func (def SimplePeaDefinition) GetScope() string {
	return def.scope
}

type PeaDefinitionRegistry interface {
	RegisterPeaDefinition(peaName string, definition PeaDefinition)
	RemovePeaDefinition(peaName string)
	ContainsPeaDefinition(peaName string) bool
	GetPeaDefinition(peaName string) PeaDefinition
	GetPeaDefinitionNames() []string
	GetPeaDefinitionCount() int
}

type DefaultPeaDefinitionRegistry struct {
	definitions map[string]PeaDefinition
	mu          sync.RWMutex
}

func NewDefaultPeaDefinitionRegistry() *DefaultPeaDefinitionRegistry {
	return &DefaultPeaDefinitionRegistry{
		definitions: make(map[string]PeaDefinition, 0),
		mu:          sync.RWMutex{},
	}
}

func (registry *DefaultPeaDefinitionRegistry) RegisterPeaDefinition(peaName string, definition PeaDefinition) {
	registry.mu.Lock()
	registry.definitions[peaName] = definition
	registry.mu.Unlock()
}

func (registry *DefaultPeaDefinitionRegistry) RemovePeaDefinition(peaName string) {
	registry.mu.Lock()
	if _, ok := registry.definitions[peaName]; ok {
		delete(registry.definitions, peaName)
	}
	registry.mu.Unlock()
}

func (registry *DefaultPeaDefinitionRegistry) ContainsPeaDefinition(peaName string) bool {
	var result bool
	registry.mu.Lock()
	_, result = registry.definitions[peaName]
	registry.mu.Unlock()
	return result
}

func (registry *DefaultPeaDefinitionRegistry) GetPeaDefinition(peaName string) PeaDefinition {
	var def PeaDefinition
	registry.mu.Lock()
	if val, ok := registry.definitions[peaName]; ok {
		def = val
	}
	registry.mu.Unlock()
	return def
}

func (registry *DefaultPeaDefinitionRegistry) GetPeaDefinitionNames() []string {
	return core.GetMapKeys(registry.definitions)
}

func (registry *DefaultPeaDefinitionRegistry) GetPeaDefinitionCount() int {
	return len(registry.definitions)
}
