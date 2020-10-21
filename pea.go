package peas

import (
	"github.com/codnect/goo"
)

type ConfigurablePeaFactory interface {
	SharedPeaRegistry
	PeaFactory
	RegisterTypeAsOnlyReadable(typ goo.Type) error
	AddPeaProcessor(processor PeaProcessor) error
	GetPeaProcessors() []PeaProcessor
	GetPeaProcessorsCount() int
	RegisterScope(scopeName string, scope PeaScope) error
	RegisterTypeToScope(typ goo.Type, scope PeaScope) error
	GetRegisteredScopes() []string
	GetRegisteredScope(scopeName string) PeaScope
	SetParentPeaFactory(parent PeaFactory)
	PreInstantiateSharedPeas()
}

type PeaInitializer interface {
	InitializePea() error
}

type PeaFactoryAware interface {
	SetPeaFactory(factory PeaFactory)
}

type PeaNameGenerator interface {
	GenerateName(peaDefinition PeaDefinition) string
}
