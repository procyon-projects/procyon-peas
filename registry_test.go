package peas

import (
	"github.com/codnect/goo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type sharedPeaRegistryMock struct {
	mock.Mock
}

func (registry sharedPeaRegistryMock) RegisterSharedPea(peaName string, sharedObject interface{}) error {
	results := registry.Called(peaName, sharedObject)
	return results.Error(0)
}

func (registry sharedPeaRegistryMock) GetSharedPea(peaName string) interface{} {
	results := registry.Called(peaName)
	return results.Get(0)
}

func (registry sharedPeaRegistryMock) ContainsSharedPea(peaName string) bool {
	results := registry.Called(peaName)
	return results.Bool(0)
}

func (registry sharedPeaRegistryMock) GetSharedPeaNames() []string {
	results := registry.Called()
	if results == nil {
		return nil
	}
	return results.Get(0).([]string)
}

func (registry sharedPeaRegistryMock) GetSharedPeaCount() int {
	results := registry.Called()
	return results.Int(0)
}

func (registry sharedPeaRegistryMock) GetSharedPeaType(requiredType goo.Type) interface{} {
	results := registry.Called(requiredType)
	return results.Get(0)
}

func (registry sharedPeaRegistryMock) GetSharedPeasByType(requiredType goo.Type) []interface{} {
	results := registry.Called(requiredType)
	if results == nil {
		return nil
	}
	return results.Get(0).([]interface{})
}

func (registry sharedPeaRegistryMock) GetSharedPeaWithObjectFunc(peaName string, objFunc GetObjectFunc) (interface{}, error) {
	results := registry.Called(peaName, objFunc)
	return results.Get(0), results.Error(1)
}

func TestDefaultSharedPeaRegistry_RegisterSharedPeaWithEmptyPeaNameOrNilInstance(t *testing.T) {
	peaRegistry := NewDefaultSharedPeaRegistry()
	instance := testStruct{}
	err := peaRegistry.RegisterSharedPea("", instance)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "pea name or shared object must not be null or empty")

	err = peaRegistry.RegisterSharedPea("test", nil)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "pea name or shared object must not be null or empty")
}

func TestDefaultSharedPeaRegistry_RegisterSharedPeaWithNonInterfaceOrNonStruct(t *testing.T) {
	peaRegistry := NewDefaultSharedPeaRegistry()
	err := peaRegistry.RegisterSharedPea("test", "test-instance")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "pea object must be only instance of struct")

	instance := testStruct{}
	err = peaRegistry.RegisterSharedPea("test", instance)
	assert.Nil(t, err)
}

func TestDefaultSharedPeaRegistry_RegisterSharedPeaWithSamePeaName(t *testing.T) {
	peaRegistry := NewDefaultSharedPeaRegistry()

	instance1 := testStruct{}
	err := peaRegistry.RegisterSharedPea("test", instance1)
	assert.Nil(t, err)

	instance2 := testStruct{}
	err = peaRegistry.RegisterSharedPea("test", instance2)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "could not register shared object with same name")
}

func TestDefaultSharedPeaRegistry_RegisterSharedPea(t *testing.T) {
	peaRegistry := NewDefaultSharedPeaRegistry()

	instance1 := testStruct{}
	err := peaRegistry.RegisterSharedPea("test1", instance1)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(peaRegistry.sharedObjects))
	assert.Equal(t, instance1, peaRegistry.sharedObjects["test1"])

	instance2 := testStruct{}
	err = peaRegistry.RegisterSharedPea("test2", instance2)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(peaRegistry.sharedObjects))
	assert.Equal(t, instance2, peaRegistry.sharedObjects["test2"])
}

func TestDefaultSharedPeaRegistry_GetSharedPea(t *testing.T) {
	peaRegistry := NewDefaultSharedPeaRegistry()

	instance1 := testStruct{}
	err := peaRegistry.RegisterSharedPea("test1", instance1)
	assert.Nil(t, err)
	assert.Equal(t, instance1, peaRegistry.GetSharedPea("test1"))
	assert.Nil(t, peaRegistry.GetSharedPea("test2"))
}

func TestDefaultSharedPeaRegistry_ContainsSharedPea(t *testing.T) {
	peaRegistry := NewDefaultSharedPeaRegistry()

	instance1 := testStruct{}
	err := peaRegistry.RegisterSharedPea("test1", instance1)
	assert.Nil(t, err)
	assert.True(t, peaRegistry.ContainsSharedPea("test1"))
	assert.False(t, peaRegistry.ContainsSharedPea("test2"))
}

