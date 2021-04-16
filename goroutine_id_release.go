// +build release

package context

// ID is the ID number of a goroutine.
type goroutineId uint64

func curID() goroutineId {
	return goroutineId(0)
}

func (self goroutineId) isSameGoroutine() bool {
	return true
}
