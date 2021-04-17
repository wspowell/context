package context_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/wspowell/context"
)

type immutableContextKey struct{}
type localContextKey struct{}
type localValueContextKey struct{}
type duplicateContextKey struct{}

var (
	immutableValue = "immutableValue"
	localValue     = "localValue"
	duplicateValue = "myValue"
)

func checkContext(t *testing.T, ctx context.Context) {
	if ctx.Value(localContextKey{}) != nil {
		panic("expected 'localContextKey{}' to be nil")
	}

	if ctx.Value(duplicateContextKey{}).(string) != immutableValue {
		panic(fmt.Sprintf("expected 'duplicatedKey' to be %v but was %v", immutableValue, ctx.Value(duplicateContextKey{})))
	}

	if ctx.Value(immutableContextKey{}).(string) != immutableValue {
		panic(fmt.Sprintf("expected 'immutable' to be %v but was %v", immutableContextKey{}, ctx.Value(immutableContextKey{})))
	}
}

func checkLocal(t *testing.T, ctx context.Context) {
	if ctx.Value(localContextKey{}).(string) != localValue {
		panic(fmt.Sprintf("expected 'localContextKey{}' to be %v but was %v", localValue, ctx.Value(localContextKey{})))
	}

	if ctx.Value(duplicateContextKey{}).(string) != duplicateValue {
		panic(fmt.Sprintf("expected 'duplicatedKey' to be %v but was %v", duplicateValue, ctx.Value(duplicateContextKey{})))
	}

	if ctx.Value(immutableContextKey{}).(string) != immutableValue {
		panic(fmt.Sprintf("expected 'immutable' to be %v but was %v", immutableContextKey{}, ctx.Value(immutableContextKey{})))
	}
}

func Test_Localize(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctx = context.WithValue(ctx, immutableContextKey{}, immutableValue)
	ctx = context.WithValue(ctx, duplicateContextKey{}, immutableValue)

	localCtx := context.Localize(ctx)

	localCtx = context.WithLocalValue(localCtx, localContextKey{}, localValue)
	localCtx = context.WithLocalValue(localCtx, duplicateContextKey{}, duplicateValue)

	checkContext(t, ctx)
	checkLocal(t, localCtx)
}

func Test_Localize_with_WithLocalValue(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctx = context.WithValue(ctx, immutableContextKey{}, immutableValue)
	ctx = context.WithValue(ctx, duplicateContextKey{}, immutableValue)

	localCtx := context.WithLocalValue(ctx, localContextKey{}, localValue)
	localCtx = context.WithLocalValue(localCtx, duplicateContextKey{}, duplicateValue)

	checkContext(t, ctx)
	checkLocal(t, localCtx)
}

func Test_Context_ThreadSafety_Correct_Usage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx = context.WithValue(ctx, immutableContextKey{}, immutableValue)
	ctx = context.WithLocalValue(ctx, localContextKey{}, &localValue)
	ctx = context.WithLocalValue(ctx, localValueContextKey{}, 15)

	if ctx.Value(immutableContextKey{}) != immutableValue {
		t.Errorf("expected immutableContextKey{} == immutableValue")
	}

	if ctx.Value(localContextKey{}) != &localValue {
		t.Errorf("expected localContextKey{} == localValue")
	}

	if ctx.Value(localValueContextKey{}) != 15 {
		t.Errorf("expected localValueContextKey{} == 15")
	}

	var paniced bool

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(boundaryCtx context.Context, wg *sync.WaitGroup) {
			defer func() {
				if err := recover(); err != nil {
					paniced = true
				}
			}()
			defer wg.Done()

			localCtx := context.Localize(boundaryCtx)
			if localCtx.Value(localContextKey{}) == &localValue {
				t.Errorf("local value should not be copied")
			}
			if localCtx.Value(localValueContextKey{}) != 15 {
				t.Errorf("local value should be copied")
			}

			if localCtx.Value(localContextKey{}) == "goroutineValue" {
				t.Errorf("local value should not be shared")
			}

			localCtx = context.WithLocalValue(localCtx, localContextKey{}, "goroutineValue")

			if localCtx.Value(localContextKey{}) != "goroutineValue" {
				t.Errorf("local value should be set")
			}
		}(ctx, &wg)
	}

	wg.Wait()

	if paniced {
		t.Errorf("unexpected panic")
	}

	if ctx.Value(immutableContextKey{}) != immutableValue {
		t.Errorf("expected immutableContextKey{} == immutableValue")
	}

	if ctx.Value(localContextKey{}) != &localValue {
		t.Errorf("expected localContextKey{} == localValue")
	}

	if ctx.Value(localValueContextKey{}) != 15 {
		t.Errorf("expected localValueContextKey{} == 15")
	}
}
