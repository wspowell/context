package local_test

import (
	"context"
	"testing"

	"github.com/wspowell/local"
)

func Benchmark_Baseline(b *testing.B) {

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx := context.Background()
			ctx = context.WithValue(ctx, "key", "value")
			ctx.Value("key")
		}
	})
}

func Benchmark_Localized(b *testing.B) {

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			localCtx := local.NewLocalized()
			localCtx.Localize("key", "value")
			localCtx.Value("key")
		}
	})
}
