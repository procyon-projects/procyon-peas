package peas

import (
	"github.com/codnect/goo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimplePeaDefinition(t *testing.T) {
	peaType := goo.GetType(testStruct{})
	peaDefinition1 := NewSimplePeaDefinition(peaType)
	assert.Equal(t, SharedScope, peaDefinition1.GetScope())
	assert.Equal(t, peaType, peaDefinition1.GetPeaType())
	assert.Equal(t, "testStruct", peaDefinition1.GetTypeName())

	peaDefinition2 := NewSimplePeaDefinition(peaType, WithScope(PrototypeScope))
	assert.Equal(t, PrototypeScope, peaDefinition2.GetScope())
	assert.Equal(t, peaType, peaDefinition2.GetPeaType())
	assert.Equal(t, "testStruct", peaDefinition2.GetTypeName())
}

func TestDefaultPeaDefinitionRegistry_RegisterPeaDefinition(t *testing.T) {
	peaDefinitionRegistry := NewDefaultPeaDefinitionRegistry()

	peaType1 := goo.GetType(testStruct{})
	peaDefinition1 := NewSimplePeaDefinition(peaType1)
	peaDefinitionRegistry.RegisterPeaDefinition("testPea1", peaDefinition1)
	assert.Equal(t, 1, len(peaDefinitionRegistry.definitions))

	peaType2 := goo.GetType(testStruct2{})
	peaDefinition2 := NewSimplePeaDefinition(peaType2)
	peaDefinitionRegistry.RegisterPeaDefinition("testPea2", peaDefinition2)
	assert.Equal(t, 2, len(peaDefinitionRegistry.definitions))
}

func TestDefaultPeaDefinitionRegistry_ContainsPeaDefinition(t *testing.T) {
	peaDefinitionRegistry := NewDefaultPeaDefinitionRegistry()

	peaType1 := goo.GetType(testStruct{})
	peaDefinition1 := NewSimplePeaDefinition(peaType1)
	peaDefinitionRegistry.RegisterPeaDefinition("testPea1", peaDefinition1)
	assert.True(t, peaDefinitionRegistry.ContainsPeaDefinition("testPea1"))

	peaType2 := goo.GetType(testStruct2{})
	peaDefinition2 := NewSimplePeaDefinition(peaType2)
	peaDefinitionRegistry.RegisterPeaDefinition("testPea2", peaDefinition2)
	assert.True(t, peaDefinitionRegistry.ContainsPeaDefinition("testPea2"))

	assert.False(t, peaDefinitionRegistry.ContainsPeaDefinition("testPea3"))
}

func TestDefaultPeaDefinitionRegistry_GetPeaDefinition(t *testing.T) {
	peaDefinitionRegistry := NewDefaultPeaDefinitionRegistry()

	peaType1 := goo.GetType(testStruct{})
	peaDefinition1 := NewSimplePeaDefinition(peaType1)
	peaDefinitionRegistry.RegisterPeaDefinition("testPea1", peaDefinition1)
	assert.NotNil(t, peaDefinitionRegistry.GetPeaDefinition("testPea1"))

	peaType2 := goo.GetType(testStruct2{})
	peaDefinition2 := NewSimplePeaDefinition(peaType2)
	peaDefinitionRegistry.RegisterPeaDefinition("testPea2", peaDefinition2)
	assert.NotNil(t, peaDefinitionRegistry.GetPeaDefinition("testPea2"))

	assert.Nil(t, peaDefinitionRegistry.GetPeaDefinition("testPea3"))
}

func TestDefaultPeaDefinitionRegistry_GetPeaDefinitionCount(t *testing.T) {
	peaDefinitionRegistry := NewDefaultPeaDefinitionRegistry()

	peaType1 := goo.GetType(testStruct{})
	peaDefinition1 := NewSimplePeaDefinition(peaType1)
	peaDefinitionRegistry.RegisterPeaDefinition("testPea1", peaDefinition1)
	assert.Equal(t, 1, peaDefinitionRegistry.GetPeaDefinitionCount())

	peaType2 := goo.GetType(testStruct2{})
	peaDefinition2 := NewSimplePeaDefinition(peaType2)
	peaDefinitionRegistry.RegisterPeaDefinition("testPea2", peaDefinition2)
	assert.Equal(t, 2, peaDefinitionRegistry.GetPeaDefinitionCount())
}

