package model

import (
	"fmt"
	"time"
)

type Model interface {
	GetId() interface{}
	SetId(interface{})
	Validate() []error
}

type DbGateway interface {
	Create(Model) error        // update injected model
	Save(Model) error          // update injected model
	FindBy(Model, Query) error // update injected model
	FindById(Model) error      // update injected model
	FindAll(Model) (interface{}, error)
	FindAllBy(Model, Query) (interface{}, error)
	Update(Model, map[string]interface{}) error // update injected model
	UpdateAll(Model, map[string]interface{}) (int, error)
	UpdateAllWhere(Model, map[string]interface{}, map[string]interface{}) (int, error)
	DeleteWhere(Model, map[string]interface{}) error
	DeleteById(Model) error
	DeleteAll(Model) (int, error)
	DeleteAllWhere(Model, map[string]interface{}) (int, error)
}

type Query struct {
	Conditions map[string]interface{}
	Order      []string
	Limit      *int
	Offset     *int
}

//TODO: set these from config
const modelTimeFormat = time.RFC3339Nano
const modalTimePrecision = time.Millisecond
const modalTimeZone = "UTC"

type modelTime struct {
	time.Time
}

func (mt modelTime) format() modelTime {
	location, _ := time.LoadLocation(modalTimeZone)
	return modelTime{mt.Time.In(location).Round(modalTimePrecision)}
}

func (mt modelTime) formatString() string {
	return mt.format().Time.Format(modelTimeFormat)
}

func (mt modelTime) String() string {
	return mt.formatString()
}

func (mt modelTime) MarshalText() ([]byte, error) {
	return []byte(mt.formatString()), nil
}

func (mt modelTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + mt.formatString() + `"`), nil
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf(e.Message)
}
