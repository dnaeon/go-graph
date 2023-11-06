package graph_test

import (
	"errors"
	"slices"
	"testing"

	"gopkg.in/dnaeon/go-graph.v1"
)

func TestWalkTopoOrder(t *testing.T) {
	dummyWalker := func(v *graph.Vertex[int]) error {
		return nil
	}

	// Topo sorting of undirected graphs is not support
	g1 := graph.New[int](graph.KindUndirected)
	err := graph.WalkTopoOrder(g1, dummyWalker)
	if err != graph.ErrIsNotDirectedGraph {
		t.Fatal("WalkTopoOrder: topo sort should fail on undirected graphs")
	}

	// Test topo sorting a simple acylic graph
	g2 := graph.New[int](graph.KindDirected)
	g2.AddEdge(1, 2)
	g2.AddEdge(2, 3)
	g2.AddEdge(3, 4)
	collector := g2.NewCollector()
	if err := graph.WalkTopoOrder(g2, collector.WalkFunc); err != nil {
		t.Fatal(err)
	}
	gotValues := make([]int, 0)
	for _, v := range collector.Get() {
		gotValues = append(gotValues, v.Value)
	}
	wantValues := []int{4, 3, 2, 1}

	if !slices.Equal(gotValues, wantValues) {
		t.Fatalf("g2: want topo order %v, got %v", wantValues, gotValues)
	}

	// Short-circuit collecting by signalling ErrStopWalking
	result := make([]int, 0)
	shortCircuitWalker := func(v *graph.Vertex[int]) error {
		if v.Value == 3 {
			return graph.ErrStopWalking
		}

		result = append(result, v.Value)
		return nil
	}
	if err := graph.WalkTopoOrder(g2, shortCircuitWalker); err != nil {
		t.Fatal(err)
	}

	// We should have collected only a single vertex so far
	if !slices.Equal([]int{4}, result) {
		t.Fatal("g2: collected vertex values do not match")
	}

	// Test with a walker which signals an error
	myErr := errors.New("my custom error")
	errWalker := func(v *graph.Vertex[int]) error {
		return myErr
	}
	if err := graph.WalkTopoOrder(g2, errWalker); err != myErr {
		t.Fatal("g2: walker did not return correct error")
	}

	// Test topo order with a graph containing a cycle
	g3 := graph.New[int](graph.KindDirected)
	g3.AddEdge(1, 2)
	g3.AddEdge(2, 3)
	g3.AddEdge(3, 4)
	g3.AddEdge(4, 1) // Cycle
	if err := graph.WalkTopoOrder(g3, dummyWalker); err != graph.ErrCycleDetected {
		t.Fatal("g3: graph should contain a cycle")
	}
}
