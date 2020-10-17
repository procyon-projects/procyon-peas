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
	GetSharedPeaWithObjectFunc(peaName string, objFunc GetObjectFunc) (interface{}, error)
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
	defer func() {
		registry.muSharedObjects.Unlock()
	}()
	var result interface{}
	registry.muSharedObjects.Lock()
	if sharedObj, ok := registry.sharedObjects[peaName]; ok {
		result = sharedObj
	} else if sharedObjInPreparation, ok := registry.sharedObjectsInPreparation[peaName]; ok {
		log.Print(sharedObjInPreparation)
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
	names := core.GetMapKeys(registry.sharedObjects)
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
