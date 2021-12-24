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
Benchmark_Background-8                            313935              3782 ns/op              64 B/op          1 allocs/op
Benchmark_golang_Background-8                   1000000000               0.4807 ns/op          0 B/op          0 allocs/op
Benchmark_Background_WithValue-8                  302839              3820 ns/op             112 B/op          2 allocs/op
Benchmark_golang_Background_WithValue-8         24433995                47.67 ns/op           48 B/op          1 allocs/op
Benchmark_Background_Value-8                    225462111                5.275 ns/op           0 B/op          0 allocs/op
Benchmark_golang_Background_Value-8             277593678                4.358 ns/op           0 B/op          0 allocs/op
Benchmark_Localized_Value-8                        72796             16625 ns/op             440 B/op          9 allocs/op
Benchmark_Background_WithLocalValue-8             138301              8642 ns/op             376 B/op          7 allocs/op
```

Release
```
go test -bench=. -benchmem -count=1 -cpu 8 -parallel 8 -tags release

goos: linux
goarch: amd64
pkg: github.com/wspowell/context
cpu: AMD Ryzen 9 4900HS with Radeon Graphics         
Benchmark_Background-8                          26362558                47.00 ns/op           64 B/op          1 allocs/op
Benchmark_golang_Background-8                   1000000000               0.4766 ns/op          0 B/op          0 allocs/op
Benchmark_Background_WithValue-8                13373020                94.64 ns/op          112 B/op          2 allocs/op
Benchmark_golang_Background_WithValue-8         23825217                47.44 ns/op           48 B/op          1 allocs/op
Benchmark_Background_Value-8                    278421498                4.357 ns/op           0 B/op          0 allocs/op
Benchmark_golang_Background_Value-8             273102477                4.369 ns/op           0 B/op          0 allocs/op
Benchmark_Localized_Value-8                      6209235               202.2 ns/op           440 B/op          9 allocs/op
Benchmark_Background_WithLocalValue-8            2752718               480.5 ns/op           376 B/op          7 allocs/op
```
