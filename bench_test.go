package context_test

import (
	"testing"

	"github.com/wspowell/context"
)

type contextKey struct{}

var key = contextKey{}

func Benchmark_Background_New(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			context.Background()
		}
	})
}

func Benchmark_Local(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			context.Local()
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
			context.WithValue(ctx, key, "value")
		}
	})
}

func Benchmark_Local_WithLocalValue(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx := context.Local()
			context.WithLocalValue(ctx, key, "value")
		}
	})
}

func Benchmark_Background_Value(b *testing.B) {
	ctx := context.Background()
	context.WithValue(ctx, key, "value")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx.Value(key)
		}
	})
}

func Benchmark_Localized_Value(b *testing.B) {
	ctx := context.Local()
	context.WithLocalValue(ctx, key, "value")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			localCtx := context.Localize(ctx)
			localCtx.Value(key)
		}
	})
}
