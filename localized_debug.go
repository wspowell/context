//go:build !release
// +build !release

package context

type localized struct {
	Context

	localValues locals

	goroutineOrigin goroutineId
}

// Value returns the value stored at key in the context.
// First check local values, then checks stored context.
// Returns nil if key does not exist.
func (self *localized) Value(key interface{}) interface{} {
	if key == (localsKey{}) {
		return self
	}

	if localValue, exists := self.localValues[key]; exists {
		if !self.goroutineOrigin.isSameGoroutine() {
			panic("localized value accessed outside original goroutine")
		}
		return localValue
	}

	return self.Context.Value(key)
}
