package local_test

import (
	"context"
	"testing"

	"github.com/wspowell/local"
)

type contextKey struct{}

var key = contextKey{}

func Benchmark_Context_New(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			context.Background()
		}
	})
}

func Benchmark_Localized_New(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			local.NewLocalized()
		}
	})
}

func Benchmark_Context_WithValue(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx := context.Background()
			context.WithValue(ctx, key, "value")
		}
	})
}

func Benchmark_Localized_Localize(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			localCtx := local.NewLocalized()
			localCtx.Localize(key, "value")
		}
	})
}

func Benchmark_Context_Value(b *testing.B) {
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
	localCtx := local.NewLocalized()
	localCtx.Localize(key, "value")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			localCtx.Value(key)
		}
	})
}
