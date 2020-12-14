package peas

import (
	"errors"
	"github.com/procyon-projects/goo"
	"reflect"
	"sync"
)

type PeaFactory interface {
	GetPea(name string) (interface{}, error)
	GetPeaByNameAndType(name string, typ goo.Type) (interface{}, error)
	GetPeaByNameAndArgs(name string, args ...interface{}) (interface{}, error)
	GetPeaByType(typ goo.Type) (interface{}, error)
	ContainsPea(name string) bool
}

type DefaultPeaFactory struct {
	SharedPeaRegistry
	PeaDefinitionRegistry
	peaProcessors *PeaProcessors
	readableTypes map[string]goo.Type
	muScopes      *sync.RWMutex
}

func NewDefaultPeaFactory() DefaultPeaFactory {
	return DefaultPeaFactory{
		SharedPeaRegistry:     NewDefaultSharedPeaRegistry(),
		PeaDefinitionRegistry: NewDefaultPeaDefinitionRegistry(),
		peaProcessors:         NewPeaProcessors(),
		readableTypes:         make(map[string]goo.Type, 0),
		muScopes:              &sync.RWMutex{},
	}
}

func (factory DefaultPeaFactory) GetPea(name string) (interface{}, error) {
	return factory.getPeaWith(name, nil)
}

func (factory DefaultPeaFactory) GetPeaByNameAndType(name string, typ goo.Type) (interface{}, error) {
	return factory.getPeaWith(name, typ)
}

func (factory DefaultPeaFactory) GetPeaByNameAndArgs(name string, args ...interface{}) (interface{}, error) {
	return factory.getPeaWith(name, nil, args...)
}

func (factory DefaultPeaFactory) GetPeaByType(typ goo.Type) (interface{}, error) {
	return factory.getPeaWith("", typ)
}

func (factory DefaultPeaFactory) ContainsPea(name string) bool {
	return factory.ContainsSharedPea(name)
}

func (factory DefaultPeaFactory) getPeaWith(name string, requiredType goo.Type, args ...interface{}) (interface{}, error) {
	if name == "" && requiredType == nil {
		return nil, errors.New("one of the pea name or type must not be nil at least")
	}

	if name == "" {
		candidatePeaNames := factory.GetPeaNamesByType(requiredType)
		candidatePeaCount := len(candidatePeaNames)
		if candidatePeaCount > 1 {
			return nil, errors.New("there is more than one candidate pea definition for the required type, it cannot be distinguished : " + requiredType.GetPackageFullName())
		} else if candidatePeaCount == 0 {
			return nil, errors.New("pea definition couldn't be found for the required type : " + requiredType.GetPackageFullName())
		}
		name = candidatePeaNames[0]
	}

	sharedPea := factory.GetSharedPea(name)
	if sharedPea != nil && args == nil {

		if requiredType != nil {
			peaDefinition := factory.GetPeaDefinition(name)

			var peaType goo.Type
			if peaDefinition == nil {
				peaType = goo.GetType(sharedPea)
			} else {
				peaType = peaDefinition.GetPeaType()
			}

			if factory.matches(peaType, requiredType) {
				return sharedPea, nil
			}

			return nil, errors.New("instance's type does not match the required type")
		}

		return sharedPea, nil
	}

	peaDefinition := factory.GetPeaDefinition(name)
	if peaDefinition == nil {
		return nil, errors.New("pea definition couldn't be found : " + name)
	}

	if requiredType != nil && !factory.matches(peaDefinition.GetPeaType(), requiredType) {
		return nil, errors.New("pea definition type does not match the required type")
	}

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
		peaType := peaDefinition.GetPeaType()
		if peaType.IsFunction() {
			fun := peaType.ToFunctionType()
			if fun.GetFunctionReturnTypeCount() == 1 {
				peaType = fun.GetFunctionReturnTypes()[0]
			} else {
				return nil, errors.New("pea must have only one return type")
			}
		}
		instance, err := factory.createPeaInstance(name, peaType, args)
		return instance, err
	}

	return nil, errors.New("instance couldn't be created")
}

func (factory DefaultPeaFactory) matches(peaType goo.Type, requiredType goo.Type) bool {
	if peaType.IsFunction() {
		fun := peaType.ToFunctionType()
		if fun.GetFunctionReturnTypeCount() == 1 {
			peaType = fun.GetFunctionReturnTypes()[0]
		}
	}
	match := false
	if peaType.Equals(requiredType) || peaType.GetGoType().ConvertibleTo(requiredType.GetGoType()) {
		match = true
	} else if requiredType.IsInterface() && peaType.ToStructType().Implements(requiredType.ToInterfaceType()) {
		match = true
	} else if requiredType.IsStruct() && peaType.ToStructType().EmbeddedStruct(requiredType.ToStructType()) {
		match = true
	}
	return match
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
			instance := factory.getDefaultValue(parameterType)
			if instance != nil {
				instanceType := goo.GetType(instance)
				if factory.isOnlyReadableType(instanceType) && instanceType.IsPointer() {
					instance = reflect.ValueOf(instance).Elem().Interface()
				} else if instanceType != nil && instanceType.IsPointer() && !parameterType.IsPointer() && !parameterType.IsInterface() {
					instance = reflect.ValueOf(instance).Elem().Interface()
				}
			}
			argumentArray[parameterIndex] = instance
		} else if peaObjectCount == 1 {
			instance := peas[0]
			if instance != nil {
				instanceType := goo.GetType(instance)
				if factory.isOnlyReadableType(instanceType) && instanceType.IsPointer() {
					instance = reflect.ValueOf(instance).Elem().Interface()
				} else if instanceType != nil && instanceType.IsPointer() && !parameterType.IsPointer() && !parameterType.IsInterface() {
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
		names := factory.GetPeaNamesByType(parameterType)
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
	} else if parameterType.IsPointer() {
		return nil
	} else if parameterType.IsStruct() {
		return parameterType.ToStructType().NewInstance()
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

func (factory DefaultPeaFactory) PreInstantiateSharedPeas() {
	peaNames := factory.GetPeaDefinitionNames()
	for _, peaName := range peaNames {
		factory.GetPea(peaName)
	}
}
