package storage

import (
	"errors"
	"fmt"
)

// AddURLError is duplicating error.
type AddURLError struct {
	DuplicateValue string
	Err            error
}

// Error is provided to implement Error interface.
func (ae *AddURLError) Error() string {
	return fmt.Sprintf("%v %v", ae.DuplicateValue, ae.Err)
}

// NewDuplicateError return new AddURLError.
func NewDuplicateError(err error, shortenedURL string) error {
	return &AddURLError{
		DuplicateValue: shortenedURL,
		Err:            err,
	}
}

// ErrDuplicate is duplicate error.
var ErrDuplicate = errors.New("url already exists")

// .............................................................

// NoContentError is error when there is no wanted data.
type NoContentError struct {
	UUID string
	Err  error
}

// Error is provided to implement Error interface.
func (nc *NoContentError) Error() string {
	return fmt.Sprintf("%v %v", nc.UUID, nc.Err)
}

// NewNoContentError returns new NoContentError.
func NewNoContentError(err error, uuid string) error {
	return &NoContentError{
		UUID: uuid,
		Err:  err,
	}
}

// ErrNoContent is  error when there is no wanted data.
var ErrNoContent = errors.New("no urls created by current user ")

// .............................................................
