//go:build !release
// +build !release

package context

import "sync"

type locals map[any]any
type localsKey struct{}

// Localize a Context to the current goroutine.
// Any local values set on the Context via WithLocalValue become inaccessible to the returned Context.
func Localize(ctx Context) Context {
	var localValues locals

	if local, ok := ctx.Value(localsKey{}).(*localCtx); ok {
		if local.goroutineOrigin.isSameGoroutine() {
			panic("context localized twice in the same goroutine")
		}

		// Values are shadowed by the local context to prevent access to any locals in a parent context.
		local.localsMutex.Lock()
		localValues = make(locals, len(local.localValues))
		for key, value := range local.localValues {
			if localizer, ok := value.(interface{ Localize() any }); ok {
				// Use localized value.
				localValues[key] = localizer.Localize()
			} else {
				// Shadowed local value reset to nil.
				localValues[key] = nil
			}
		}
		local.localsMutex.Unlock()
	} else {
		localValues = make(locals, 0)
	}

	return &localCtx{
		Context:         ctx,
		localsMutex:     &sync.RWMutex{},
		localValues:     localValues,
		goroutineOrigin: curID(),
	}
}

// WithLocalValue wraps the parent Context and adds the key-value pair
// as a value local to the current goroutine.
func WithLocalValue(parent Context, key any, value any) {
	if local, ok := parent.Value(localsKey{}).(*localCtx); ok {
		if !local.goroutineOrigin.isSameGoroutine() {
			panic("context not localized to the current goroutine")
		}

		local.localValues[key] = value

		return
	}

	panic("context not localized to the current goroutine")
}
