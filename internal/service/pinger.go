package service

import "context"

// Pinger is an interface that pings DB.
type Pinger interface {
	PingDB(context.Context) error
}
