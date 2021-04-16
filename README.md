# local

Localized Golang Context

## Context

Golang `context.Context` is a feature that is easily abused. A `context.Context` should only be used for immutable data and are meant to be passed between API boundaries (and therefore must be thread safe). However, it is extremely tempting (and easy) to violate this contract and use it as a generic variable store for values used throughout an goroutines lifetime. 

`context.Context` attempts to address these issues. A `context.Context` is both a `context.Context` and a variable store for goroutine local data. The difference is that `context.Context` provides behavior to localize data to the goroutine. Localized data is not immutable and must never be sent across API boundaries (and therefore not thread safe). Passing a local context to a goroutine will automatically cut out the local data and only pass the immutable context data.

## Building

The package utilizes goroutine identification (that Golang authors created) to catch threading issues during development. The downside is that tracking goroutine IDs adds significant overhead. In order to solve this, the default build will use goroutine tracking. To build a "release" build, use -tags option with "release".
`go build -tags release ./...`

## Example

```
// Create a new localized context.
ctx := context.NewLocalized()

// Add immutable data.
context.WithValue(ctx, "log_level", "info")

// Add local data.
// Stored map may be accessed and altered at any time during the goroutine.
ctx.Localize("local", map[string]string{})

// Start a goroutine to process data.
go processData(ctx)

...

func processData(boundaryCtx context.Context) {
    // Create a new context local to this goroutine.
    // Context no longer has access to the "local" key.
    ctx := context.FromContext(boundaryCtx)

    // Get the log level to create a new logger.
    logLevel := ctx.Value("log_level")

    ...
}
```