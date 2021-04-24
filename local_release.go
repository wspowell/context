// +build release

package context

type locals map[interface{}]interface{}
type localsKey struct{}

// Localize a Context to the current goroutine.
// Any local values set on the Context via WithLocalValue become inaccessable to the returned Context.
func Localize(ctx Context) Context {
	var localValues locals

	if local, ok := ctx.Value(localsKey{}).(*localized); ok {
		// Values are shadowed by the local context to prevent access to any locals in a parent context.
		localValues = make(locals, len(local.localValues))
		for key, value := range local.localValues {
			if localizer, ok := value.(interface{ Localize() interface{} }); ok {
				// Use localized value.
				localValues[key] = localizer.Localize()
			} else {
				// Shadowed local value reset to nil.
				localValues[key] = nil
			}
		}
	} else {
		localValues = locals{}
	}

	return &localized{
		Context:     ctx,
		localValues: localValues,
	}
}

// WithLocalValue wraps the parent Context and adds the key-value pair
// as a value local to the current goroutine.
func WithLocalValue(parent Context, key interface{}, value interface{}) {
	if local, ok := parent.Value(localsKey{}).(*localized); ok {
		local.localValues[key] = value
		return
	}

	panic("context not localized to the current goroutine")
}
