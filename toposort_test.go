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
	err = graph.WalkTopoOrder(g3, dummyWalker)
	if !errors.Is(err, graph.ErrCycleDetected) {
		t.Fatal("g3: graph should contain a cycle")
	}

	// The returned error should be a *CycleError carrying the
	// exact cycle.
	var cycleErr *graph.CycleError[int]
	if !errors.As(err, &cycleErr) {
		t.Fatalf("g3: expected a *graph.CycleError, got %T", err)
	}
	assertValidCycle(t, g3, cycleErr.Cycle)
}

// assertValidCycle verifies that the given slice is a well-formed cycle
// in g: it starts and ends on the same vertex and every consecutive
// pair is a real edge.
func assertValidCycle(t *testing.T, g graph.Graph[int], cycle []*graph.Vertex[int]) {
	t.Helper()

	if len(cycle) < 2 {
		t.Fatalf("cycle too short: %v", cycle)
	}

	if cycle[0].Value != cycle[len(cycle)-1].Value {
		t.Fatalf("cycle must start and end on the same vertex, got %v", cycleValues(cycle))
	}

	for i := 0; i < len(cycle)-1; i++ {
		from := cycle[i].Value
		to := cycle[i+1].Value
		if !g.EdgeExists(from, to) {
			t.Fatalf("cycle contains non-existent edge %d -> %d (cycle %v)", from, to, cycleValues(cycle))
		}
	}
}

func cycleValues(cycle []*graph.Vertex[int]) []int {
	values := make([]int, 0, len(cycle))
	for _, v := range cycle {
		values = append(values, v.Value)
	}

	return values
}

// TestWalkTopoOrderSharedSuccessor is a regression test for a
// false-positive cycle detection bug. A shared successor reachable via
// more than one path must not be reported as a cycle. We build the graph
// with both neighbour orderings, because the bug was order-dependent.
func TestWalkTopoOrderSharedSuccessor(t *testing.T) {
	dummyWalker := func(v *graph.Vertex[string]) error {
		return nil
	}

	// Edges: P->L, P->M, M->L. L is reached from P both directly and via
	// M (a transitive edge over the P->M->L path), so it has two incoming
	// paths but the graph is acyclic.
	build := func(lFirst bool) graph.Graph[string] {
		g := graph.New[string](graph.KindDirected)
		if lFirst {
			g.AddEdge("P", "L")
			g.AddEdge("P", "M")
		} else {
			g.AddEdge("P", "M")
			g.AddEdge("P", "L")
		}
		g.AddEdge("M", "L")
		return g
	}

	for _, lFirst := range []bool{true, false} {
		g := build(lFirst)
		if err := graph.WalkTopoOrder(g, dummyWalker); err != nil {
			t.Fatalf("lFirst=%v: unexpected error on acyclic graph: %v", lFirst, err)
		}
	}
}

// TestWalkTopoOrderMultiPathDAG verifies that a larger acyclic DAG
// where several vertices are reachable by more than one path is walked
// without a false cycle and that every vertex is emitted exactly once.
func TestWalkTopoOrderMultiPathDAG(t *testing.T) {
	// 1 -> {2,3}; 2 -> {4,5}; 3 -> {4,5}; 4 -> 6; 5 -> 6
	g := graph.New[int](graph.KindDirected)
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 4)
	g.AddEdge(2, 5)
	g.AddEdge(3, 4)
	g.AddEdge(3, 5)
	g.AddEdge(4, 6)
	g.AddEdge(5, 6)

	collector := g.NewCollector()
	if err := graph.WalkTopoOrder(g, collector.WalkFunc); err != nil {
		t.Fatalf("unexpected error on acyclic graph: %v", err)
	}

	got := cycleValues(collector.Get())
	slices.Sort(got)
	want := []int{1, 2, 3, 4, 5, 6}
	if !slices.Equal(got, want) {
		t.Fatalf("each vertex should be emitted exactly once, want %v, got %v", want, got)
	}
}

// TestWalkTopoOrderTwoCycle verifies that a genuine two-node cycle is
// still detected after the false-positive fix, guarding against
// over-correcting into a false negative.
func TestWalkTopoOrderTwoCycle(t *testing.T) {
	dummyWalker := func(v *graph.Vertex[int]) error {
		return nil
	}

	// 1 -> 2 and 2 -> 1 form a mutual dependency (2-cycle).
	g := graph.New[int](graph.KindDirected)
	g.AddEdge(1, 2)
	g.AddEdge(2, 1)

	err := graph.WalkTopoOrder(g, dummyWalker)
	if !errors.Is(err, graph.ErrCycleDetected) {
		t.Fatalf("expected a cycle to be detected, got %v", err)
	}

	var cycleErr *graph.CycleError[int]
	if !errors.As(err, &cycleErr) {
		t.Fatalf("expected a *graph.CycleError, got %T", err)
	}
	assertValidCycle(t, g, cycleErr.Cycle)
}

