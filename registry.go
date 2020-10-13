package peas

import (
	"errors"
	core "github.com/procyon-projects/procyon-core"
	"log"
	"sync"
)

const defaultSharedObjectsMapSize = 0

type SharedPeaRegistry interface {
	RegisterSharedPea(peaName string, sharedObject interface{}) error
	GetSharedPea(peaName string) interface{}
	ContainsSharedPea(peaName string) bool
	GetSharedPeaNames() []string
	GetSharedPeaCount() int
	GetSharedPeaWithObjectFactory(peaName string, objFunc GetObjectFunc) (interface{}, error)
}

type DefaultSharedPeaRegistry struct {
	sharedObjects              map[string]interface{}
	sharedObjectsInPreparation map[string]interface{}
	sharedObjectsCompleted     map[string]interface{}
	muSharedObjects            sync.RWMutex
}

func NewDefaultSharedPeaRegistry() *DefaultSharedPeaRegistry {
	return &DefaultSharedPeaRegistry{
		sharedObjects:              make(map[string]interface{}, defaultSharedObjectsMapSize),
		sharedObjectsInPreparation: make(map[string]interface{}, defaultSharedObjectsMapSize),
		sharedObjectsCompleted:     make(map[string]interface{}, defaultSharedObjectsMapSize),
		muSharedObjects:            sync.RWMutex{},
	}
}

func (registry *DefaultSharedPeaRegistry) RegisterSharedPea(peaName string, sharedObject interface{}) error {
	if peaName == "" || sharedObject == nil {
		return errors.New("pea name or shared object must not be null or empty")
	}
	registry.muSharedObjects.Lock()
	if _, ok := registry.sharedObjects[peaName]; ok {
		return errors.New("could not register shared object with same name")
	}
	registry.sharedObjects[peaName] = sharedObject
	registry.muSharedObjects.Unlock()
	return nil
}

func (registry *DefaultSharedPeaRegistry) GetSharedPea(peaName string) interface{} {
	var result interface{}
	registry.muSharedObjects.Lock()
	if sharedObj, ok := registry.sharedObjects[peaName]; ok {
		result = sharedObj
	} else if sharedObjInPreparation, ok := registry.sharedObjectsInPreparation[peaName]; ok {
		log.Print(sharedObjInPreparation)
	}
	registry.muSharedObjects.Unlock()
	return result
}

func (registry *DefaultSharedPeaRegistry) ContainsSharedPea(peaName string) bool {
	if _, ok := registry.sharedObjects[peaName]; ok {
		return true
	}
	return false
}

func (registry *DefaultSharedPeaRegistry) GetSharedPeaNames() []string {
	return core.GetMapKeys(registry.sharedObjects)
}

func (registry *DefaultSharedPeaRegistry) GetSharedPeaCount() int {
	return len(registry.sharedObjects)
}

func (registry *DefaultSharedPeaRegistry) GetSharedPeaWithObjectFactory(peaName string, objFunc GetObjectFunc) (interface{}, error) {
	registry.muSharedObjects.Lock()
	if sharedObj, ok := registry.sharedObjects[peaName]; ok {
		registry.muSharedObjects.Unlock()
		return sharedObj, nil
	}
	if registry.isSharedPeaInPreparation(peaName) {
		return nil, NewPeaInPreparationError(peaName)
	} else {
		registry.addSharedPeaToPreparation(peaName)
	}
	newSharedObj, err := objFunc()
	if err != nil {
		registry.muSharedObjects.Unlock()
		return nil, err
	}
	if !registry.isSharedPeaInPreparation(peaName) {
		return nil, NewPeaPreparationError(peaName, "Pea isn't currently in preparation")
	} else {
		registry.removedSharedPeaFromPreparation(peaName)
	}
	registry.muSharedObjects.Unlock()
	return newSharedObj, nil
}

func (registry *DefaultSharedPeaRegistry) isSharedPeaInPreparation(peaName string) bool {
	if _, ok := registry.sharedObjectsInPreparation[peaName]; ok {
		return true
	}
	return false
}

func (registry *DefaultSharedPeaRegistry) addSharedPeaToPreparation(peaName string) {
	registry.sharedObjectsInPreparation[peaName] = nil
}

func (registry *DefaultSharedPeaRegistry) removedSharedPeaFromPreparation(peaName string) {
	delete(registry.sharedObjectsInPreparation, peaName)
}
