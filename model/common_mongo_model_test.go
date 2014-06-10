package model

import (
	"github.com/repp/armadillo/test"
	"labix.org/v2/mgo/bson"
	"time"
	"testing"
)

type MockMongoModel struct {
	CommonMongoModel `bson:",inline"`
	Name string
}

func (m *MockMongoModel) Validate() []error {
	return []error{}
}

func TestMongoModelImplementsModel(t *testing.T) {
	model := Model(&MockMongoModel{})
	test.AssertNotEqual(t, model, nil)
}

func TestModelInitialize(t *testing.T) {
	mock := &MockMongoModel{Name: "Mogi"}
	mock.Initialize()

	test.AssertEqual(t, mock.Name, "Mogi")
	test.AssertNotEqual(t, mock.GetId(), nil)
	test.AssertTypeMatch(t, mock.GetId(), bson.NewObjectId())
	test.AssertTypeMatch(t, mock.GetUpdated(), time.Now())
	test.AssertTypeMatch(t, mock.GetCreated(), time.Now())
}
