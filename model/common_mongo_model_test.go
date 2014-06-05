package model

import (
	"github.com/repp/armadillo/test"
	"testing"
)

type MockModel struct {
	CommonMongoModel
	Name string
}

func (m *MockModel) Validate() []error {
	return []error{}
}

func TestMongoModelImplementsModel(t *testing.T) {
	model := Model(&MockModel{})
	test.AssertNotEqual(t, model, nil)
}
