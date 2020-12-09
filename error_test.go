package peas

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPeaPreparationError_Error(t *testing.T) {
	err := NewPeaPreparationError("test-pea", "test error message")
	assert.Equal(t, "test-pea : test error message", err.Error())
}

func TestPeaPreparationError_GetPeaName(t *testing.T) {
	err := NewPeaPreparationError("test-pea", "test error message")
	assert.Equal(t, "test-pea", err.GetPeaName())
}

func TestPeaPreparationError_GetMessage(t *testing.T) {
	err := NewPeaPreparationError("test-pea", "test error message")
	assert.Equal(t, "test error message", err.GetMessage())
}

func TestPeaInPreparationError_Error(t *testing.T) {
	err := NewPeaInPreparationError("test-pea")
	assert.Equal(t, "test-pea : Pea is currently in preparation, maybe it has got circular dependency cycle", err.Error())
}

func TestPeaInPreparationError_GetPeaName(t *testing.T) {
	err := NewPeaInPreparationError("test-pea")
	assert.Equal(t, "test-pea", err.GetPeaName())
}

func TestPeaInPreparationError_GetMessage(t *testing.T) {
	err := NewPeaInPreparationError("test-pea")
	assert.Equal(t, "Pea is currently in preparation, maybe it has got circular dependency cycle", err.GetMessage())
}
