package peas

import (
	"errors"
	"fmt"
	core "github.com/procyon-projects/procyon-core"
	"reflect"
	"sync"
)

type PeaFactory interface {
	GetPea(name string) (interface{}, error)
	GetPeaByNameAndType(name string, typ *core.Type) (interface{}, error)
	GetPeaByNameAndArgs(name string, args ...interface{}) (interface{}, error)
	GetPeaByType(typ *core.Type) (interface{}, error)
	ContainsPea(name string) (interface{}, error)
	ClonePeaFactory() PeaFactory
}

type DefaultPeaFactory struct {
	SharedPeaRegistry
	PeaDefinitionRegistry
	peaProcessors    *PeaProcessors
	parentPeaFactory PeaFactory
	peaScopes        map[string]PeaScope
	peaTypeScopes    map[reflect.Type]PeaScope
	muScopes         *sync.RWMutex
}

func NewDefaultPeaFactory(parentPeaFactory PeaFactory) DefaultPeaFactory {
	return DefaultPeaFactory{
		SharedPeaRegistry:     NewDefaultSharedPeaRegistry(),
		PeaDefinitionRegistry: NewDefaultPeaDefinitionRegistry(),
		peaProcessors:         NewPeaProcessors(),
		parentPeaFactory:      parentPeaFactory,
		peaScopes:             make(map[string]PeaScope, 0),
		peaTypeScopes:         make(map[reflect.Type]PeaScope, 0),
		muScopes:              &sync.RWMutex{},
	}
}

func (factory DefaultPeaFactory) SetParentPeaFactory(parent PeaFactory) {
	factory.parentPeaFactory = parent
}

func (factory DefaultPeaFactory) ClonePeaFactory() PeaFactory {
	return DefaultPeaFactory{
		SharedPeaRegistry:     factory.SharedPeaRegistry,
		PeaDefinitionRegistry: factory.PeaDefinitionRegistry,
		peaProcessors:         factory.peaProcessors,
		parentPeaFactory:      factory.parentPeaFactory,
		peaScopes:             factory.peaScopes,
		muScopes:              factory.muScopes,
	}
}

func (factory DefaultPeaFactory) GetPea(name string) (interface{}, error) {
	val, err := factory.getPeaWith(name, nil, nil)
	if err != nil {
		return val, err
	}
	if val == nil && factory.parentPeaFactory != nil {
		if parentPeaFactory, ok := factory.parentPeaFactory.(DefaultPeaFactory); ok {
			return parentPeaFactory.getPeaWith(name, nil, nil)
		}
	}
	return val, nil
}

func (factory DefaultPeaFactory) GetPeaByNameAndType(name string, typ *core.Type) (interface{}, error) {
	val, err := factory.getPeaWith(name, typ, nil)
	if err != nil {
		return val, err
	}
	if val == nil && factory.parentPeaFactory != nil {
		if parentPeaFactory, ok := factory.parentPeaFactory.(DefaultPeaFactory); ok {
			return parentPeaFactory.getPeaWith(name, typ, nil)
		}
	}
	return val, nil
}

func (factory DefaultPeaFactory) GetPeaByNameAndArgs(name string, args ...interface{}) (interface{}, error) {
	val, err := factory.getPeaWith(name, nil, args)
	if err != nil {
		return val, err
	}
	if val == nil && factory.parentPeaFactory != nil {
		if parentPeaFactory, ok := factory.parentPeaFactory.(DefaultPeaFactory); ok {
			return parentPeaFactory.getPeaWith(name, nil, args)
		}
	}
	return val, nil
}

func (factory DefaultPeaFactory) GetPeaByType(typ *core.Type) (interface{}, error) {
	val, err := factory.getPeaWith("", typ, nil)
	if err != nil {
		return val, err
	}
	if val == nil && factory.parentPeaFactory != nil {
		if parentPeaFactory, ok := factory.parentPeaFactory.(DefaultPeaFactory); ok {
			return parentPeaFactory.getPeaWith("", typ, nil)
		}
	}
	return val, nil
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

func (factory DefaultPeaFactory) getPeaMetadataInfo(name string) {

}

func (factory DefaultPeaFactory) createPeaObj(name string, typ *core.Type, args ...interface{}) (result interface{}, error error) {
	var instance interface{}
	defer func() {
		if r := recover(); r != nil {
			error = errors.New(fmt.Sprintf("while creating an pea object, an error occurred : %s", name))
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
	if factory.GetPeaProcessorsCount() > 0 {
		for _, processor := range factory.GetPeaProcessors() {
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
	if factory.GetPeaProcessorsCount() > 0 {
		for _, processor := range factory.GetPeaProcessors() {
			result, err = processor.AfterInitialization(name, result)
			if err != nil {
				return result, err
			}
		}
	}
	return result, nil
}

/* Pea Processors */
func (factory DefaultPeaFactory) AddPeaProcessor(processor PeaProcessor) error {
	return factory.peaProcessors.AddPeaProcessor(processor)
}

func (factory DefaultPeaFactory) GetPeaProcessors() []PeaProcessor {
	return factory.peaProcessors.GetProcessors()
}

func (factory DefaultPeaFactory) GetPeaProcessorsCount() int {
	return factory.peaProcessors.GetProcessorsCount()
}

/* Pea Scope */
func (factory DefaultPeaFactory) RegisterScope(scopeName string, scope PeaScope) error {
	if scopeName == "" {
		return errors.New("scopeName must not be null")
	}
	if scope == nil {
		return errors.New("scope must not be null")
	}
	if SharedScope == scopeName || PrototypeScope == scopeName {
		return errors.New("existing scopes shared and prototype cannot be replaced")
	}
	factory.muScopes.Lock()
	factory.peaScopes[scopeName] = scope
	factory.muScopes.Unlock()
	return nil
}

func (factory DefaultPeaFactory) GetRegisteredScopes() []string {
	return core.GetMapKeys(factory.peaScopes)
}

func (factory DefaultPeaFactory) GetRegisteredScope(scopeName string) PeaScope {
	var scope PeaScope
	factory.muScopes.Lock()
	if val, ok := factory.peaScopes[scopeName]; ok {
		scope = val
	}
	factory.muScopes.Unlock()
	return scope
}

func (factory DefaultPeaFactory) RegisterTypeToScope(typ *core.Type, scope PeaScope) error {
	if typ == nil {
		return errors.New("type must not be null")
	}
	if scope == nil {
		return errors.New("scope must not be null")
	}
	factory.muScopes.Lock()
	scopeType := typ.Typ
	if scopeType.Kind() == reflect.Ptr {
		scopeType = scopeType.Elem()
	}
	factory.peaTypeScopes[scopeType] = scope
	factory.muScopes.Unlock()
	return nil
}
