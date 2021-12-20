//go:build release
// +build release

package context

import (
	"sync"
)

type localCtx struct {
	Context

	localsMutex *sync.RWMutex
	localValues locals
}

// Value returns the value stored at key in the context.
// First check local values, then checks stored context.
// Returns nil if key does not exist.
func (self *localCtx) Value(key any) any {
	if key == (localsKey{}) {
		return self
	}

	self.localsMutex.RLock()
	localValue, exists := self.localValues[key]
	self.localsMutex.RUnlock()
	if exists {
		return localValue
	}

	return self.Context.Value(key)
}
