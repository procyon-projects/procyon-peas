package peas

import (
	"errors"
	"fmt"
	"github.com/codnect/goo"
	core "github.com/procyon-projects/procyon-core"
	"reflect"
	"sync"
)

type PeaFactory interface {
	GetPea(name string) (interface{}, error)
	GetPeaByNameAndType(name string, typ goo.Type) (interface{}, error)
	GetPeaByNameAndArgs(name string, args ...interface{}) (interface{}, error)
	GetPeaByType(typ goo.Type) (interface{}, error)
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
	readableTypes    map[string]goo.Type
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
		readableTypes:         make(map[string]goo.Type, 0),
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
	val, err := factory.getPeaWith(name, nil)
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

func (factory DefaultPeaFactory) GetPeaByNameAndType(name string, typ goo.Type) (interface{}, error) {
	val, err := factory.getPeaWith(name, typ)
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

func (factory DefaultPeaFactory) GetPeaByType(typ goo.Type) (interface{}, error) {
	val, err := factory.getPeaWith("", typ)
	if err != nil {
		return val, err
	}
	if val == nil && factory.parentPeaFactory != nil {
		if parentPeaFactory, ok := factory.parentPeaFactory.(DefaultPeaFactory); ok {
			return parentPeaFactory.getPeaWith("", typ)
		}
	}
	return val, nil
}

func (factory DefaultPeaFactory) ContainsPea(name string) (interface{}, error) {
	return nil, nil
}

func (factory DefaultPeaFactory) getPeaWith(name string, typ goo.Type, args ...interface{}) (interface{}, error) {
	if name == "" {
		return nil, errors.New("pea name must not be null")
	}
	sharedPea := factory.GetSharedPea(name)
	if sharedPea != nil && args == nil {
		return sharedPea, nil
	} else {
		peaDefinition := factory.GetPeaDefinition(name)
		if SharedScope == peaDefinition.GetScope() {
			instance, err := factory.GetSharedPeaWithObjectFunc(name, func() (instance interface{}, err error) {
				defer func() {
					if r := recover(); r != nil {
						err = NewPeaPreparationError(name, "Creation of pea is failed")
					}
				}()
				instance, err = factory.createPea(name, peaDefinition, args)
				return
			})
			return instance, err
		} else if PrototypeScope == peaDefinition.GetScope() {

		}
	}
	return nil, nil
}

func (factory DefaultPeaFactory) createPea(name string, definition PeaDefinition, args []interface{}) (interface{}, error) {
	instance, err := factory.createPeaInstance(name, definition.GetPeaType(), args)
	if err == nil && definition.GetScope() == SharedScope {
		err = factory.RegisterSharedPea(name, instance)
	}
	return instance, err
}

func (factory DefaultPeaFactory) createPeaInstance(name string, typ goo.Type, args []interface{}) (result interface{}, error error) {
	var instance interface{}
	defer func() {
		if r := recover(); r != nil {
			error = errors.New(fmt.Sprintf("while creating an pea object, an error occurred : %s", name))
		}
	}()
	if typ.IsFunction() {
		constructorFunction := typ.ToFunctionType()
		parameterCount := constructorFunction.GetFunctionParameterCount()
		if parameterCount != 0 && args == nil {
			parameterTypes := constructorFunction.GetFunctionParameterTypes()
			resolvedArguments := factory.createArgumentArray(name, parameterTypes)
			instance, error = CreateInstance(typ, resolvedArguments)
		} else if (parameterCount == 0 && args == nil) || (args != nil && parameterCount == len(args)) {
			instance, error = CreateInstance(typ, args)
		} else {
			error = errors.New("argument count does not match with the parameter count which constructor function has go")
		}
	} else {
		instance, error = CreateInstance(typ, nil)
	}
	if error != nil {
		return
	}
	return factory.initializePea(name, instance)
}

func (factory DefaultPeaFactory) createArgumentArray(name string, parameterTypes []goo.Type) []interface{} {
	argumentArray := make([]interface{}, len(parameterTypes))
	for parameterIndex, parameterType := range parameterTypes {
		peas := factory.resolveDependency(parameterType)
		peaObjectCount := len(peas)
		if peaObjectCount == 0 {
			argumentArray[parameterIndex] = factory.getDefaultValue(parameterType)
		} else if peaObjectCount == 1 {
			instance := peas[0]
			if instance != nil {
				instanceType := goo.GetType(instance)
				if factory.isOnlyReadableType(instanceType) && instanceType.IsPointer() {
					instance = reflect.ValueOf(instance).Elem().Interface()
				} else if instanceType != nil && instanceType.IsPointer() && !parameterType.IsPointer() && parameterType.IsStruct() {
					instance = reflect.ValueOf(instance).Elem().Interface()
				}
			}
			argumentArray[parameterIndex] = instance
		} else {
			panic("Determining which dependency is used cannot be distinguished : " + name)
		}
	}
	return argumentArray
}

func (factory DefaultPeaFactory) resolveDependency(parameterType goo.Type) []interface{} {
	candidateProcessedMap := make(map[string]bool, 0)
	candidates := make([]interface{}, 0)
	if parameterType.IsStruct() || parameterType.IsInterface() {
		names := factory.GetPeaNamesForType(parameterType)
		for _, name := range names {
			candidate, err := factory.GetPea(name)
			if err == nil {
				candidates = append(candidates, candidate)
				candidateProcessedMap[goo.GetType(candidate).GetFullName()] = true
			}
		}
	}
	typeCandidates := factory.GetSharedPeasByType(parameterType)
	for _, typeCandidate := range typeCandidates {
		if _, ok := candidateProcessedMap[goo.GetType(typeCandidate).GetFullName()]; ok {
			continue
		}
		candidates = append(candidates, typeCandidate)
	}
	return candidates
}

func (factory DefaultPeaFactory) getDefaultValue(parameterType goo.Type) interface{} {
	if parameterType.IsInterface() || parameterType.IsArray() || parameterType.IsSlice() || parameterType.IsMap() {
		return nil
	} else if parameterType.IsStruct() {
		if parameterType.IsPointer() {
			return nil
		} else {
			return parameterType.ToStructType().NewInstance()
		}
	} else if parameterType.IsString() {
		return parameterType.ToStringType().NewInstance()
	} else if parameterType.IsBoolean() {
		return parameterType.ToBooleanType().NewInstance()
	} else if parameterType.IsNumber() {
		return parameterType.ToNumberType().NewInstance()
	} else if parameterType.IsFunction() {
		return nil
	}
	panic("Default value cannot be determined, it is not supported :" + parameterType.GetFullName())
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
			result, err = processor.BeforePeaInitialization(name, result)
			if err != nil {
				return result, err
			}
		}
	}
	return result, nil
}

