package peas

import (
	core "github.com/procyon-projects/procyon-core"
	"log"
	"sync"
)

const defaultSharedObjectsMapSize = 0

type SharedPeaRegistry interface {
	RegisterSharedPea(peaName string, sharedObject interface{})
	GetSharedPea(peaName string) interface{}
	ContainsSharedPea(peaName string) bool
	GetSharedPeaNames() []string
	GetSharedPeaCount() int
}

type DefaultSharedPeaRegistry struct {
	sharedObjects              map[string]interface{}
	sharedObjectsInPreparation map[string]interface{}
	sharedObjectsCompleted     map[string]interface{}
	muSharedObjects            sync.RWMutex
}

func NewDefaultSharedPeaRegistry() DefaultSharedPeaRegistry {
	return DefaultSharedPeaRegistry{
		sharedObjects:              make(map[string]interface{}, defaultSharedObjectsMapSize),
		sharedObjectsInPreparation: make(map[string]interface{}, defaultSharedObjectsMapSize),
		sharedObjectsCompleted:     make(map[string]interface{}, defaultSharedObjectsMapSize),
		muSharedObjects:            sync.RWMutex{},
	}
}

func (registry DefaultSharedPeaRegistry) RegisterSharedPea(peaName string, sharedObject interface{}) {
	if peaName == "" || sharedObject == nil {
		core.Log.Error("Pea name or shared object must not be null or empty")
		return
	}
	registry.muSharedObjects.Lock()
	if _, ok := registry.sharedObjects[peaName]; ok {
		panic("Could not register shared object")
	}
	registry.sharedObjects[peaName] = sharedObject
	registry.muSharedObjects.Unlock()
}

func (registry DefaultSharedPeaRegistry) GetSharedPea(peaName string) interface{} {
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

func (registry DefaultSharedPeaRegistry) ContainsSharedPea(peaName string) bool {
	if _, ok := registry.sharedObjects[peaName]; ok {
		return true
	}
	return false
}

func (registry DefaultSharedPeaRegistry) GetSharedPeaNames() []string {
	return core.GetMapKeys(registry.sharedObjects)
}

func (registry DefaultSharedPeaRegistry) GetSharedPeaCount() int {
	return len(registry.sharedObjects)
}
