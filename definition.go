package peas

import (
	core "github.com/procyon-projects/procyon-core"
	"sync"
)

type PeaDefinition interface {
	GetPeaType() *core.Type
	GetScope() string
}

type SimplePeaDefinitionOption func(definition *SimplePeaDefinition)

type SimplePeaDefinition struct {
	typ   *core.Type
	scope string
}

func NewSimplePeaDefinition(typ *core.Type, options ...SimplePeaDefinitionOption) *SimplePeaDefinition {
	def := &SimplePeaDefinition{
		typ: typ,
	}
	for _, option := range options {
		option(def)
	}
	if def.scope == "" {
		def.scope = SharedScope
	}
	return def
}

func (def *SimplePeaDefinition) GetPeaType() *core.Type {
	return def.typ
}

func (def *SimplePeaDefinition) GetScope() string {
	return def.scope
}

func WithScope(scope string) SimplePeaDefinitionOption {
	return func(definition *SimplePeaDefinition) {
		definition.scope = scope
	}
}

type PeaDefinitionRegistry interface {
	RegisterPeaDefinition(peaName string, definition PeaDefinition)
	RemovePeaDefinition(peaName string)
	ContainsPeaDefinition(peaName string) bool
	GetPeaDefinition(peaName string) PeaDefinition
	GetPeaDefinitionNames() []string
	GetPeaDefinitionCount() int
	GetPeaNamesForType(typ *core.Type) []string
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

func (registry *DefaultPeaDefinitionRegistry) GetPeaNamesForType(typ *core.Type) []string {
	result := make([]string, 0)
	for peaName, peaDefinition := range registry.definitions {
		peaType := peaDefinition.GetPeaType()
		if (core.IsInterface(typ) && peaType.Typ.Implements(typ.Typ)) ||
			(core.IsStruct(typ) && (typ.Typ == peaType.Typ)) ||
			(core.IsStruct(typ) && core.IsEmbeddedStruct(typ, peaType)) {
			result = append(result, peaName)
		}
	}
	return result
}
