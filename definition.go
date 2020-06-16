package peas

import core "github.com/procyon-projects/procyon-core"

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
