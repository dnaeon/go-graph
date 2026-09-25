// Copyright (c) 2026 Marin Atanasov Nikolov <dnaeon@gmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//   1. Redistributions of source code must retain the above copyright
//      notice, this list of conditions and the following disclaimer.
//   2. Redistributions in binary form must reproduce the above copyright
//      notice, this list of conditions and the following disclaimer in the
//      documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package graph

import "gopkg.in/dnaeon/go-deque.v1"

// tarjanFrame is a single work item on the explicit DFS stack used by the
// iterative Tarjan traversal in StronglyConnectedComponents. It tracks the
// vertex being visited and the position within its neighbour list that has been
// processed so far, so the traversal can be suspended and resumed without
// recursion.
type tarjanFrame[T comparable] struct {
	vertex     *Vertex[T]
	neighbours []*Vertex[T]
	next       int // index into neighbours to process next
}

// StronglyConnectedComponents computes the strongly connected components (SCCs)
// of a directed graph using Tarjan's algorithm.
//
// A strongly connected component is a maximal set of vertices such that every
// vertex is reachable from every other vertex in the set, following edge
// directions. In a directed acyclic graph every SCC is a single vertex; a
// component with more than one vertex (or a single vertex with a self-loop) is a
// cyclic cluster.
//
// StronglyConnectedComponents returns [ErrIsNotDirectedGraph] if the graph is
// not directed.
//
// The order of the returned components, and the order of the vertices within
// each component, is unspecified and should not be relied upon, as it depends on
// the internal iteration order of the graph.
//
// The algorithm is implemented iteratively, so deep graphs do not exhaust the
// goroutine stack.
func StronglyConnectedComponents[T comparable](g Graph[T]) ([][]*Vertex[T], error) {
	if g.Kind() != KindDirected {
		return nil, ErrIsNotDirectedGraph
	}

	// Per-vertex bookkeeping used by Tarjan's algorithm. These are kept in
	// local maps rather than on the Vertex struct, since they are private to
	// this algorithm.
	//
	//   - index:   the order in which a vertex was discovered (its DFS index)
	//   - lowLink: the smallest index reachable from the vertex (including
	//              itself) via the current DFS subtree and at most one back edge
	//   - onStack: whether the vertex is currently on the Tarjan output stack
	index := make(map[*Vertex[T]]int)
	lowLink := make(map[*Vertex[T]]int)
	onStack := make(map[*Vertex[T]]bool)

	// The Tarjan output stack, holding the vertices of the SCC currently being
	// assembled.
	stack := make([]*Vertex[T], 0)

	// A monotonically increasing counter, handing out DFS indices.
	nextIndex := 0

	result := make([][]*Vertex[T], 0)

	// strongConnect performs the iterative Tarjan DFS starting from the given
	// source vertex.
	strongConnect := func(source *Vertex[T]) {
		enter := func(v *Vertex[T]) tarjanFrame[T] {
			index[v] = nextIndex
			lowLink[v] = nextIndex
			nextIndex++
			stack = append(stack, v)
			onStack[v] = true
			return tarjanFrame[T]{
				vertex:     v,
				neighbours: g.GetNeighbourVertices(v.Value),
				next:       0,
			}
		}

		dfsStack := []tarjanFrame[T]{enter(source)}

		for len(dfsStack) > 0 {
			top := &dfsStack[len(dfsStack)-1]
			v := top.vertex

			if top.next < len(top.neighbours) {
				// Process the next neighbour of v.
				u := top.neighbours[top.next]
				top.next++

				if _, seen := index[u]; !seen {
					// u has not been visited yet: descend into it. Its
					// low-link will be folded back into v's when we pop
					// u's frame (handled below).
					dfsStack = append(dfsStack, enter(u))
				} else if onStack[u] {
					// u is on the stack, hence in the current SCC being
					// assembled: this is a back/cross edge within the
					// component. Update v's low-link with u's index.
					if index[u] < lowLink[v] {
						lowLink[v] = index[u]
					}
				}
				// If u is visited but not on the stack, it already belongs
				// to a finished SCC; ignore it.
				continue
			}

			// All neighbours of v processed. v's frame is done.
			dfsStack = dfsStack[:len(dfsStack)-1]

			// Fold v's low-link into its parent's low-link, mirroring the
			// return step of the recursive formulation.
			if len(dfsStack) > 0 {
				parent := dfsStack[len(dfsStack)-1].vertex
				if lowLink[v] < lowLink[parent] {
					lowLink[parent] = lowLink[v]
				}
			}

			// If v is the root of an SCC (its low-link equals its own
			// index), pop the whole component off the stack.
			if lowLink[v] == index[v] {
				component := make([]*Vertex[T], 0)
				for {
					w := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					onStack[w] = false
					component = append(component, w)
					if w == v {
						break
					}
				}
				result = append(result, component)
			}
		}
	}

	for _, v := range g.GetVertices() {
		if _, seen := index[v]; !seen {
			strongConnect(v)
		}
	}

	return result, nil
}

