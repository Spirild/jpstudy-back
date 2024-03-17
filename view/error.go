package view

type HttpError struct {
	err  error
	code int
}

func (he *HttpError) Error() string {
	return he.err.Error()
}
