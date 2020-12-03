package peas

import (
	"errors"
	"github.com/codnect/goo"
	"sync"
)

const defaultSharedObjectsMapSize = 0

type SharedPeaRegistry interface {
	RegisterSharedPea(peaName string, sharedObject interface{}) error
	GetSharedPea(peaName string) interface{}
	ContainsSharedPea(peaName string) bool
	GetSharedPeaNames() []string
	GetSharedPeaCount() int
	GetSharedPeaType(requiredType goo.Type) interface{}
	GetSharedPeasByType(requiredType goo.Type) []interface{}
	GetSharedPeaWithObjectFunc(peaName string, objFunc GetObjectFunc) (interface{}, error)
}

type DefaultSharedPeaRegistry struct {
	sharedObjects              map[string]interface{}
	sharedObjectsInPreparation map[string]interface{}
	sharedObjectsType          map[string]goo.Type
	muSharedObjects            sync.RWMutex
}

func NewDefaultSharedPeaRegistry() *DefaultSharedPeaRegistry {
	return &DefaultSharedPeaRegistry{
		sharedObjects:              make(map[string]interface{}, defaultSharedObjectsMapSize),
		sharedObjectsInPreparation: make(map[string]interface{}, defaultSharedObjectsMapSize),
		sharedObjectsType:          make(map[string]goo.Type, defaultSharedObjectsMapSize),
		muSharedObjects:            sync.RWMutex{},
	}
}

func (registry *DefaultSharedPeaRegistry) RegisterSharedPea(peaName string, sharedObject interface{}) error {
	if peaName == "" || sharedObject == nil {
		return errors.New("pea name or shared object must not be null or empty")
	}
	sharedObjectType := goo.GetType(sharedObject)
	if !sharedObjectType.IsInterface() && !sharedObjectType.IsStruct() {
		return errors.New("pea object must be only instance of struct")
	}
	registry.muSharedObjects.Lock()
	if _, ok := registry.sharedObjects[peaName]; ok {
		registry.muSharedObjects.Unlock()
		return errors.New("could not register shared object with same name")
	}
	registry.sharedObjects[peaName] = sharedObject
	registry.muSharedObjects.Unlock()
	registry.addInstanceSharedObjectsType(peaName, sharedObjectType)
	return nil
}

func (registry *DefaultSharedPeaRegistry) GetSharedPea(peaName string) interface{} {
	defer func() {
		registry.muSharedObjects.Unlock()
	}()
	var result interface{}
	registry.muSharedObjects.Lock()
	if sharedObj, ok := registry.sharedObjects[peaName]; ok {
		result = sharedObj
	}
	return result
}

func (registry *DefaultSharedPeaRegistry) ContainsSharedPea(peaName string) bool {
	defer func() {
		registry.muSharedObjects.Unlock()
	}()
	registry.muSharedObjects.Lock()
	if _, ok := registry.sharedObjects[peaName]; ok {
		return true
	}
	return false
}

func (registry *DefaultSharedPeaRegistry) GetSharedPeaNames() []string {
	defer func() {
		registry.muSharedObjects.Unlock()
	}()
	registry.muSharedObjects.Lock()
	names := getStringMapKeys(registry.sharedObjects)
	return names
}

func (registry *DefaultSharedPeaRegistry) GetSharedPeaCount() int {
	defer func() {
		registry.muSharedObjects.Unlock()
	}()
	registry.muSharedObjects.Lock()
	objectLength := len(registry.sharedObjects)
	return objectLength
}

func (registry *DefaultSharedPeaRegistry) GetSharedPeaType(requiredType goo.Type) interface{} {
	instances := registry.GetSharedPeasByType(requiredType)
	if len(instances) > 1 {
		panic("Instances of required type cannot be distinguished")
	} else if len(instances) == 0 {
		return nil
	}
	return instances[0]
}

func (registry *DefaultSharedPeaRegistry) GetSharedPeasByType(requiredType goo.Type) []interface{} {
	if requiredType == nil {
		panic("Required type must not be nil")
	}
	defer func() {
		if r := recover(); r != nil {
			registry.muSharedObjects.Unlock()
		}
	}()
	instances := make([]interface{}, 0)
	registry.muSharedObjects.Lock()
	for peaName, peaType := range registry.sharedObjectsType {
		match := false
		if peaType.Equals(requiredType) || peaType.GetGoType().ConvertibleTo(requiredType.GetGoType()) {
			match = true
		} else if requiredType.IsInterface() && peaType.ToStructType().Implements(requiredType.ToInterfaceType()) {
			match = true
		} else if requiredType.IsStruct() && peaType.ToStructType().EmbeddedStruct(requiredType.ToStructType()) {
			match = true
		}
		if match {
			instances = append(instances, registry.sharedObjects[peaName])
		}
	}
	registry.muSharedObjects.Unlock()
	return instances
}

func (registry *DefaultSharedPeaRegistry) GetSharedPeaWithObjectFunc(peaName string, objFunc GetObjectFunc) (interface{}, error) {
	sharedPea := registry.GetSharedPea(peaName)
	if sharedPea != nil {
		return sharedPea, nil
	}
	err := registry.addSharedPeaToPreparation(peaName)
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			registry.removedSharedPeaFromPreparation(peaName)
		}
	}()
	var newSharedObj interface{}
	newSharedObj, err = objFunc()
	registry.removedSharedPeaFromPreparation(peaName)
	return newSharedObj, err
}

func (registry *DefaultSharedPeaRegistry) addSharedPeaToPreparation(peaName string) error {
	defer func() {
		registry.muSharedObjects.Unlock()
	}()
	registry.muSharedObjects.Lock()
	if _, ok := registry.sharedObjectsInPreparation[peaName]; ok {
		return NewPeaInPreparationError(peaName)
	}
	registry.sharedObjectsInPreparation[peaName] = nil
	return nil
}

func (registry *DefaultSharedPeaRegistry) removedSharedPeaFromPreparation(peaName string) {
	registry.muSharedObjects.Lock()
	delete(registry.sharedObjectsInPreparation, peaName)
	registry.muSharedObjects.Unlock()
}

func (registry *DefaultSharedPeaRegistry) addInstanceSharedObjectsType(peaName string, typ goo.Type) {
	registry.muSharedObjects.Lock()
	registry.sharedObjectsType[peaName] = typ
	registry.muSharedObjects.Unlock()
}
