//go:build !release && !race
// +build !release,!race

// Do not use the race detector on this file. These tests are expected to have data races.
// The whole purpose of these test are cause a data race and test that they panic.
package context_test

import (
	"sync"
	"testing"

	"github.com/wspowell/context"
)

func Test_Context_ThreadSafety_Incorrect_Usage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx = context.WithValue(ctx, immutableContextKey{}, immutableValue)

	context.WithLocalValue(ctx, localContextKey{}, localValue)

	if ctx.Value(immutableContextKey{}) != immutableValue {
		t.Errorf("expected immutableContextKey{} == immutableValue")
	}

	if ctx.Value(localContextKey{}) != localValue {
		t.Errorf("expected localContextKey{} == localValue")
	}

	var paniced bool

	var wg sync.WaitGroup
	wg.Add(1)
	go func(boundaryCtx context.Context, wg *sync.WaitGroup) {
		defer func() {
			if err := recover(); err != nil {
				paniced = true
			}
		}()
		defer wg.Done()

		// Context should have been Localized().
		boundaryCtx.Value(localContextKey{})
	}(ctx, &wg)

	wg.Wait()

	if !paniced {
		t.Errorf("expected panic")
	}
}

func Test_Context_ThreadSafety_double_Localized(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx = context.WithValue(ctx, immutableContextKey{}, immutableValue)

	context.WithLocalValue(ctx, localContextKey{}, localValue)

	if ctx.Value(immutableContextKey{}) != immutableValue {
		t.Errorf("expected immutableContextKey{} == immutableValue")
	}

	if ctx.Value(localContextKey{}) != localValue {
		t.Errorf("expected localContextKey{} == localValue")
	}

	var paniced bool

	var wg sync.WaitGroup
	wg.Add(1)
	go func(boundaryCtx context.Context, wg *sync.WaitGroup) {
		defer func() {
			if err := recover(); err != nil {
				paniced = true
			}
		}()
		defer wg.Done()

		localCtx := context.Localize(boundaryCtx)

		// Context is already Localized().
		context.Localize(localCtx)
	}(ctx, &wg)

	wg.Wait()

	if !paniced {
		t.Errorf("expected panic")
	}
}

func Test_WithLocalValue_missing_Localize(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		defer func() {
			if err := recover(); err != nil {
				// Tests passed.
				return
			}
			t.Errorf("expected panic, but found none")
		}()

		ctx = context.WithValue(ctx, immutableContextKey{}, immutableValue)
		ctx = context.WithValue(ctx, duplicateContextKey{}, immutableValue)

		context.WithLocalValue(ctx, localContextKey{}, localValue)
		context.WithLocalValue(ctx, duplicateContextKey{}, duplicateValue)
	}()
	wg.Wait()
}
