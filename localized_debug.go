// +build !release

package context

type localized struct {
	Context

	key interface{}

	goroutineOrigin goroutineId
}

// Value returns the value stored at key in the context.
// First check local values, then checks stored context.
// Returns nil if key does not exist.
func (self *localized) Value(key interface{}) interface{} {
	if !self.goroutineOrigin.isSameGoroutine() {
		panic("localized value accessed outside original goroutine")
	}

	if localValues, ok := self.Context.Value(localsKey).(locals); ok {
		if localValue, exists := localValues[key]; exists {
			return localValue
		}
	}

	return self.Context.Value(key)
}
