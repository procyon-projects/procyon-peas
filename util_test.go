package peas

import (
	"github.com/codnect/goo"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testInterface interface {
	testMethod()
}

type baseTestStruct struct {
}

type testStruct2 struct {
}

type testStruct struct {
	baseTestStruct
}

func (testStruct) testMethod() {

}

func newStructFunction() testStruct {
	return testStruct{}
}

func newStructFunctionWithParameters(str string, i int, obj interface{}) testStruct {
	return testStruct{}
}

func newStructFunctionWithMoreReturnValuesThanOne() (testStruct, error) {
	return testStruct{}, nil
}

func TestCreateInstance_WhenIsInvokedWithFunctionNotHavingAnyParameters(t *testing.T) {
	instance, err := CreateInstance(goo.GetType(newStructFunction), nil)
	assert.NotNil(t, instance)
	assert.Nil(t, err)
}

func TestCreateInstance_WhenIsInvokedWithFunctionHavingParameters(t *testing.T) {
	args := make([]interface{}, 0)
	args = append(args, "test-arg")
	args = append(args, 0)
	args = append(args, testStruct{})
	instance, err := CreateInstance(goo.GetType(newStructFunctionWithParameters), args)
	assert.NotNil(t, instance)
	assert.Nil(t, err)
}

func TestCreateInstance_WhenIsInvokedWithFunctionReturningMoreValuesThanOne(t *testing.T) {
	args := make([]interface{}, 0)
	args = append(args, "test-arg")
	instance, err := CreateInstance(goo.GetType(newStructFunctionWithMoreReturnValuesThanOne), nil)
	assert.Nil(t, instance)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "it only supports the construction functions with one return parameter")
}

func TestCreateInstance_WhenIsInvokedWithStructAndArgs(t *testing.T) {
	args := make([]interface{}, 0)
	args = append(args, "test-arg")
	instance, err := CreateInstance(goo.GetType(testStruct{}), args)
	assert.Nil(t, instance)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "struct type does not support args")
}

func TestCreateInstance_WhenIsInvokedWithNonStructAndFunction(t *testing.T) {
	instance, err := CreateInstance(goo.GetType("test"), nil)
	assert.Nil(t, instance)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "you can only pass Struct or Func types")
}

func TestGetStringMapKeys_WhenIsInvokedWithNil(t *testing.T) {
	assert.Nil(t, getStringMapKeys(nil))
}

func TestGetStringMapKeys_WhenIsInvokedWithMapHavingNonStringKeyType(t *testing.T) {
	assert.Panics(t, func() {
		getStringMapKeys(new(map[interface{}]interface{}))
	})
}

func TestGetStringMapKeys_WhenIsInvokedWithNonMapType(t *testing.T) {
	assert.Panics(t, func() {
		getStringMapKeys("test")
	})
}

func TestGetStringMapKeys_WhenIsInvokedWithMapHavingStringKeyType(t *testing.T) {
	testMap := make(map[string]interface{})
	testMap["key1"] = nil
	testMap["key2"] = nil
	testMap["key3"] = nil
	assert.NotPanics(t, func() {
		keys := getStringMapKeys(testMap)
		assert.Equal(t, 3, len(keys))
		assert.Contains(t, keys, "key1")
		assert.Contains(t, keys, "key2")
		assert.Contains(t, keys, "key3")
		assert.NotContains(t, keys, "key4")
	})
}
