package peas

import (
	"github.com/procyon-projects/goo"
	"sync"
)

type PeaDefinition interface {
	GetTypeName() string
	GetPeaType() goo.Type
	GetScope() PeaScope
}

type SimplePeaDefinitionOption func(definition *SimplePeaDefinition)

type SimplePeaDefinition struct {
	typ   goo.Type
	scope PeaScope
}

func NewSimplePeaDefinition(typ goo.Type, options ...SimplePeaDefinitionOption) *SimplePeaDefinition {
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

func (def *SimplePeaDefinition) GetTypeName() string {
	if def.typ == nil {
		return ""
	}

	if def.typ.IsFunction() {
		fun := def.typ.ToFunctionType()
		if fun.GetFunctionReturnTypeCount() == 1 {
			return fun.GetFunctionReturnTypes()[0].GetFullName()
		}
	}

	return def.typ.String()
}

func (def *SimplePeaDefinition) GetPeaType() goo.Type {
	return def.typ
}

func (def *SimplePeaDefinition) GetScope() PeaScope {
	return def.scope
}

func WithScope(scope PeaScope) SimplePeaDefinitionOption {
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
	GetPeaNamesByType(typ goo.Type) []string
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
	return getStringMapKeys(registry.definitions)
}

func (registry *DefaultPeaDefinitionRegistry) GetPeaDefinitionCount() int {
	return len(registry.definitions)
}

func (registry *DefaultPeaDefinitionRegistry) GetPeaNamesByType(typ goo.Type) []string {
	result := make([]string, 0)
	for peaName, peaDefinition := range registry.definitions {
		peaType := peaDefinition.GetPeaType()

		if peaType.IsFunction() {
			fun := peaType.ToFunctionType()
			if fun.GetFunctionReturnTypeCount() == 1 {
				peaType = fun.GetFunctionReturnTypes()[0]
			} else {
				continue
			}
		}
		match := false

		if typ.IsInterface() && peaType.IsStruct() && peaType.ToStructType().Implements(typ.ToInterfaceType()) {
			match = true
		} else if typ.IsStruct() && peaType.IsStruct() {
			if typ.GetGoType() == peaType.GetGoType() && peaType.IsPointer() && !typ.IsPointer() {
				match = true
			} else if peaType.ToStructType().EmbeddedStruct(typ.ToStructType()) {
				match = true
			} else if typ.GetGoType() == peaType.GetGoType() {
				match = true
			}
		}

		if match {
			result = append(result, peaName)
		}
	}
	return result
}
