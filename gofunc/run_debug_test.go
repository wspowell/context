//go:build !release
// +build !release

package gofunc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wspowell/context"
	"github.com/wspowell/context/gofunc"
)

func Test_Run_local_check(t *testing.T) {
	ctx := context.Local()
	err := gofunc.Run(ctx, func(ctx context.Context) error {
		// Should panic since Run() already localized the context.
		context.Localize(ctx)
		return nil
	})
	assert.NotNil(t, <-err)
}

func Test_Exec_local_check(t *testing.T) {
	ctx := context.Local()
	job := newTask(true)
	err := gofunc.Exec(ctx, job)
	assert.NotNil(t, <-err)
}
