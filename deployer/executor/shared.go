package executor

import (
	"context"
	cd "github.com/containerd/containerd"
)

type state uint

const (
	unspecified state = iota
	created
	deleted
	started
	stopped
)

type cdApi struct {
	client  *cd.Client
	context context.Context
}
