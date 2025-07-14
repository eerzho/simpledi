package simpledi

import "sync"

var (
	c    *Container
	once sync.Once
)

// Register a dependency by key
//   - key:     unique name for the dependency
//   - needs:   list of dependency keys this object depends on
//   - builder: function that returns the object instance
func Register(key string, needs []string, builder func() any) {
	defaultC().Register(key, needs, builder)
}

// Get a dependency by key
//   - key: unique name of the dependency
func Get(key string) (any, bool) {
	return defaultC().Get(key)
}

// MustGet retrieves a dependency by key or panics if not found
//   - key: unique name of the dependency
func MustGet(key string) any {
	return defaultC().MustGet(key)
}

// Resolve all dependencies
func Resolve() error {
	return defaultC().Resolve()
}

func defaultC() *Container {
	once.Do(func() {
		c = NewContainer()
	})
	return c
}
