package throwable

type InvalidFilter struct {
	err string
}

func (e InvalidFilter) Error() string {
	return e.err
}

func NewInvalidFilter(err string) *InvalidFilter {
	return &InvalidFilter{err: err}
}