func TestDefaultPeaDefinitionRegistry_GetPeaDefinitionNames(t *testing.T) {
	peaDefinitionRegistry := NewDefaultPeaDefinitionRegistry()

	peaType1 := goo.GetType(testStruct{})
	peaDefinition1 := NewSimplePeaDefinition(peaType1)
	peaDefinitionRegistry.RegisterPeaDefinition("testPea1", peaDefinition1)

	peaType2 := goo.GetType(testStruct2{})
	peaDefinition2 := NewSimplePeaDefinition(peaType2)
	peaDefinitionRegistry.RegisterPeaDefinition("testPea2", peaDefinition2)

	peaDefinitionNames := peaDefinitionRegistry.GetPeaDefinitionNames()
	assert.Equal(t, 2, len(peaDefinitionNames))
	assert.Contains(t, peaDefinitionNames, "testPea1")
	assert.Contains(t, peaDefinitionNames, "testPea2")
}

func TestDefaultPeaDefinitionRegistry_GetPeaNamesForType(t *testing.T) {
	peaDefinitionRegistry := NewDefaultPeaDefinitionRegistry()

	peaType1 := goo.GetType(testStruct{})
	peaDefinition1 := NewSimplePeaDefinition(peaType1)
	peaDefinitionRegistry.RegisterPeaDefinition("testPea1", peaDefinition1)

	peaType2 := goo.GetType(testStruct{})
	peaDefinition2 := NewSimplePeaDefinition(peaType2)
	peaDefinitionRegistry.RegisterPeaDefinition("testPea2", peaDefinition2)

	peaNames := peaDefinitionRegistry.GetPeaNamesByType(goo.GetType(testStruct{}))
	assert.NotNil(t, peaNames)
	assert.Equal(t, 2, len(peaNames))
	assert.Contains(t, peaNames, "testPea1")
	assert.Contains(t, peaNames, "testPea2")

	peaNames = peaDefinitionRegistry.GetPeaNamesByType(goo.GetType(testStruct2{}))
	assert.NotNil(t, peaNames)
	assert.Equal(t, 0, len(peaNames))

	peaNames = peaDefinitionRegistry.GetPeaNamesByType(goo.GetType((*testInterface)(nil)))
	assert.NotNil(t, peaNames)
	assert.Equal(t, 2, len(peaNames))
	assert.Contains(t, peaNames, "testPea1")
	assert.Contains(t, peaNames, "testPea2")

	peaNames = peaDefinitionRegistry.GetPeaNamesByType(goo.GetType(baseTestStruct{}))
	assert.NotNil(t, peaNames)
	assert.Equal(t, 2, len(peaNames))
	assert.Contains(t, peaNames, "testPea1")
	assert.Contains(t, peaNames, "testPea2")
}

func TestDefaultPeaDefinitionRegistry_RemovePeaDefinition(t *testing.T) {
	peaDefinitionRegistry := NewDefaultPeaDefinitionRegistry()

	peaType1 := goo.GetType(testStruct{})
	peaDefinition1 := NewSimplePeaDefinition(peaType1)
	peaDefinitionRegistry.RegisterPeaDefinition("testPea1", peaDefinition1)

	peaType2 := goo.GetType(testStruct2{})
	peaDefinition2 := NewSimplePeaDefinition(peaType2)
	peaDefinitionRegistry.RegisterPeaDefinition("testPea2", peaDefinition2)
	assert.Equal(t, 2, peaDefinitionRegistry.GetPeaDefinitionCount())

	peaDefinitionRegistry.RemovePeaDefinition("testPea1")
	peaDefinitionRegistry.RemovePeaDefinition("testPea2")
	assert.Equal(t, 0, peaDefinitionRegistry.GetPeaDefinitionCount())
}
