// +build !release

package context

// type Localized struct {
// 	context.Context

// 	// Store locals in a map that do not have a defined variable.
// 	locals map[interface{}]interface{}

// 	goroutineOrigin goroutineId
// }

// func NewLocalized() *Localized {
// 	return FromContext(context.Background())
// }

// // FromContext created a new Localized context using the given parent context.
// func FromContext(ctx context.Context) *Localized {
// 	// If the context is a contextualizer, then use its Context() value instead.
// 	// This prevents Localizer value from being copied across goroutines.
// 	if contextualizerCtx, ok := ctx.(contextualizer); ok {
// 		ctx = contextualizerCtx.GetContext()
// 	}

// 	return &Localized{
// 		Context:         ctx,
// 		locals:          map[interface{}]interface{}{},
// 		goroutineOrigin: curID(),
// 	}
// }

// // Value returns the value stored at key in the context.
// // First check local values, then checks stored context.
// // Returns nil if key does not exist.
// func (self *Localized) Value(key interface{}) interface{} {
// 	self.threadSafetyCheck()
// 	if localValue, exists := self.locals[key]; exists {
// 		return localValue
// 	}
// 	return self.Context.Value(key)
// }

// func (self *Localized) Localize(key interface{}, value interface{}) {
// 	self.threadSafetyCheck()
// 	self.locals[key] = value
// }

// func (self *Localized) threadSafetyCheck() {
// 	if !self.goroutineOrigin.isSameGoroutine() {
// 		// Panic instead of error because what are you doing developer???
// 		panic("local context used outside original goroutine")
// 	}
// }

// func (self *Localized) GetContext() context.Context {
// 	return self.Context
// }

// func (self *Localized) setContext(ctx context.Context) {
// 	self.Context = ctx
// }

// // Lock so `go vet` gives a warning if this struct is copied.
// func (*Localized) Lock() {}

// // Unlock so `go vet` gives a warning if this struct is copied.
// func (*Localized) Unlock() {}
