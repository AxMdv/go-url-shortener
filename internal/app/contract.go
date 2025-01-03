package app

// Closer is contract of closing repository
type Closer interface {
	Close() error
}
