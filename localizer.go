package local

// Localizer stores data local to a goroutine. A localized context.
// This works differently than context.Context in that it is not meant to
//   cross API boundaries and is not immutable. However, it is designed to
//   be able to work alongside context.Context. It is also meant to be
//   wrapped by developers to allow for direct access of endpoint local data.
// Not thread safe.
type Localizer interface {
	// Do not allow copy.
	Lock()
	Unlock()

	// Localize a key/value pair to the current context scope.
	// Not thread safe.
	Localize(key interface{}, value interface{})
}

var _ Context = (*Localized)(nil)
var _ contextualizer = (*Localized)(nil)
