package app

type Closer interface {
	Close() error
}