func TestDefaultSharedPeaRegistry_GetSharedPeaNames(t *testing.T) {
	peaRegistry := NewDefaultSharedPeaRegistry()

	instance1 := testStruct{}
	err := peaRegistry.RegisterSharedPea("test1", instance1)
	assert.Nil(t, err)
	instance2 := testStruct{}
	err = peaRegistry.RegisterSharedPea("test2", instance2)
	assert.Nil(t, err)

	sharedPeaNames := peaRegistry.GetSharedPeaNames()
	assert.Equal(t, 2, len(sharedPeaNames))
	assert.Contains(t, sharedPeaNames, "test1")
	assert.Contains(t, sharedPeaNames, "test2")
}

func TestDefaultSharedPeaRegistry_GetSharedPeaCount(t *testing.T) {
	peaRegistry := NewDefaultSharedPeaRegistry()

	instance1 := testStruct{}
	err := peaRegistry.RegisterSharedPea("test1", instance1)
	assert.Nil(t, err)
	instance2 := testStruct{}
	err = peaRegistry.RegisterSharedPea("test2", instance2)
	assert.Nil(t, err)

	assert.Equal(t, 2, peaRegistry.GetSharedPeaCount())
}

func TestDefaultSharedPeaRegistry_GetSharedPeasByType_WhenIsInvokedWithNil(t *testing.T) {
	peaRegistry := NewDefaultSharedPeaRegistry()
	assert.Panics(t, func() {
		peaRegistry.GetSharedPeasByType(nil)
	})
}

func TestDefaultSharedPeaRegistry_GetSharedPeasByType(t *testing.T) {
	peaRegistry := NewDefaultSharedPeaRegistry()

	instance1 := testStruct{}
	err := peaRegistry.RegisterSharedPea("test1", instance1)
	assert.Nil(t, err)
	instance2 := testStruct{}
	err = peaRegistry.RegisterSharedPea("test2", instance2)
	assert.Nil(t, err)

	peas := peaRegistry.GetSharedPeasByType(goo.GetType(testStruct{}))
	assert.Equal(t, 2, len(peas))

	peas = peaRegistry.GetSharedPeasByType(goo.GetType((*testInterface)(nil)))
	assert.Equal(t, 2, len(peas))

	peas = peaRegistry.GetSharedPeasByType(goo.GetType(baseTestStruct{}))
	assert.Equal(t, 2, len(peas))
}

func TestDefaultSharedPeaRegistry_GetSharedPeaType(t *testing.T) {
	peaRegistry := NewDefaultSharedPeaRegistry()

	instance1 := testStruct{}
	err := peaRegistry.RegisterSharedPea("test1", instance1)
	assert.Nil(t, err)

	pea := peaRegistry.GetSharedPeaType(goo.GetType(testStruct{}))
	assert.NotNil(t, pea)

	pea = peaRegistry.GetSharedPeaType(goo.GetType(testStruct2{}))
	assert.Nil(t, pea)

	pea = peaRegistry.GetSharedPeaType(goo.GetType((*testInterface)(nil)))
	assert.NotNil(t, pea)

	pea = peaRegistry.GetSharedPeaType(goo.GetType(baseTestStruct{}))
	assert.NotNil(t, pea)
}

func TestDefaultSharedPeaRegistry_GetSharedPeaTypeForInstancesDistinguished(t *testing.T) {
	peaRegistry := NewDefaultSharedPeaRegistry()

	instance1 := testStruct{}
	err := peaRegistry.RegisterSharedPea("test1", instance1)
	assert.Nil(t, err)
	instance2 := testStruct{}
	err = peaRegistry.RegisterSharedPea("test2", instance2)
	assert.Nil(t, err)

	assert.Panics(t, func() {
		peaRegistry.GetSharedPeaType(goo.GetType(testStruct{}))
	})

	assert.Panics(t, func() {
		peaRegistry.GetSharedPeaType(goo.GetType((*testInterface)(nil)))
	})

	assert.Panics(t, func() {
		peaRegistry.GetSharedPeaType(goo.GetType(baseTestStruct{}))
	})
}

func TestDefaultSharedPeaRegistry_GetSharedPeaWithObjectFunc(t *testing.T) {
	peaRegistry := NewDefaultSharedPeaRegistry()

	pea, err := peaRegistry.GetSharedPeaWithObjectFunc("test1", func() (i interface{}, err error) {
		return testStruct{}, nil
	})
	assert.NotNil(t, pea)
	assert.Nil(t, err)
	assert.Equal(t, 1, peaRegistry.GetSharedPeaCount())
}
