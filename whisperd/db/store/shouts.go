package store

import (
	"context"

	"whisperd.io/whisperd/whisperd"
)

type Shouts interface {
	GetAll(context.Context) ([]whisperd.Shout, error)
	Get(context.Context, int64) (whisperd.Shout, error)
	Insert(context.Context, whisperd.Shout) (whisperd.Shout, error)
	Delete(context.Context, whisperd.Shout) (whisperd.Shout, error)
}
