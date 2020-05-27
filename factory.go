package peas

import (
	core "github.com/Rollcomp/procyon-core"
)

type PeaFactory interface {
	GetPea(name string) (interface{}, error)
	GetPeaByNameAndType(name string, typ *core.Type) (interface{}, error)
	GetPeaByNameAndArgs(name string, args ...interface{}) (interface{}, error)
	GetPeaByType(typ *core.Type) (interface{}, error)
	ContainsPea(name string) (interface{}, error)
}

type DefaultPeaFactory struct {
	SharedPeaRegistry
	peaProcessors *PeaProcessors
}

func NewDefaultPeaFactory() DefaultPeaFactory {
	return DefaultPeaFactory{
		SharedPeaRegistry: NewDefaultSharedPeaRegistry(),
		peaProcessors:     NewPeaProcessors(),
	}
}

func (factory DefaultPeaFactory) GetPea(name string) (interface{}, error) {
	return nil, nil
}

func (factory DefaultPeaFactory) GetPeaByNameAndType(name string, typ *core.Type) (interface{}, error) {
	return nil, nil
}

func (factory DefaultPeaFactory) GetPeaByNameAndArgs(name string, args ...interface{}) (interface{}, error) {
	return nil, nil
}

func (factory DefaultPeaFactory) GetPeaByType(typ *core.Type) (interface{}, error) {
	return nil, nil
}

func (factory DefaultPeaFactory) ContainsPea(name string) (interface{}, error) {
	return nil, nil
}

func (factory DefaultPeaFactory) AddPeaProcessor(processor PeaProcessor) {
	factory.peaProcessors.AddPeaProcessor(processor)
}

func (factory DefaultPeaFactory) GetProcessors() []PeaProcessor {
	return factory.peaProcessors.GetProcessors()
}

func (factory DefaultPeaFactory) GetProcessorsCount() int {
	return factory.peaProcessors.GetProcessorsCount()
}
