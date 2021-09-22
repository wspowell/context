package gofunc

import (
	"github.com/wspowell/context"
	"github.com/wspowell/errors"
)

type RunFn func(ctx context.Context) error

func Run(ctx context.Context, fn RunFn) <-chan error {
	err := make(chan error)

	go func(ctx context.Context) {
		err <- errors.Catch(func() {
			ctx = context.Localize(ctx)
			err <- fn(ctx)
		})
	}(ctx)

	return err
}

type Runnable interface {
	Run(ctx context.Context) error
}

func Exec(ctx context.Context, runnable Runnable) <-chan error {
	return Run(ctx, runnable.Run)
}
