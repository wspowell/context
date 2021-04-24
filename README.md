# context 

Replacement for Golang Context

## Context

Golang `context.Context` is a feature that is easily abused. A `context.Context` should only be used for immutable data and are meant to be passed between API boundaries (and therefore must be thread safe). However, it is extremely tempting (and easy) to violate this contract and use it as a generic variable store for values used throughout an goroutines lifetime. 

`context.Context` attempts to address these issues. A `context.Context` is both a `context.Context` and a variable store for goroutine local data. The difference is that `context.Context` provides behavior to localize data to the goroutine. Localized data is not thread safe and must never be sent across API boundaries. Localizing a context to a goroutine will cut out the local data and only allow access to the immutable context data. If localized data implements `Localize() interface{}`, then the value will be cloned in the localized context. `Localize() interface{}` must return a thread safe value.

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
func (self Log) Localize() interface{} {
    ...
}

func processData(ctx context.Context) {
    // Create a new context local to this goroutine.
    // Context no longer has access to the localized keys.
    // Values are shadowed as nil, unless values implement Localize() interface{}.
    ctx := context.Localize(ctx)
    
    ...
}
```
