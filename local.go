package context

type locals map[interface{}]interface{}
type contextKey struct{}

var localsKey contextKey

func WithLocalValue(parent Context, key interface{}, value interface{}) Context {
	if localValues, ok := parent.Value(localsKey).(locals); ok {
		localValues[key] = value

		return &localized{
			Context: parent,
			key:     key,
			//goroutineOrigin: curID(),
		}
	}

	return WithLocalValue(Localize(parent), key, value)
}

func Localize(ctx Context) Context {
	if localValues, ok := ctx.Value(localsKey).(locals); ok {
		// Values are shadowed by the local context to prevent access to any locals
		// in a parent context.
		shadowValues := make(locals, len(localValues))
		for key := range localValues {
			// All shadowed local values reset to nil.
			shadowValues[key] = nil
		}

		return WithValue(ctx, localsKey, shadowValues)
	}

	return WithValue(ctx, localsKey, locals{})
}

type localized struct {
	Context

	key interface{}

	//goroutineOrigin goroutineId
}

// Value returns the value stored at key in the context.
// First check local values, then checks stored context.
// Returns nil if key does not exist.
func (self *localized) Value(key interface{}) interface{} {
	// if !self.goroutineOrigin.isSameGoroutine() {
	// 	panic("localized value accessed outside original goroutine")
	// }

	if localValues, ok := self.Context.Value(localsKey).(locals); ok {
		if localValue, exists := localValues[key]; exists {
			return localValue
		}
	}

	return self.Context.Value(key)
}
