package model

type Model interface {
	GetId() interface{}
	SetId(interface{})
	Validate() []error
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf(e.Message)
}
