package app

import "context"

// RepoCloser is contract of closing repository
type RepoCloser interface {
	Close(context.Context) error
}
