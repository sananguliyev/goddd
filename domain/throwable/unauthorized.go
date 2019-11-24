package throwable

type Unauthorized struct {
	err string
}

func (e Unauthorized) Error() string {
	return e.err
}

func NewUnauthorized(err string) *Unauthorized {
	return &Unauthorized{err: err}
}
