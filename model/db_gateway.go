package model

type DbGateway interface {
	Create(Model) error        // update injected model
	Save(Model) error          // update injected model
	FindBy(Model, Query) error // update injected model
	FindById(Model) error      // update injected model
	FindAll(Model) ([]Model, error)
	FindAllBy(Model, Query) ([]Model, error)
	Update(Model) error // update injected model
	UpdateAll(Model) (int, error)
	UpdateAllBy(Model, Query) (int, error)
	DeleteById(Model) error
	DeleteBy(Model, Query) (int, error)
	DeleteAll(Model) (int, error)
}

type Query struct {
	Conditions map[string]interface{}
	Order      interface{}
	Limit      int
	Offset     int
}
