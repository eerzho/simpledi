package simpledi

import (
	"errors"
	"sync"
)

type Definition struct {
	ID        string
	DependsOn []string
	New       func() any
	Close     func() error
}

type container struct {
	resolved    bool
	definitions []Definition
	instances   map[string]any
}

var ctr = sync.OnceValue(func() *container {
	return &container{
		definitions: make([]Definition, 0),
		instances:   make(map[string]any),
	}
})

// ResetForTesting сбрасывает контейнер для целей тестирования
// ВАЖНО: Используйте только в тестах!
func ResetForTesting() {
	c := ctr()
	c.definitions = make([]Definition, 0)
	c.instances = make(map[string]any)
	c.resolved = false
}

// GetDefinitionsForTesting возвращает список definitions в отсортированном порядке
// ВАЖНО: Используйте только в тестах!
func GetDefinitionsForTesting() []Definition {
	return ctr().definitions
}

func Set(d Definition) {
	if d.ID == "" {
		panic("empty id")
	}
	if d.New == nil {
		panic("empty new")
	}
	if ctr().resolved {
		panic("container resolved")
	}
	ctr().definitions = append(ctr().definitions, d)
}

func Get[T any](id string) T {
	if !ctr().resolved {
		panic("container not resolved")
	}
	object, ok := ctr().instances[id]
	if !ok {
		panic("not found key")
	}
	typedObject, ok := object.(T)
	if !ok {
		panic("invalid type")
	}
	return typedObject
}

func Close() error {
	if !ctr().resolved {
		return errors.New("container not resolved")
	}

	errs := make([]error, 0)
	for i := len(ctr().definitions) - 1; i >= 0; i-- {
		definition := ctr().definitions[i]
		if definition.Close != nil {
			if err := definition.Close(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func Resolve() {
	if ctr().resolved {
		panic("container resolved")
	}

	err := ctr().sort()
	if err != nil {
		panic(err)
	}

	for _, definition := range ctr().definitions {
		instance := definition.New()
		ctr().instances[definition.ID] = instance
	}
	ctr().resolved = true
}

func (c *container) sort() error {
	queue := make([]Definition, 0, len(c.definitions))

	indegree := make(map[string]int)
	graph := make(map[string][]Definition)
	for _, definition := range c.definitions {
		count := len(definition.DependsOn)
		if count == 0 {
			queue = append(queue, definition)
			continue
		}
		indegree[definition.ID] = count
		for _, depend := range definition.DependsOn {
			graph[depend] = append(graph[depend], definition)
		}
	}

	definitions := make([]Definition, 0, len(c.definitions))
	for len(queue) > 0 {
		definition := queue[0]
		queue = queue[1:]
		definitions = append(definitions, definition)
		for _, subdefinition := range graph[definition.ID] {
			indegree[subdefinition.ID]--
			if indegree[subdefinition.ID] == 0 {
				queue = append(queue, subdefinition)
			}
		}
	}

	if len(c.definitions) != len(definitions) {
		return errors.New("cyclic depends")
	}

	c.definitions = definitions
	return nil
}
