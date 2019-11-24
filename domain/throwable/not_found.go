package throwable

type NotFound struct {
	err string
}

func (e NotFound) Error() string {
	return e.err
}

func NewNotFound(err string) *NotFound {
	return &NotFound{err: err}
}
