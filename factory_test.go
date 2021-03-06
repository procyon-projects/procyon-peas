package peas

import (
	"errors"
	"github.com/procyon-projects/goo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultPeaFactory_GetPeaWithEmptyString(t *testing.T) {
	peaFactory := NewDefaultPeaFactory()
	_, err := peaFactory.GetPea("")
	assert.NotNil(t, err)
}

func TestDefaultPeaFactory_GetPeaByNameAndTypeWithEmptyString(t *testing.T) {
	peaFactory := NewDefaultPeaFactory()
	_, err := peaFactory.GetPeaByNameAndType("", nil)
	assert.NotNil(t, err)
}

func TestDefaultPeaFactory_GetPeaForExistingSharedPea(t *testing.T) {
	peaFactory := NewDefaultPeaFactory()
	testPea := &testStruct{}
	peaFactory.RegisterSharedPea("testPea", testPea)

	pea, err := peaFactory.GetPea("testPea")
	assert.Nil(t, err)
	assert.Equal(t, testPea, pea)

	pea, err = peaFactory.GetPea("testPea")
	assert.Nil(t, err)
	assert.Equal(t, testPea, pea)
}

func TestDefaultPeaFactory_GetPeaByNameAndTypeForExistingSharedPea(t *testing.T) {
	peaFactory := NewDefaultPeaFactory()
	testPea := &testStruct{}
	peaType := goo.GetType(testPea)
	peaFactory.RegisterSharedPea("testPea", testPea)

	pea, err := peaFactory.GetPeaByNameAndType("testPea", peaType)
	assert.Nil(t, err)
	assert.Equal(t, testPea, pea)

	pea, err = peaFactory.GetPeaByNameAndType("testPea", peaType)
	assert.Nil(t, err)
	assert.Equal(t, testPea, pea)

	notMatchPeaType := goo.GetType(testStruct2{})

	pea, err = peaFactory.GetPeaByNameAndType("testPea", notMatchPeaType)
	assert.NotNil(t, err)
	assert.Equal(t, "instance's type does not match the required type", err.Error())

	pea, err = peaFactory.GetPeaByNameAndType("testPea", notMatchPeaType)
	assert.NotNil(t, err)
	assert.Equal(t, "instance's type does not match the required type", err.Error())

	interfaceType := goo.GetType((*testInterface)(nil))
	pea, err = peaFactory.GetPeaByNameAndType("testPea", interfaceType)
	assert.Nil(t, err)
	assert.Equal(t, testPea, pea)

	pea, err = peaFactory.GetPeaByNameAndType("testPea", interfaceType)
	assert.Nil(t, err)
	assert.Equal(t, testPea, pea)

	embeddedStructType := goo.GetType(baseTestStruct{})
	pea, err = peaFactory.GetPeaByNameAndType("testPea", embeddedStructType)
	assert.Nil(t, err)
	assert.Equal(t, testPea, pea)

	pea, err = peaFactory.GetPeaByNameAndType("testPea", embeddedStructType)
	assert.Nil(t, err)
	assert.Equal(t, testPea, pea)
}

func TestDefaultPeaFactory_GetPeaByNameAndTypeForDefinitionWithFunction(t *testing.T) {
	peaFactory := NewDefaultPeaFactory()

	peaType := goo.GetType(aStruct{})
	peaDefinition := NewSimplePeaDefinition(goo.GetType(newAStruct))
	peaFactory.RegisterPeaDefinition("aPea", peaDefinition)

	pea, err := peaFactory.GetPeaByNameAndType("aPea", peaType)
	assert.Nil(t, err)
	assert.NotNil(t, pea)

	pea, err = peaFactory.GetPeaByNameAndType("aPea", peaType)
	assert.Nil(t, err)
	assert.NotNil(t, pea)
}

func TestDefaultPeaFactory_GetPeaByNameAndArgs(t *testing.T) {
	peaFactory := NewDefaultPeaFactory()

	peaType := goo.GetType(newStructFunctionWithParameters)
	peaDefinition := NewSimplePeaDefinition(peaType)
	peaFactory.RegisterPeaDefinition("testPea", peaDefinition)

	pea, err := peaFactory.GetPeaByNameAndArgs("testPea", "test-arg", 10, "test-message")
	assert.Nil(t, err)
	assert.NotNil(t, pea)

	pea2, err := peaFactory.GetPeaByNameAndArgs("testPea", "test-arg", 10, "test-message")
	assert.Nil(t, err)
	assert.Equal(t, pea, pea2)
}

func TestDefaultPeaFactory_GetPeaByTypeForExistingSharedPeaDefinition(t *testing.T) {
	peaFactory := NewDefaultPeaFactory()

	peaType := goo.GetType(testStruct{})
	peaDefinition := NewSimplePeaDefinition(peaType)
	peaFactory.RegisterPeaDefinition("testPea", peaDefinition)

	pea1, err := peaFactory.GetPeaByType(peaType)
	assert.Nil(t, err)
	assert.NotNil(t, pea1)

	pea2, err := peaFactory.GetPeaByType(peaType)
	assert.Nil(t, err)
	assert.Equal(t, &pea1, &pea2)

	embeddedStructType := goo.GetType(baseTestStruct{})
	pea2, err = peaFactory.GetPeaByType(embeddedStructType)
	assert.Nil(t, err)
	assert.Equal(t, &pea1, &pea2)

	pea2, err = peaFactory.GetPeaByType(embeddedStructType)
	assert.Nil(t, err)
	assert.Equal(t, &pea1, &pea2)
}

