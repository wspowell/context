package context_test

import (

	// nolint:depguard // reason: imported for bench comparison
	gocontext "context"
	"testing"

	"github.com/wspowell/context"
)

type contextKey struct{}

func Benchmark_Background(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			context.Background()
		}
	})
}

func Benchmark_golang_Background(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gocontext.Background()
		}
	})
}

func Benchmark_Localize_New(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			context.Localize(ctx)
		}
	})
}

func Benchmark_Background_WithValue(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx := context.Background()
			context.WithValue(ctx, contextKey{}, "value")
		}
	})
}

func Benchmark_golang_Background_WithValue(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx := gocontext.Background()
			// nolint:govet // reason: benchmark
			gocontext.WithValue(ctx, contextKey{}, "value")
		}
	})
}

func Benchmark_Background_WithLocalValue(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx := context.Background()
			context.WithLocalValue(ctx, contextKey{}, "value")
		}
	})
}

func Benchmark_Background_Value(b *testing.B) {
	ctx := context.Background()
	context.WithValue(ctx, contextKey{}, "value")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx.Value(contextKey{})
		}
	})
}

func Benchmark_golang_Background_Value(b *testing.B) {
	ctx := gocontext.Background()
	ctx = gocontext.WithValue(ctx, contextKey{}, "value")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx.Value(contextKey{})
		}
	})
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
