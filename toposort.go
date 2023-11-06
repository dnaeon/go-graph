package graph

import (
	"errors"

	"gopkg.in/dnaeon/go-deque.v1"
)

// ErrCycleDetected is returned whenever a cycle has been detected in
// the graph.
var ErrCycleDetected = errors.New("cycle detected")

// ErrIsNotDirectedGraph is returned whenever an operation cannot be
// performed, because the graph is not directed.
var ErrIsNotDirectedGraph = errors.New("graph is not directed")

// WalkTopoOrder performs a topological sort and walks over the
// vertices in topological order.
//
// In case a cycle exists in the graph, WalkTopoOrder will return
// ErrCycleDetected.
//
// In case ErrCycleDetected is returned, the vertices which remained
// Gray are forming a cyclic path in the graph.
func WalkTopoOrder[T comparable](g Graph[T], walkFunc WalkFunc[T]) error {
	if g.Kind() != KindDirected {
		return ErrIsNotDirectedGraph
	}

	// Make sure to reset all vertex attributes
	g.ResetVertexAttributes()

	// A helper function, which performs post-order Depth-first
	// Search (DFS) traversal of the graph, starting from the
	// given source vertex.
	//
	// If a cycle is found, then this function will return
	// ErrCycleDetected.
	//
	// This function almost identical to WalkPostOrderDFS, except
	// for the fact that we don't reset the vertex attributes
	// while performing DFS on each vertex, and also we return
	// ErrCycleDetected whenever we detect a cycle in the graph.
	dfsPostOrder := func(source *Vertex[T]) ([]*Vertex[T], error) {
		result := make([]*Vertex[T], 0)

		// Vertex has already been visited
		if source.Color == Black {
			return result, nil
		}

		// Push source vertex to the stack and paint it
		source.Color = Gray
		stack := deque.New[*Vertex[T]]()
		stack.PushFront(source)

		for !stack.IsEmpty() {
			// Don't pop a vertex from the stack yet, just peek to
			// see if this vertex is ready
			v, err := stack.PeekFront()
			if err != nil {
				panic(err)
			}

			isReady := true
			neighbours := g.GetNeighbourVertices(v.Value)
			for _, u := range neighbours {
				if u.Color == White {
					// First time seeing this neighbour
					isReady = false
					u.Color = Gray
					u.DistanceFromSource = v.DistanceFromSource + 1
					u.Parent = v
					stack.PushFront(u)
				} else if u.Color == Gray {
					// Seen this neighbour before, cycle
					// has been detected
					return result, ErrCycleDetected
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
