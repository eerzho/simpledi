package simpledi

type Container struct {
	ts       *topoSort
	objects  map[string]any
	builders map[string]func() any
}

func NewContainer() *Container {
	return &Container{
		ts:       newTopoSort(),
		objects:  make(map[string]any),
		builders: make(map[string]func() any),
	}
}

func (c *Container) Register(key string, needs []string, builder func() any) *Container {
	if builder == nil {
		panic("builder is nil")
	}
	if err := c.ts.append(key, needs); err != nil {
		panic(err)
	}
	c.builders[key] = builder
	return c
}

func (c *Container) Get(key string) any {
	if object, ok := c.objects[key]; ok {
		return object
	}
	panic("object not found: " + key)
}

func (c *Container) Resolve() {
	sorted, err := c.ts.sort()
	if err != nil {
		panic(err)
	}
	for _, key := range sorted {
		c.objects[key] = c.builders[key]()
	}
}
