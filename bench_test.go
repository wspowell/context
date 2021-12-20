package context_test

import (

	// nolint:depguard // reason: imported for bench comparison
	gocontext "context"
	"fmt"
	"io"
	"testing"

	"github.com/wspowell/context"
)

type contextKey struct{}

func Benchmark_Background(b *testing.B) {
	var ctx context.Context

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx = context.Background()
	}

	b.StopTimer()
	fmt.Fprintf(io.Discard, "%v", ctx)
	b.StartTimer()
}

func Benchmark_golang_Background(b *testing.B) {
	var ctx gocontext.Context

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx = gocontext.Background()
	}

	b.StopTimer()
	fmt.Fprintf(io.Discard, "%v", ctx)
	b.StartTimer()
}

func Benchmark_Background_WithValue(b *testing.B) {
	var ctx context.Context

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx = context.Background()
		ctx = context.WithValue(ctx, contextKey{}, "value")
	}

	b.StopTimer()
	fmt.Fprintf(io.Discard, "%v", ctx)
	b.StartTimer()
}

func Benchmark_golang_Background_WithValue(b *testing.B) {
	var ctx gocontext.Context

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx = gocontext.Background()
		// nolint:govet // reason: benchmark
		ctx = gocontext.WithValue(ctx, contextKey{}, "value")
	}

	b.StopTimer()
	fmt.Fprintf(io.Discard, "%v", ctx)
	b.StartTimer()
}

func Benchmark_Background_Value(b *testing.B) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextKey{}, "value")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx.Value(contextKey{})
	}

	b.StopTimer()
	fmt.Fprintf(io.Discard, "%v", ctx)
	b.StartTimer()

}

func Benchmark_golang_Background_Value(b *testing.B) {
	ctx := gocontext.Background()
	ctx = gocontext.WithValue(ctx, contextKey{}, "value")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx.Value(contextKey{})
	}

	b.StopTimer()
	fmt.Fprintf(io.Discard, "%v", ctx)
	b.StartTimer()
}

func Benchmark_Localized_Value(b *testing.B) {
	ctx := context.Background()
	context.WithLocalValue(ctx, contextKey{}, "value")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			localCtx := context.Localize(ctx)
			localCtx.Value(contextKey{})
		}
	})
}

func Benchmark_Background_WithLocalValue(b *testing.B) {
	var ctx gocontext.Context

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx = context.Background()
		context.WithLocalValue(ctx, contextKey{}, "value")
	}

	b.StopTimer()
	fmt.Fprintf(io.Discard, "%v", ctx)
	b.StartTimer()
}
