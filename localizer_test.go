package local_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/wspowell/local"
)

type immutableContextKey struct{}
type localContextKey struct{}
type duplicateContextKey struct{}

var (
	immutableKey = immutableContextKey{}
	localKey     = localContextKey{}
	duplicateKey = duplicateContextKey{}

	immutableValue = "immutableValue"
	localValue     = "localValue"
	duplicateValue = "myValue"
)

func checkContext(t *testing.T, ctx context.Context) {
	if ctx.Value(localKey) != nil {
		panic("expected 'localKey' to be nil")
	}

	if ctx.Value(duplicateKey).(string) != immutableValue {
		panic(fmt.Sprintf("expected 'duplicatedKey' to be %v but was %v", immutableValue, ctx.Value(duplicateKey)))
	}

	if ctx.Value(immutableKey).(string) != immutableValue {
		panic(fmt.Sprintf("expected 'immutable' to be %v but was %v", immutableKey, ctx.Value(immutableKey)))
	}
}

func checkLocal(t *testing.T, ctx local.Context) {
	if ctx.Value(localKey).(string) != localValue {
		panic(fmt.Sprintf("expected 'localKey' to be %v but was %v", localValue, ctx.Value(localKey)))
	}

	if ctx.Value(duplicateKey).(string) != duplicateValue {
		panic(fmt.Sprintf("expected 'duplicatedKey' to be %v but was %v", duplicateValue, ctx.Value(duplicateKey)))
	}

	if ctx.Value(immutableKey).(string) != immutableValue {
		panic(fmt.Sprintf("expected 'immutable' to be %v but was %v", immutableKey, ctx.Value(immutableKey)))
	}
}

func Test_NewLocalized(t *testing.T) {
	t.Parallel()

	localCtx := local.NewLocalized()

	local.WithValue(localCtx, immutableKey, immutableValue)
	local.WithValue(localCtx, duplicateKey, immutableValue)

	localCtx.Localize(localKey, localValue)
	localCtx.Localize(duplicateKey, duplicateValue)

	checkContext(t, localCtx.GetContext())
	checkLocal(t, localCtx)
}

func Test_FromContext(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), immutableKey, immutableValue)
	ctx = context.WithValue(ctx, duplicateKey, immutableValue)

	localCtx := local.FromContext(ctx)
	localCtx.Localize(localKey, localValue)
	localCtx.Localize(duplicateKey, duplicateValue)

	checkContext(t, localCtx.GetContext())
	checkLocal(t, localCtx)
}

func Test_Context_ThreadSafety_Correct_Usage(t *testing.T) {
	t.Parallel()

	localCtx := local.NewLocalized()
	localCtx.Localize(localKey, localValue)

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

		ctx := local.FromContext(boundaryCtx)
		if ctx.Value(localKey) == localValue {
			t.Errorf("local value should not be copied")
		}

		ctx.Localize(localKey, "goroutineValue")

	}(localCtx, &wg)

	wg.Wait()

	if paniced {
		t.Errorf("unexpected panic")
	}

	if localCtx.Value(localKey) != localValue {
		t.Errorf("expected localKey == localValue")
	}
}
