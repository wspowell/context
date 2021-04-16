package context

// // contextualizer is a context.Context that allows the underlying immutable
// //   context.Context to be accessed and overridden.
// type contextualizer interface {
// 	// Context embedded behavior.
// 	context.Context

// 	// setContext sets the underlying context.Context.
// 	setContext(context.Context)

// 	// GetContext returns the underlying context.Context.
// 	// Returned value must be thread safe.
// 	GetContext() context.Context
// }
