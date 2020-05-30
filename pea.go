package peas

import core "github.com/Rollcomp/procyon-core"

type ConfigurablePeaFactory interface {
	PeaFactory
	SharedPeaRegistry
}

type PeaInitializer interface {
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
