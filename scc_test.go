package graph_test

import (
	"slices"
	"testing"

	"gopkg.in/dnaeon/go-graph.v1"
)

// sccValueSets returns the components as sorted sets of int values, so that
// they can be compared independently of the (unspecified) component and
// within-component ordering.
func sccValueSets(sccs [][]*graph.Vertex[int]) [][]int {
	result := make([][]int, 0, len(sccs))
	for _, scc := range sccs {
		values := make([]int, 0, len(scc))
		for _, v := range scc {
			values = append(values, v.Value)
		}
		slices.Sort(values)
		result = append(result, values)
	}

	// Sort the components themselves for a stable comparison.
	slices.SortFunc(result, func(a, b []int) int {
		return slices.Compare(a, b)
	})

	return result
}

// nonTrivialSCCs returns only the components that represent a cyclic cluster:
// those with more than one vertex, or a single vertex with a self-loop.
func nonTrivialSCCs(t *testing.T, g graph.Graph[int], sccs [][]*graph.Vertex[int]) [][]*graph.Vertex[int] {
	t.Helper()

	result := make([][]*graph.Vertex[int], 0)
	for _, scc := range sccs {
		if len(scc) > 1 {
			result = append(result, scc)
			continue
		}
		v := scc[0]
		if g.EdgeExists(v.Value, v.Value) {
			result = append(result, scc)
		}
	}

	return result
}

func TestStronglyConnectedComponentsNotDirected(t *testing.T) {
	g := graph.New[int](graph.KindUndirected)
	g.AddEdge(1, 2)
	if _, err := graph.StronglyConnectedComponents(g); err != graph.ErrIsNotDirectedGraph {
		t.Fatalf("expected ErrIsNotDirectedGraph, got %v", err)
	}
}

func TestStronglyConnectedComponentsAcyclic(t *testing.T) {
	// A DAG spine plus two disjoint paths reconverging on a shared
	// successor. Every vertex must be its own SCC and there must be no
	// non-trivial component.
	g := graph.New[int](graph.KindDirected)
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)
	g.AddEdge(3, 4)
	g.AddEdge(1, 5)
	g.AddEdge(5, 4) // reconverge on 4: 1->2->3->4 and 1->5->4

	sccs, err := graph.StronglyConnectedComponents(g)
	if err != nil {
		t.Fatal(err)
	}

	if got := len(sccs); got != 5 {
		t.Fatalf("expected 5 singleton SCCs, got %d: %v", got, sccValueSets(sccs))
	}
	for _, scc := range sccs {
		if len(scc) != 1 {
			t.Fatalf("expected all SCCs to be singletons, got %v", sccValueSets(sccs))
		}
	}
	if n := len(nonTrivialSCCs(t, g, sccs)); n != 0 {
		t.Fatalf("expected no non-trivial SCC in an acyclic graph, got %d", n)
	}
}

func TestStronglyConnectedComponentsTwoCycle(t *testing.T) {
	// 1 <-> 2 is one non-trivial SCC; 3 is a trivial tail.
	g := graph.New[int](graph.KindDirected)
	g.AddEdge(1, 2)
	g.AddEdge(2, 1)
	g.AddEdge(2, 3)

	sccs, err := graph.StronglyConnectedComponents(g)
	if err != nil {
		t.Fatal(err)
	}

	nonTrivial := nonTrivialSCCs(t, g, sccs)
	if len(nonTrivial) != 1 {
		t.Fatalf("expected exactly one non-trivial SCC, got %v", sccValueSets(nonTrivial))
	}
	got := sccValueSets(nonTrivial)
	want := [][]int{{1, 2}}
	if !slices.EqualFunc(got, want, slices.Equal) {
		t.Fatalf("want non-trivial SCC %v, got %v", want, got)
	}
}

func TestStronglyConnectedComponentsThreeCycle(t *testing.T) {
	// 1 -> 2 -> 3 -> 1 is a single non-trivial SCC of {1, 2, 3}.
	g := graph.New[int](graph.KindDirected)
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)
	g.AddEdge(3, 1)

	sccs, err := graph.StronglyConnectedComponents(g)
	if err != nil {
		t.Fatal(err)
	}

	nonTrivial := nonTrivialSCCs(t, g, sccs)
	got := sccValueSets(nonTrivial)
	want := [][]int{{1, 2, 3}}
	if !slices.EqualFunc(got, want, slices.Equal) {
		t.Fatalf("want non-trivial SCC %v, got %v", want, got)
	}
}

// TestStronglyConnectedComponentsMixed models the cyclic-sample fixture: a DAG
// spine with a shared-successor diamond, plus two independent cycles. It must
// report exactly two non-trivial SCCs.
func TestStronglyConnectedComponentsMixed(t *testing.T) {
	// Vertices: 1..5 form a DAG spine with a shared-successor diamond;
	// 10,11 form a 2-cycle; 20,21,22 form a 3-cycle.
	g := graph.New[int](graph.KindDirected)

	// DAG spine + shared successor (diamond)
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 4)
	g.AddEdge(3, 4)  // shared successor
	g.AddEdge(4, 5)  // leaf

	// 2-cycle
	g.AddEdge(1, 10)
	g.AddEdge(10, 11)
	g.AddEdge(11, 10)

	// 3-cycle
	g.AddEdge(1, 20)
	g.AddEdge(20, 21)
	g.AddEdge(21, 22)
	g.AddEdge(22, 20)

	sccs, err := graph.StronglyConnectedComponents(g)
	if err != nil {
		t.Fatal(err)
	}

	nonTrivial := nonTrivialSCCs(t, g, sccs)
	got := sccValueSets(nonTrivial)
	want := [][]int{{10, 11}, {20, 21, 22}}
	if !slices.EqualFunc(got, want, slices.Equal) {
		t.Fatalf("want non-trivial SCCs %v, got %v", want, got)
	}
}

func TestStronglyConnectedComponentsSelfLoop(t *testing.T) {
	// A single vertex with a self-loop is a non-trivial (cyclic) SCC.
	g := graph.New[int](graph.KindDirected)
	g.AddEdge(1, 1)
	g.AddEdge(1, 2)

	sccs, err := graph.StronglyConnectedComponents(g)
	if err != nil {
		t.Fatal(err)
	}

	nonTrivial := nonTrivialSCCs(t, g, sccs)
	got := sccValueSets(nonTrivial)
	want := [][]int{{1}}
	if !slices.EqualFunc(got, want, slices.Equal) {
		t.Fatalf("want non-trivial SCC %v, got %v", want, got)
	}
}
