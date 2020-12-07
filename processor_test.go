package peas

import (
	"github.com/codnect/goo"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testPeaProcessor struct {
}

func newTestPeaProcessor() testPeaProcessor {
	return testPeaProcessor{}
}

func (processor testPeaProcessor) BeforePeaInitialization(peaName string, pea interface{}) (interface{}, error) {
	return nil, nil
}

func (processor testPeaProcessor) AfterPeaInitialization(peaName string, pea interface{}) (interface{}, error) {
	return nil, nil
}

func TestPeaProcessors_AddPeaProcessor(t *testing.T) {
	peaProcessors := NewPeaProcessors()
	testPeaProcessor := newTestPeaProcessor()
	err := peaProcessors.AddPeaProcessor(testPeaProcessor)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(peaProcessors.processors))
}

func TestPeaProcessors_AddPeaProcessor_WhenIsInvokedWithNil(t *testing.T) {
	peaProcessors := NewPeaProcessors()
	err := peaProcessors.AddPeaProcessor(nil)
	assert.NotNil(t, err)
	assert.Equal(t, "processor cannot be null", err.Error())
}

func TestPeaProcessors_AddPeaProcessor_WhenIsInvokedWithTheSameProcessor(t *testing.T) {
	peaProcessors := NewPeaProcessors()
	testPeaProcessor := newTestPeaProcessor()
	err := peaProcessors.AddPeaProcessor(testPeaProcessor)
	assert.Nil(t, err)

	processorType := goo.GetType(testPeaProcessor)
	err = peaProcessors.AddPeaProcessor(testPeaProcessor)
	assert.NotNil(t, err)
	assert.Equal(t, "You have already registered this processor : "+processorType.GetFullName(), err.Error())
}

func TestPeaProcessors_RemoveProcessor(t *testing.T) {
	peaProcessors := NewPeaProcessors()
	testPeaProcessor := newTestPeaProcessor()
	err := peaProcessors.AddPeaProcessor(testPeaProcessor)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(peaProcessors.processors))

	peaProcessors.RemoveProcessor(testPeaProcessor)
	assert.Equal(t, 0, len(peaProcessors.processors))
}

func TestPeaProcessors_RemoveAllProcessor(t *testing.T) {
	peaProcessors := NewPeaProcessors()
	testPeaProcessor := newTestPeaProcessor()
	err := peaProcessors.AddPeaProcessor(testPeaProcessor)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(peaProcessors.processors))

	peaProcessors.RemoveAllProcessor()
	assert.Equal(t, 0, len(peaProcessors.processors))
}

func TestPeaProcessors_GetProcessorsCount(t *testing.T) {
	peaProcessors := NewPeaProcessors()
	testPeaProcessor := newTestPeaProcessor()
	err := peaProcessors.AddPeaProcessor(testPeaProcessor)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(peaProcessors.processors))

	assert.Equal(t, 1, peaProcessors.GetProcessorsCount())
}

func TestPeaProcessors_GetProcessors(t *testing.T) {
	peaProcessors := NewPeaProcessors()
	testPeaProcessor := newTestPeaProcessor()
	err := peaProcessors.AddPeaProcessor(testPeaProcessor)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(peaProcessors.GetProcessors()))
}
