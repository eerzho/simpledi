package simpledi

import (
	"fmt"
)

type topoSort struct {
	adjacency map[string][]string
}

func newTopoSort() *topoSort {
	return &topoSort{
		adjacency: make(map[string][]string),
	}
}

func (ts *topoSort) append(node string, neighbours []string) error {
	if _, ok := ts.adjacency[node]; ok {
		return fmt.Errorf("node already exists")
	}
	ts.adjacency[node] = neighbours
	return nil
}

func (ts *topoSort) sort() ([]string, error) {
	queue := make([]string, 0)
	graph := make(map[string][]string)
	inDegree := make(map[string]int)
	for node, neighbours := range ts.adjacency {
		if len(neighbours) == 0 {
			queue = append(queue, node)
		} else {
			inDegree[node] = len(neighbours)
		}
		for _, neighbour := range neighbours {
			if _, ok := ts.adjacency[neighbour]; !ok {
				return nil, fmt.Errorf("node [%s] not declared", neighbour)
			}
			graph[neighbour] = append(graph[neighbour], node)
		}
	}
	sorted := make([]string, 0)
	for len(queue) != 0 {
		key := queue[0]
		queue = queue[1:]
		sorted = append(sorted, key)
		for _, neighbour := range graph[key] {
			inDegree[neighbour]--
			if inDegree[neighbour] == 0 {
				queue = append(queue, neighbour)
			}
		}
	}
	if len(sorted) != len(ts.adjacency) {
		cycles := make([]string, 0)
		for node, degree := range inDegree {
			if degree > 0 {
				cycles = append(cycles, node)
			}
		}
		return nil, fmt.Errorf("cyclic detected: %v", cycles)
	}
	return sorted, nil
}