// CycleWalkFunc is a function which receives a single cycle while walking over
// the cycles in a graph. The cycle is given as an ordered, closed path, e.g. for
// the mutual dependency 1 <-> 2 the cycle is [1 2 1].
type CycleWalkFunc[T comparable] func(cycle []*Vertex[T]) error

// WalkCycles walks over the cycles in a directed graph, invoking walkFunc once
// per cyclic cluster with a representative cycle for that cluster.
//
// WalkCycles first computes the strongly connected components (see
// [StronglyConnectedComponents]) and then, for each non-trivial component (one
// with more than one vertex, or a single vertex with a self-loop), extracts a
// single representative cycle confined to that component and passes it to
// walkFunc. Trivial components (a single vertex with no self-loop) are not
// cycles and are skipped.
//
// Note that WalkCycles yields one cycle per cyclic cluster, not every simple
// cycle within it. A single strongly connected component may contain many simple
// cycles; enumerating all of them is exponential in the worst case and is
// intentionally not attempted here.
//
// A walkFunc may return [ErrStopWalking] to stop the traversal early, in which
// case WalkCycles returns nil. Any other error returned by walkFunc is
// propagated to the caller.
//
// WalkCycles returns [ErrIsNotDirectedGraph] if the graph is not directed.
//
// The order in which cycles are visited is unspecified and should not be relied
// upon.
func WalkCycles[T comparable](g Graph[T], walkFunc CycleWalkFunc[T]) error {
	sccs, err := StronglyConnectedComponents(g)
	if err != nil {
		return err
	}

	for _, scc := range sccs {
		var cycle []*Vertex[T]

		if len(scc) == 1 {
			// A single-vertex component is only a cycle if the vertex has
			// a self-loop.
			v := scc[0]
			if !g.EdgeExists(v.Value, v.Value) {
				continue
			}
			cycle = []*Vertex[T]{v, v}
		} else {
			// A component with more than one vertex is strongly connected,
			// so a cycle confined to it is guaranteed to exist.
			cycle = cycleWithin(g, scc)
		}

		walkErr := walkFunc(cycle)
		if walkErr == ErrStopWalking {
			return nil
		}
		if walkErr != nil {
			return walkErr
		}
	}

	return nil
}

// cycleWithin extracts a single cycle confined to the vertices of the given
// strongly connected component, which must contain more than one vertex.
//
// It performs the same three-color, Gray-on-entry depth-first search used by
// WalkTopoOrder, but only follows edges whose target is a member of the
// component. Because the subgraph induced by a non-trivial SCC is strongly
// connected, this DFS is guaranteed to encounter a back edge to a Gray
// (on-path) vertex, at which point the cycle is reconstructed from the parent
// chain via buildCycle.
func cycleWithin[T comparable](g Graph[T], scc []*Vertex[T]) []*Vertex[T] {
	// The set of vertex values that make up this component, used to confine
	// the traversal to the component's induced subgraph.
	members := make(map[T]bool, len(scc))
	for _, v := range scc {
		members[v.Value] = true
	}

	// Reset the traversal attributes for the component's vertices so that the
	// coloring below starts from a clean state.
	for _, v := range scc {
		v.Color = White
		v.DistanceFromSource = 0.0
		v.Parent = nil
	}

	stack := deque.New[*Vertex[T]]()
	stack.PushFront(scc[0])

	for !stack.IsEmpty() {
		v, err := stack.PeekFront()
		if err != nil {
			panic(err)
		}

		if v.Color == Black {
			if _, err := stack.PopFront(); err != nil {
				panic(err)
			}
			continue
		}

		// First time this vertex's frame becomes active: paint it Gray to
		// mark it as being on the current DFS path.
		v.Color = Gray

		isReady := true
		for _, u := range g.GetNeighbourVertices(v.Value) {
			// Confine the traversal to the component.
			if !members[u.Value] {
				continue
			}

			switch u.Color {
			case White:
				isReady = false
				u.DistanceFromSource = v.DistanceFromSource + 1
				u.Parent = v
				stack.PushFront(u)
			case Gray:
				// Back edge to an ancestor on the current path: this
				// closes a cycle confined to the component.
				return buildCycle(v, u)
			case Black:
				// Already fully explored within this component.
			}
		}

		if isReady {
			popped, err := stack.PopFront()
			if err != nil {
				panic(err)
			}
			popped.Color = Black
		}
	}

	// Unreachable for a genuine, strongly connected component with more than
	// one vertex: a back edge is always found before the stack empties.
	return nil
}
