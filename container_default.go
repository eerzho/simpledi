package simpledi

// Default container
var c = NewContainer()

// Register a dependency by key
//   - key:     unique name for the dependency
//   - needs:   list of dependency keys this object depends on
//   - builder: function that returns the object instance
func Register(key string, needs []string, builder func() any) {
	c.Register(key, needs, builder)
}

// Get a dependency by key
//   - key: unique name of the dependency
func Get(key string) any {
	return c.Get(key)
}

// Resolve all dependencies
func Resolve() error {
	return c.Resolve()
}
