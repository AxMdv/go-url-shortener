package service

// Closer is an interface that can close DB.
type Closer interface {
	Close() error
}
