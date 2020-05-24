package peas

import core "github.com/Rollcomp/procyon-core"

type PeaFactory interface {
	GetPea(name string) (interface{}, error)
	GetPeaByNameAndType(name string, typ *core.Type) (interface{}, error)
	GetPeaByNameAndArgs(name string, args ...interface{}) (interface{}, error)
	GetPeaByType(typ *core.Type) (interface{}, error)
	ContainsPea(name string) (interface{}, error)
}

type SimplePeaFactory struct {
	sharedObjects core.SyncMap
}

func NewSimplePeaFactory() SimplePeaFactory {
	return SimplePeaFactory{
		sharedObjects: core.NewSyncMap(),
	}
}

func (factory SimplePeaFactory) GetPea(name string) (interface{}, error) {
	return factory.getPea(name, nil, nil)
}

func (factory SimplePeaFactory) GetPeaByNameAndType(name string, typ *core.Type) (interface{}, error) {
	return factory.getPea(name, typ, nil)
}

func (factory SimplePeaFactory) GetPeaByNameAndArgs(name string, args ...interface{}) (interface{}, error) {
	return factory.getPea(name, nil, args)
}

func (factory SimplePeaFactory) GetPeaByType(typ *core.Type) (interface{}, error) {
	return factory.getPea("<nil>", typ, nil)
}

func (factory SimplePeaFactory) ContainsPea(name string) (interface{}, error) {
	return factory.getPea(name, nil, nil)
}

func (factory SimplePeaFactory) getPea(name string, typ *core.Type, args ...interface{}) (interface{}, error) {
	peaName := name
	sharedObject := factory.getSharedObject(peaName)
	if sharedObject != nil && len(args) == 0 {
		return sharedObject, nil
	}
	if typ != nil {

	}
	return nil, nil
}

func (factory SimplePeaFactory) getSharedObject(name string) interface{} {
	sharedObj := factory.sharedObjects.Get(name)
	if sharedObj == nil {

	}
	return sharedObj
}
