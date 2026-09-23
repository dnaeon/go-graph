package graph

import (
	"errors"
	"slices"
	"strings"

	"gopkg.in/dnaeon/go-deque.v1"
)

// ErrCycleDetected is returned whenever a cycle has been detected in
// the graph.
var ErrCycleDetected = errors.New("cycle detected")

// ErrIsNotDirectedGraph is returned whenever an operation cannot be
// performed, because the graph is not directed.
var ErrIsNotDirectedGraph = errors.New("graph is not directed")

// CycleError is returned whenever a cycle has been detected during a
// topological sort. It wraps [ErrCycleDetected] and carries the ordered list of
// vertices forming the cycle.
//
// Callers can extract the cycle using [errors.As], e.g.
//
//	var cycleErr *graph.CycleError[int]
//	if errors.As(err, &cycleErr) {
//		fmt.Println(cycleErr.Cycle)
//	}
type CycleError[T comparable] struct {
	// Cycle is the ordered set of vertices forming the cycle. The slice
	// starts and ends with the same vertex, e.g. for the edges 1->2->3->1
	// the cycle is [1 2 3 1].
	Cycle []*Vertex[T]
}

// Error implements the error interface.
func (e *CycleError[T]) Error() string {
	labels := make([]string, 0, len(e.Cycle))
	for _, v := range e.Cycle {
		labels = append(labels, v.Label)
	}

	return "cycle detected: " + strings.Join(labels, " -> ")
}

// Unwrap returns the wrapped [ErrCycleDetected] error, allowing callers to use
// errors.Is(err, ErrCycleDetected).
func (e *CycleError[T]) Unwrap() error {
	return ErrCycleDetected
}

// buildCycle reconstructs the cycle closed by the back-edge from -> to, where
// `to` is an ancestor of `from` on the current DFS path. It walks `from` up its
// parent chain until it reaches `to` vertex, then closes the loop, returning
// the cycle as [to ... from to].
func buildCycle[T comparable](from, to *Vertex[T]) []*Vertex[T] {
	// Collect the path from `from` up to `to` (inclusive), which comes out
	// in reverse (leaf-to-root) order.
	path := make([]*Vertex[T], 0)
	for v := from; v != nil; v = v.Parent {
		path = append(path, v)
		if v == to {
			break
		}
	}

	// Reverse so the cycle reads root-first: [to ... from]
	slices.Reverse(path)

	// Close the loop back to `to`
	return append(path, to)
}

// WalkTopoOrder performs a topological sort and walks over the vertices in
// topological order.
//
// In case a cycle exists in the graph, WalkTopoOrder returns a *CycleError,
// which wraps ErrCycleDetected. The CycleError.Cycle field contains the ordered
// set of vertices forming the cycle. Use [errors.As] to extract it, or
// errors.Is(err, ErrCycleDetected) to test for a cycle without inspecting the
// path.
func WalkTopoOrder[T comparable](g Graph[T], walkFunc WalkFunc[T]) error {
	if g.Kind() != KindDirected {
		return ErrIsNotDirectedGraph
	}

	// Make sure to reset all vertex attributes
	g.ResetVertexAttributes()

	// A helper function, which performs post-order Depth-first Search (DFS)
	// traversal of the graph, starting from the given source vertex.
	//
	// If a cycle is found, then this function will return a *CycleError.
	//
	// This function is almost identical to WalkPostOrderDFS, except for the
	// fact that we don't reset the vertex attributes while performing DFS
	// on each vertex, and also we return a *CycleError whenever we detect a
	// cycle in the graph.
	//
	// The three vertex colors carry a strict meaning here, which is what
	// makes cycle detection correct:
	//
	//   - White: the vertex has not been discovered yet.
	//   - Gray:  the vertex has been entered and is currently on the
	//            active DFS path, i.e. it is an ancestor of the
	//            vertex being scanned.
	//   - Black: the vertex and all of its descendants are fully
	//            explored.
	//
	// A vertex is painted Gray only when its own frame becomes active (the
	// first time it is peeked), NOT when it is first discovered as a
	// neighbour. This ensures that a Gray neighbour is always a genuine
	// ancestor on the current path, so a back edge (and hence a cycle) is
	// reported only when one truly exists. Painting Gray at discovery time
	// would misreport a shared successor (e.g. the L in P->L, P->M, M->L)
	// as a cycle.
	dfsPostOrder := func(source *Vertex[T]) ([]*Vertex[T], error) {
		result := make([]*Vertex[T], 0)

		// Vertex has already been fully explored
		if source.Color == Black {
			return result, nil
		}

		stack := deque.New[*Vertex[T]]()
		stack.PushFront(source)

		for !stack.IsEmpty() {
			// Don't pop a vertex from the stack yet, just peek to
			// see if this vertex is ready
			v, err := stack.PeekFront()
			if err != nil {
				panic(err)
			}

			// The vertex may have been finished via another path
			// after it was pushed, or pushed more than once by
			// different parents. If it is already Black, discard
			// this stale stack entry.
			if v.Color == Black {
				if _, err := stack.PopFront(); err != nil {
					panic(err)
				}
				continue
			}

			// First time we peek this vertex: its frame becomes
			// active, so paint it Gray to mark it as being on the
			// current DFS path.
			v.Color = Gray

			isReady := true
			neighbours := g.GetNeighbourVertices(v.Value)
			for _, u := range neighbours {
				switch u.Color {
				case White:
					// First time seeing this neighbour. Only
					// discover it (push it); it is painted Gray
					// when its own frame becomes active.
					isReady = false
					u.DistanceFromSource = v.DistanceFromSource + 1
					u.Parent = v
					stack.PushFront(u)
				case Gray:
					// The neighbour is on the current DFS path,
					// so this is a back edge: a cycle has been
					// detected. Reconstruct the exact cycle from
					// the Parent chain.
					return result, &CycleError[T]{Cycle: buildCycle(v, u)}
				case Black:
					// Already fully explored, not part of a
					// cycle.
				}
			}

			if isReady {
				// The vertex is ready, pop it out
				popped, err := stack.PopFront()
				if err != nil {
					panic(err)
				}

				// We are done with vertex V
				popped.Color = Black
				result = append(result, popped)
			}
		}

		return result, nil
	}

	// Enqueue all vertices and perform post-order DFS on each
	queue := deque.New[*Vertex[T]]()
	for _, v := range g.GetVertices() {
		queue.PushBack(v)
	}

	for !queue.IsEmpty() {
		v, err := queue.PopFront()
		if err != nil {
			panic(err)
		}

		ready, err := dfsPostOrder(v)
		if err != nil {
			return err
		}

		for _, u := range ready {
			err := walkFunc(u)
			if err == ErrStopWalking {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// FindCycle returns the first cycle found in a directed graph, or nil if the
// graph is acyclic. The returned slice starts and ends with the same vertex,
// e.g. for the edges 1->2->3->1 it is [1 2 3 1].
//
// FindCycle returns [ErrIsNotDirectedGraph] if the graph is not directed.
//
// If the graph contains more than one cycle, only the first one encountered is
// returned. Which cycle that is depends on the internal iteration order and
// should not be relied upon.
func FindCycle[T comparable](g Graph[T]) ([]*Vertex[T], error) {
	noop := func(v *Vertex[T]) error {
		return nil
	}

	err := WalkTopoOrder(g, noop)
	if err == nil {
		return nil, nil
	}

	var cycleErr *CycleError[T]
	if errors.As(err, &cycleErr) {
		return cycleErr.Cycle, nil
	}

	// Some other error, e.g. ErrIsNotDirectedGraph
	return nil, err
}
