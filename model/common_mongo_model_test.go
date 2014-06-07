package model

import (
	"github.com/repp/armadillo/test"
	"testing"
)

type MockCommonModel struct {
	CommonMongoModel
	Name string
}

func (m *MockCommonModel) Validate() []error {
	return []error{}
}

func TestMongoModelImplementsModel(t *testing.T) {
	model := Model(&MockCommonModel{})
	test.AssertNotEqual(t, model, nil)
}