func TestWalkCyclesNotDirected(t *testing.T) {
	noop := func(cycle []*graph.Vertex[int]) error {
		return nil
	}

	g := graph.New[int](graph.KindUndirected)
	g.AddEdge(1, 2)
	if err := graph.WalkCycles(g, noop); err != graph.ErrIsNotDirectedGraph {
		t.Fatalf("expected ErrIsNotDirectedGraph, got %v", err)
	}
}

func TestWalkCyclesAcyclic(t *testing.T) {
	// A DAG with a shared-successor diamond: the walker must never be
	// invoked.
	g := graph.New[int](graph.KindDirected)
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 4)
	g.AddEdge(3, 4)

	count := 0
	walker := func(cycle []*graph.Vertex[int]) error {
		count++
		return nil
	}
	if err := graph.WalkCycles(g, walker); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("expected no cycles in an acyclic graph, walker invoked %d times", count)
	}
}

// TestWalkCyclesMixed models the cyclic-sample fixture and asserts that
// WalkCycles emits exactly one valid cycle per cyclic cluster (two clusters).
func TestWalkCyclesMixed(t *testing.T) {
	g := graph.New[int](graph.KindDirected)

	// DAG spine + shared successor (diamond)
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 4)
	g.AddEdge(3, 4)
	g.AddEdge(4, 5)

	// 2-cycle: 10 <-> 11
	g.AddEdge(1, 10)
	g.AddEdge(10, 11)
	g.AddEdge(11, 10)

	// 3-cycle: 20 -> 21 -> 22 -> 20
	g.AddEdge(1, 20)
	g.AddEdge(20, 21)
	g.AddEdge(21, 22)
	g.AddEdge(22, 20)

	cycles := make([][]*graph.Vertex[int], 0)
	walker := func(cycle []*graph.Vertex[int]) error {
		cycles = append(cycles, cycle)
		return nil
	}
	if err := graph.WalkCycles(g, walker); err != nil {
		t.Fatal(err)
	}

	if len(cycles) != 2 {
		t.Fatalf("expected exactly 2 cycles, got %d", len(cycles))
	}

	// Each emitted cycle must be a well-formed, closed path of real edges,
	// and every vertex in it must belong to the same expected cluster.
	twoCycle := map[int]bool{10: true, 11: true}
	threeCycle := map[int]bool{20: true, 21: true, 22: true}
	for _, cycle := range cycles {
		assertValidCycle(t, g, cycle)

		values := cycleValues(cycle)
		inTwo, inThree := true, true
		for _, v := range values {
			if !twoCycle[v] {
				inTwo = false
			}
			if !threeCycle[v] {
				inThree = false
			}
		}
		if !inTwo && !inThree {
			t.Fatalf("emitted cycle %v does not belong to a single expected cluster", values)
		}
	}
}

func TestWalkCyclesStopWalking(t *testing.T) {
	// Two independent cycles; stop after the first one is seen.
	g := graph.New[int](graph.KindDirected)
	g.AddEdge(1, 2)
	g.AddEdge(2, 1)
	g.AddEdge(3, 4)
	g.AddEdge(4, 3)

	count := 0
	walker := func(cycle []*graph.Vertex[int]) error {
		count++
		return graph.ErrStopWalking
	}
	if err := graph.WalkCycles(g, walker); err != nil {
		t.Fatalf("ErrStopWalking should not be propagated, got %v", err)
	}
	if count != 1 {
		t.Fatalf("expected walking to stop after the first cycle, got %d", count)
	}
}

func TestWalkCyclesWalkerError(t *testing.T) {
	g := graph.New[int](graph.KindDirected)
	g.AddEdge(1, 2)
	g.AddEdge(2, 1)

	myErr := errors.New("my custom error")
	walker := func(cycle []*graph.Vertex[int]) error {
		return myErr
	}
	if err := graph.WalkCycles(g, walker); err != myErr {
		t.Fatalf("expected custom walker error to propagate, got %v", err)
	}
}

func TestWalkCyclesSelfLoop(t *testing.T) {
	// A self-loop is a cycle [v v].
	g := graph.New[int](graph.KindDirected)
	g.AddEdge(1, 1)
	g.AddEdge(1, 2)

	cycles := make([][]*graph.Vertex[int], 0)
	walker := func(cycle []*graph.Vertex[int]) error {
		cycles = append(cycles, cycle)
		return nil
	}
	if err := graph.WalkCycles(g, walker); err != nil {
		t.Fatal(err)
	}

	if len(cycles) != 1 {
		t.Fatalf("expected exactly one cycle for a self-loop, got %d", len(cycles))
	}
	assertValidCycle(t, g, cycles[0])
	if got := cycleValues(cycles[0]); !slices.Equal(got, []int{1, 1}) {
		t.Fatalf("expected self-loop cycle [1 1], got %v", got)
	}
}
