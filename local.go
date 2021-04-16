package context

type locals map[interface{}]interface{}
type contextKey int

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
