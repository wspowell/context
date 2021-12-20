# context 

Replacement for Golang Context

## Context

Golang `context.Context` is a feature that is easily abused. A `context.Context` should only be used for immutable data and are meant to be passed between API boundaries (and therefore must be thread safe). However, it is extremely tempting (and easy) to violate this contract and use it as a generic variable store for values used throughout an goroutines lifetime. 

`context.Context` attempts to address these issues. A `context.Context` is both a `context.Context` and a variable store for goroutine local data. The difference is that `context.Context` provides behavior to localize data to the goroutine. Localized data is not thread safe and must never be sent across API boundaries. Localizing a context to a goroutine will cut out the local data and only allow access to the immutable context data. If localized data implements `Localize() any`, then the value will be cloned in the localized context. `Localize() any` must return a thread safe value.

## Building

The package utilizes goroutine identification (that Golang authors created) to catch threading issues during development. The downside is that tracking goroutine IDs adds significant overhead. In order to solve this, the default build will use goroutine tracking. To build a "release" build, use -tags option with "release".

`go build -tags release ./...`

## Example

```
// Create a new context.
ctx := context.Local()

// Add immutable data.
ctx = context.WithValue(ctx, logConfigKey{}, logConfig)

// Add local data.
// Stored map may be accessed and altered at any time during the goroutine.
context.WithLocalValue(ctx, loggerKey{}, NewLogger())

// Start a goroutine to process data.
go processData(ctx)

...

type Log struct {
    ...
}

func NewLogger() Log {
    ...
}

// Localize Log to the new local Context.
func (self Log) Localize() any {
    ...
}

func processData(ctx context.Context) {
    // Create a new context local to this goroutine.
    // Context no longer has access to the localized keys.
    // Values are shadowed as nil, unless values implement Localize() any.
    ctx := context.Localize(ctx)
    
    ...
}
```

# Benchmarks

Take benchmarks with a bucket of salt.

Debug
```
go test -bench=. -benchmem -count=1 -cpu 8 -parallel 8

goos: linux
goarch: amd64
pkg: github.com/wspowell/context
cpu: AMD Ryzen 9 4900HS with Radeon Graphics         
Benchmark_Background-8                            288093              4170 ns/op             120 B/op          3 allocs/op
Benchmark_golang_Background-8                   1000000000               0.4853 ns/op          0 B/op          0 allocs/op
Benchmark_Background_WithValue-8                  277447              4288 ns/op             168 B/op          4 allocs/op
Benchmark_golang_Background_WithValue-8         22266259                52.83 ns/op           48 B/op          1 allocs/op
Benchmark_Background_Value-8                    295668966                4.046 ns/op           0 B/op          0 allocs/op
Benchmark_golang_Background_Value-8             294940878                4.104 ns/op           0 B/op          0 allocs/op
Benchmark_Localized_Value-8                        65911             17042 ns/op             408 B/op          4 allocs/op
Benchmark_Background_WithLocalValue-8             131917              8846 ns/op             408 B/op          4 allocs/op
```

Release
```
go test -bench=. -benchmem -count=1 -cpu 8 -parallel 8 -tags release

goos: linux
goarch: amd64
pkg: github.com/wspowell/context
cpu: AMD Ryzen 9 4900HS with Radeon Graphics         
Benchmark_Background-8                          10086432               117.0 ns/op           104 B/op          3 allocs/op
Benchmark_golang_Background-8                   1000000000               0.3613 ns/op          0 B/op          0 allocs/op
Benchmark_Background_WithValue-8                 7274389               161.5 ns/op           152 B/op          4 allocs/op
Benchmark_golang_Background_WithValue-8         23939910                51.03 ns/op           48 B/op          1 allocs/op
Benchmark_Background_Value-8                    295331836                4.093 ns/op           0 B/op          0 allocs/op
Benchmark_golang_Background_Value-8             292899957                4.085 ns/op           0 B/op          0 allocs/op
Benchmark_Localized_Value-8                      2853351               421.0 ns/op           392 B/op          4 allocs/op
Benchmark_Background_WithLocalValue-8            3905462               294.0 ns/op           392 B/op          4 allocs/op
```