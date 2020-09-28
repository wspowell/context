package local

import (
	"context"
	"time"
)

// Context is a contextualizer and a Localizer.
// Context is meant to be used as a replacement for context.Context.
// It is expected that the value of the internal context.Context is thread safe.
type Context interface {
	contextualizer
	Localizer
}

// WithValue creates a new context with the given value.
// See: context.WithValue()
func WithValue(parent Context, key interface{}, value interface{}) {
	childContext := context.WithValue(parent.GetContext(), key, value)

	parent.setContext(childContext)
}

// WithCancel creates a new context that can be cancelled.
// See: context.WithCancel()
func WithCancel(parent Context) context.CancelFunc {
	childContext, cancel := context.WithCancel(parent.GetContext())

	parent.setContext(childContext)

	return cancel
}

// WithDeadline creates a new context that has a deadline.
// See: context.WithDeadline()
func WithDeadline(parent Context, deadline time.Time) context.CancelFunc {
	childContext, cancel := context.WithDeadline(parent.GetContext(), deadline)

	parent.setContext(childContext)

	return cancel
}

// WithTimeout creates a new context that has a timeout.
// See: context.WithTimeout()
func WithTimeout(parent Context, timeout time.Duration) context.CancelFunc {
	childContext, cancel := context.WithTimeout(parent.GetContext(), timeout)

	parent.setContext(childContext)

	return cancel
}
