package peas

import (
	"errors"
	"fmt"
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
	return factory.getPeaWith(name, nil, nil)
}

func (factory DefaultPeaFactory) GetPeaByNameAndType(name string, typ *core.Type) (interface{}, error) {
	return factory.getPeaWith(name, typ, nil)
}

func (factory DefaultPeaFactory) GetPeaByNameAndArgs(name string, args ...interface{}) (interface{}, error) {
	return factory.getPeaWith(name, nil, args)
}

func (factory DefaultPeaFactory) GetPeaByType(typ *core.Type) (interface{}, error) {
	return factory.getPeaWith("", typ, nil)
}

func (factory DefaultPeaFactory) ContainsPea(name string) (interface{}, error) {
	return nil, nil
}

func (factory DefaultPeaFactory) getPeaWith(name string, typ *core.Type, args ...interface{}) (interface{}, error) {
	if name == "" {
		return nil, errors.New("pea name must not be null")
	}
	sharedPea := factory.GetSharedPea(name)
	if sharedPea != nil && args == nil {
		return sharedPea, nil
	}
	return nil, nil
}

func (factory DefaultPeaFactory) createPeaObj(name string, typ *core.Type, args ...interface{}) (interface{}, error) {
	var instance interface{}
	defer func() {
		if r := recover(); r != nil {
			fmt.Print("While creating an pea object, an error occurred : "+name+"\n", r)
		}
	}()
	return factory.initializePea(name, instance)
}

func (factory DefaultPeaFactory) initializePea(name string, obj interface{}) (interface{}, error) {
	/* first of all, invoke pea aware methods */
	factory.invokePeaAware(name, obj)
	result := obj
	var err error
	result, err = factory.applyPeaProcessorsBeforeInitialization(name, result)
	if err != nil {
		return result, err
	}
	err = factory.invokePeaInitializers(name, result)
	if err != nil {
		return result, err
	}
	result, err = factory.applyPeaProcessorsAfterInitialization(name, obj)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (factory DefaultPeaFactory) invokePeaAware(name string, obj interface{}) {
	if aware, ok := obj.(PeaFactoryAware); ok {
		aware.SetPeaFactory(factory)
	}
}

func (factory DefaultPeaFactory) applyPeaProcessorsBeforeInitialization(name string, obj interface{}) (interface{}, error) {
	result := obj
	var err error
	if factory.GetProcessorsCount() > 0 {
		for _, processor := range factory.GetProcessors() {
			result, err = processor.BeforeInitialization(name, result)
			if err != nil {
				return result, err
			}
		}
	}
	return result, nil
}

func (factory DefaultPeaFactory) invokePeaInitializers(name string, obj interface{}) error {
	if initializer, ok := obj.(PeaInitializer); ok {
		return initializer.Initialize()
	}
	return nil
}

func (factory DefaultPeaFactory) applyPeaProcessorsAfterInitialization(name string, obj interface{}) (interface{}, error) {
	result := obj
	var err error
	if factory.GetProcessorsCount() > 0 {
		for _, processor := range factory.GetProcessors() {
			result, err = processor.AfterInitialization(name, result)
			if err != nil {
				return result, err
			}
		}
	}
	return result, nil
}

/* Pea Processors */
func (factory DefaultPeaFactory) AddPeaProcessor(processor PeaProcessor) {
	factory.peaProcessors.AddPeaProcessor(processor)
}

func (factory DefaultPeaFactory) GetProcessors() []PeaProcessor {
	return factory.peaProcessors.GetProcessors()
}

func (factory DefaultPeaFactory) GetProcessorsCount() int {
	return factory.peaProcessors.GetProcessorsCount()
}
