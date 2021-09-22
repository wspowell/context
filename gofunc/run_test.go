package gofunc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wspowell/context"
	"github.com/wspowell/context/gofunc"
	"github.com/wspowell/errors"
)

var errTest = errors.New("err", "test")

func Test_Run(t *testing.T) {
	ctx := context.Local()
	err := gofunc.Run(ctx, func(ctx context.Context) error {
		return nil
	})
	assert.Nil(t, <-err)
}

func Test_Run_error(t *testing.T) {
	ctx := context.Local()
	err := gofunc.Run(ctx, func(ctx context.Context) error {
		return errTest
	})
	assert.Equal(t, errTest, <-err)
}

type task struct {
	shouldPanic bool
}

func newTask(shouldPanic bool) *task {
	return &task{
		shouldPanic: shouldPanic,
	}
}

func (self *task) Run(ctx context.Context) error {
	if self.shouldPanic {
		// Should panic since Run() already localized the context.
		context.Localize(ctx)
	}
	return nil
}

func Test_Exec(t *testing.T) {
	ctx := context.Local()
	job := newTask(false)
	err := gofunc.Exec(ctx, job)
	assert.Nil(t, <-err)
}
