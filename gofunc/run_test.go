package gofunc_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wspowell/context"
	"github.com/wspowell/context/gofunc"
	"github.com/wspowell/errors"
)

var errTest = errors.New("err", "test")

func Test_Run(t *testing.T) {
	ctx := context.Local()

	wg := sync.WaitGroup{}
	wg.Add(1)

	err := gofunc.Run(ctx, func(ctx context.Context) error {
		defer wg.Done()

		return nil
	})

	wg.Wait()

	assert.Nil(t, <-err)
}

func Test_Run_error(t *testing.T) {
	ctx := context.Local()

	wg := sync.WaitGroup{}
	wg.Add(1)

	err := gofunc.Run(ctx, func(ctx context.Context) error {
		defer wg.Done()

		return errTest
	})

	wg.Wait()

	assert.Equal(t, errTest, <-err)
}

func Test_Run_local_check(t *testing.T) {
	ctx := context.Local()

	wg := sync.WaitGroup{}
	wg.Add(1)

	err := gofunc.Run(ctx, func(ctx context.Context) error {
		defer wg.Done()

		// Should panic since Run() already localized the context.
		context.Localize(ctx)

		return nil
	})

	wg.Wait()

	assert.NotNil(t, <-err)
}

type task struct {
	wg          *sync.WaitGroup
	shouldPanic bool
}

func newTask(shouldPanic bool) *task {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	return &task{
		wg:          wg,
		shouldPanic: shouldPanic,
	}
}

func (self *task) Run(ctx context.Context) error {
	defer self.wg.Done()

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
	job.wg.Wait()

	assert.Nil(t, <-err)
}

func Test_Exec_local_check(t *testing.T) {
	ctx := context.Local()

	job := newTask(true)
	err := gofunc.Exec(ctx, job)
	job.wg.Wait()

	assert.NotNil(t, <-err)
}
