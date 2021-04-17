// +build !release

package context

import "reflect"

type locals map[interface{}]interface{}
type localsKey struct{}

// Localize a Context to the current goroutine.
// Any local values set on the Context via WithLocalValue become inaccessable to the returned Context.
func Localize(ctx Context) Context {
	var localValues locals

	if local, ok := ctx.Value(localsKey{}).(*localized); ok {
		if local.goroutineOrigin.isSameGoroutine() {
			panic("context localized twice in the same goroutine")
		}

		// Values are shadowed by the local context to prevent access to any locals
		// in a parent context.
		localValues = make(locals, len(local.localValues))
		for key, value := range local.localValues {
			// Anything that can Clone() should be cloned.
			// Cloned values must be thread safe.
			if cloner, ok := value.(interface{ Clone() interface{} }); ok {
				localValues[key] = cloner.Clone()
			} else {
				switch reflect.TypeOf(value).Kind() {
				case reflect.Array, reflect.Chan, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.Struct, reflect.UnsafePointer:
					// All shadowed local values reset to nil.
					localValues[key] = nil
				default:
					// Value should be copyable.
					localValues[key] = value
				}
			}
		}
	} else {
		localValues = locals{}
	}

	return &localized{
		Context:         ctx,
		localValues:     localValues,
		goroutineOrigin: curID(),
	}
}

// WithLocalValue wraps the parent Context and adds the key-value pair
// as a value local to the current goroutine.
func WithLocalValue(parent Context, key interface{}, value interface{}) Context {
	if local, ok := parent.Value(localsKey{}).(*localized); ok {
		if !local.goroutineOrigin.isSameGoroutine() {
			panic("context not localized to the current goroutine")
		}

		local.localValues[key] = value
		return parent
	}

	return WithLocalValue(Localize(parent), key, value)
}
