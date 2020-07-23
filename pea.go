package peas

import core "github.com/procyon-projects/procyon-core"

type ConfigurablePeaFactory interface {
	SharedPeaRegistry
	PeaFactory
	AddPeaProcessor(processor PeaProcessor) error
	GetPeaProcessors() []PeaProcessor
	GetPeaProcessorsCount() int
	RegisterScope(scopeName string, scope PeaScope) error
	RegisterTypeToScope(typ *core.Type, scope PeaScope) error
	GetRegisteredScopes() []string
	GetRegisteredScope(scopeName string) PeaScope
	SetParentPeaFactory(parent PeaFactory)
}

type PeaInitializer interface {
	AfterProperties()
	Initialize() error
}

type PeaFactoryAware interface {
	SetPeaFactory(factory PeaFactory)
}

type PeaMetadataInfo struct {
	typ          *core.Type
	dependencies map[string]interface{}
}

func newPeaMetadataInfo(typ *core.Type, dependencies map[string]interface{}) PeaMetadataInfo {
	return PeaMetadataInfo{
		typ:          typ,
		dependencies: dependencies,
	}
}

func (metadata PeaMetadataInfo) GetType() *core.Type {
	return metadata.typ
}

func (metadata PeaMetadataInfo) GetDependencies() []string {
	return core.GetMapKeys(metadata.typ)
}
