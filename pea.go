package peas

import (
	"github.com/codnect/goo"
	core "github.com/procyon-projects/procyon-core"
)

type ConfigurablePeaFactory interface {
	SharedPeaRegistry
	PeaFactory
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
	AfterProperties()
	Initialize() error
}

type PeaFactoryAware interface {
	SetPeaFactory(factory PeaFactory)
}

type PeaMetadataInfo struct {
	typ          goo.Type
	dependencies map[string]interface{}
}

func newPeaMetadataInfo(typ goo.Type, dependencies map[string]interface{}) PeaMetadataInfo {
	return PeaMetadataInfo{
		typ:          typ,
		dependencies: dependencies,
	}
}

func (metadata PeaMetadataInfo) GetType() goo.Type {
	return metadata.typ
}

func (metadata PeaMetadataInfo) GetDependencies() []string {
	return core.GetMapKeys(metadata.typ)
}

type PeaNameGenerator interface {
	GenerateName(peaDefinition PeaDefinition) string
}