func TestDefaultPeaFactory_GetPeaByTypeForPrototype(t *testing.T) {
	peaFactory := NewDefaultPeaFactory()

	peaType := goo.GetType(testStruct{})
	peaDefinition := NewSimplePeaDefinition(peaType, WithScope(PrototypeScope))
	peaFactory.RegisterPeaDefinition("testPea", peaDefinition)

	pea1, err := peaFactory.GetPeaByType(peaType)
	assert.Nil(t, err)
	assert.NotNil(t, pea1)

	pea2, err := peaFactory.GetPeaByType(peaType)
	assert.Nil(t, err)
	assert.False(t, &pea1 == &pea2)

	embeddedStructType := goo.GetType(baseTestStruct{})
	pea2, err = peaFactory.GetPeaByType(embeddedStructType)
	assert.Nil(t, err)
	assert.False(t, &pea1 == &pea2)

	pea2, err = peaFactory.GetPeaByType(embeddedStructType)
	assert.Nil(t, err)
	assert.False(t, &pea1 == &pea2)
}

func TestDefaultPeaFactory_ContainsPea(t *testing.T) {
	peaFactory := NewDefaultPeaFactory()
	testPea := testStruct{}
	peaFactory.RegisterSharedPea("testPea", testPea)

	assert.True(t, peaFactory.ContainsSharedPea("testPea"))
	assert.False(t, peaFactory.ContainsSharedPea("testPea2"))

	peaType := goo.GetType(testStruct2{})
	peaDefinition := NewSimplePeaDefinition(peaType)
	peaFactory.RegisterPeaDefinition("testPea2", peaDefinition)
	assert.False(t, peaFactory.ContainsPea("testPea2"))

	peaFactory.GetPeaByType(peaType)
	assert.True(t, peaFactory.ContainsPea("testPea2"))
}

func TestDefaultPeaFactory_PeaProcessors(t *testing.T) {
	peaFactory := NewDefaultPeaFactory()

	peaType := goo.GetType(testStruct{})
	peaDefinition := NewSimplePeaDefinition(peaType, WithScope(PrototypeScope))
	peaFactory.RegisterPeaDefinition("testPea", peaDefinition)
	testPeaProcessor := newTestPeaProcessor()
	peaFactory.AddPeaProcessor(testPeaProcessor)

	pea, err := peaFactory.GetPea("testPea")
	assert.Nil(t, err)
	assert.NotNil(t, pea)

	peaErr := errors.New("pea error")
	testPeaProcessor.errBeforePeaInitialization = peaErr
	pea, err = peaFactory.GetPea("testPea")
	assert.NotNil(t, err)
	assert.Equal(t, "pea error", err.Error())
	//assert.Nil(t, pea)

	testPeaProcessor.errBeforePeaInitialization = nil
	testPeaProcessor.errAfterPeaInitialization = peaErr
	pea, err = peaFactory.GetPea("testPea")
	assert.NotNil(t, err)
	assert.Equal(t, "pea error", err.Error())
	//assert.Nil(t, pea)
}

type aStruct struct {
}

func newAStruct() aStruct {
	return aStruct{}
}

type bStruct struct {
}

func newBStruct(a aStruct) bStruct {
	return bStruct{}
}

func TestDefaultPeaFactory_CreatePea_DependencyInjection(t *testing.T) {
	peaFactory := NewDefaultPeaFactory()

	peaType := goo.GetType(newAStruct)
	peaDefinition := NewSimplePeaDefinition(peaType)
	peaFactory.RegisterPeaDefinition("aPea", peaDefinition)

	peaType = goo.GetType(newBStruct)
	peaDefinition = NewSimplePeaDefinition(peaType)
	peaFactory.RegisterPeaDefinition("bPea", peaDefinition)

	pea, err := peaFactory.GetPea("bPea")
	assert.Nil(t, err)
	assert.NotNil(t, pea)

	pea, err = peaFactory.GetPea("aPea")
	assert.Nil(t, err)
	assert.NotNil(t, pea)
}

type cStruct struct {
}

func newCStruct(i testInterface, array []string, m map[string]interface{}, aStruct aStruct, s string, b bool, n int) cStruct {
	return cStruct{}
}

func TestDefaultPeaFactory_ResolverDependencyForDefaultValues(t *testing.T) {
	peaFactory := NewDefaultPeaFactory()

	peaType := goo.GetType(newCStruct)
	peaDefinition := NewSimplePeaDefinition(peaType)
	peaFactory.RegisterPeaDefinition("cPea", peaDefinition)

	pea, err := peaFactory.GetPea("cPea")
	assert.Nil(t, err)
	assert.NotNil(t, pea)
}

func TestDefaultPeaFactory_RegisterTypeAsOnlyReadable(t *testing.T) {
	peaFactory := NewDefaultPeaFactory()
	peaType1 := goo.GetType(testStruct2{})
	err := peaFactory.RegisterTypeAsOnlyReadable(peaType1)
	assert.Nil(t, err)

	peaType2 := goo.GetType((*testInterface)(nil))
	err = peaFactory.RegisterTypeAsOnlyReadable(peaType2)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(peaFactory.readableTypes))

	err = peaFactory.RegisterTypeAsOnlyReadable(nil)
	assert.NotNil(t, err)
	assert.Equal(t, "type must not be null", err.Error())

	assert.False(t, peaFactory.isOnlyReadableType(nil))
	assert.True(t, peaFactory.isOnlyReadableType(peaType1))
	assert.True(t, peaFactory.isOnlyReadableType(goo.GetType(testStruct{})))
}

func TestDefaultPeaFactory_PreInstantiateSharedPeas(t *testing.T) {
	peaFactory := NewDefaultPeaFactory()

	peaType := goo.GetType(newCStruct)
	peaDefinition := NewSimplePeaDefinition(peaType)
	peaFactory.RegisterPeaDefinition("cPea", peaDefinition)

	peaFactory.PreInstantiateSharedPeas()
}
