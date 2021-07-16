package service

type ServerError interface {
	Error() string
	ErrorType() int
}

type ServiceError struct {
	Cause error `json:"error"`
	Type  int   `json:"-"`
}

func (e *ServiceError) Error() string {
	return e.Cause.Error()
}

func (e *ServiceError) ErrorType() int {
	return e.Type
}
