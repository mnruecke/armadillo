package test

type MockModel struct {
	Id   int
	Name string
}

func (mm *MockModel) GetId() interface{} {
	return mm.Id
}

func (mm *MockModel) SetId(id interface{}) {
	mm.Id = id.(int)
}

func (mm *MockModel) Validate() []error {
	return []error{}
}