func (factory DefaultPeaFactory) invokePeaInitializers(name string, obj interface{}) error {
	if initializer, ok := obj.(PeaInitializer); ok {
		return initializer.InitializePea()
	}
	return nil
}

func (factory DefaultPeaFactory) applyPeaProcessorsAfterInitialization(name string, obj interface{}) (interface{}, error) {
	result := obj
	var err error
	if factory.GetPeaProcessorsCount() > 0 {
		for _, processor := range factory.GetPeaProcessors() {
			result, err = processor.AfterPeaInitialization(name, result)
			if err != nil {
				return result, err
			}
		}
	}
	return result, nil
}

func (factory DefaultPeaFactory) RegisterTypeAsOnlyReadable(typ goo.Type) error {
	if typ == nil {
		return errors.New("type must not be null")
	}
	factory.muScopes.Lock()
	factory.readableTypes[typ.GetFullName()] = typ
	factory.muScopes.Unlock()
	return nil
}

func (factory DefaultPeaFactory) isOnlyReadableType(typ goo.Type) bool {
	if typ == nil {
		return false
	}
	defer func() {
		factory.muScopes.Unlock()
	}()
	factory.muScopes.Lock()
	for _, readableType := range factory.readableTypes {
		if readableType.IsInterface() && typ.IsStruct() && typ.ToStructType().Implements(readableType.ToInterfaceType()) {
			return true
		} else if readableType.IsStruct() && typ.IsStruct() {
			if readableType.GetGoType() == typ.GetGoType() {
				return true
			}
		}
	}
	return false
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

func (factory DefaultPeaFactory) RegisterTypeToScope(typ goo.Type, scope PeaScope) error {
	if typ == nil {
		return errors.New("type must not be null")
	}
	if scope == nil {
		return errors.New("scope must not be null")
	}
	factory.muScopes.Lock()
	factory.peaTypeScopes[typ.GetGoType()] = scope
	factory.muScopes.Unlock()
	return nil
}

func (factory DefaultPeaFactory) PreInstantiateSharedPeas() {
	peaNames := factory.GetPeaDefinitionNames()
	for _, peaName := range peaNames {
		factory.GetPea(peaName)
	}
}
