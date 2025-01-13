package service

// RepoCloser is an interface that can close DB.
type RepoCloser interface {
	IRepository
	Close() error
}
