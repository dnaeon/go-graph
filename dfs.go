// Copyright (c) 2023 Marin Atanasov Nikolov <dnaeon@gmail.com>
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

import (
	"fmt"

	"gopkg.in/dnaeon/go-deque.v1"
)

// WalkPreOrderDFS performs pre-order Depth-first Search (DFS)
// traversal of the graph, starting from the given source vertex.
func WalkPreOrderDFS[T comparable](g Graph[T], source T, walkFunc WalkFunc[T]) error {
	if !g.VertexExists(source) {
		return fmt.Errorf("Source vertex %v not found in the graph", source)
	}

	// Make sure to reset all vertex attributes
	g.ResetVertexAttributes()

	// Push the source vertex to the stack and paint it
	srcVertex := g.GetVertex(source)
	srcVertex.Color = Gray
	stack := deque.New[*Vertex[T]]()
	stack.PushFront(srcVertex)

	for !stack.IsEmpty() {
		// Pop an item from the stack
		v, err := stack.PopFront()
		if err != nil {
			panic(err)
		}

		// Visit the neighbours of V
		neighbours := g.GetNeighbourVertices(v.Value)
		for _, u := range neighbours {
			// First time seeing this neighbour vertex,
			// push it to the stack
			if u.Color == White {
				u.Color = Gray
				u.DistanceFromSource = v.DistanceFromSource + 1
				u.Parent = v
				stack.PushFront(u)
			}
		}

		walkErr := walkFunc(v)
		if walkErr == ErrStopWalking {
			return nil
		}
		if walkErr != nil {
			return walkErr
		}

		// We are done with vertex V
		v.Color = Black
	}

	return nil
}

// WalkPostOrderDFS performs post-order Depth-first Search (DFS)
// traversal of the graph, starting from the given source vertex.
func WalkPostOrderDFS[T comparable](g Graph[T], source T, walkFunc WalkFunc[T]) error {
	if !g.VertexExists(source) {
		return fmt.Errorf("Source vertex %v not found in the graph", source)
	}

	// Make sure to reset all vertex attributes
	g.ResetVertexAttributes()

	// Push the source vertex to the stack and paint it
	srcVertex := g.GetVertex(source)
	srcVertex.Color = Gray
	stack := deque.New[*Vertex[T]]()
	stack.PushFront(srcVertex)

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
			// First time seeing this neighbour
			if u.Color == White {
				isReady = false
				u.Color = Gray
				u.DistanceFromSource = v.DistanceFromSource + 1
				u.Parent = v
				stack.PushFront(u)
			}
		}

		if isReady {
			// The vertex is ready, pop it out
			popped, err := stack.PopFront()
			if err != nil {
				panic(err)
			}

			walkErr := walkFunc(v)
			if walkErr == ErrStopWalking {
				return nil
			}
			if walkErr != nil {
				return walkErr
			}
			popped.Color = Black
		}
	}

	return nil
}

// WalkUnreachableVertices walks over the vertices which are
// unreachable from the given source vertex
func WalkUnreachableVertices[T comparable](g Graph[T], source T, walkFunc WalkFunc[T]) error {
	// In order to find all unreachable vertices we will first DFS
	// traverse the graph.  The vertices which remain White after
	// we've walked the graph are unreachable from the source
	// vertex.
	dummyDfsWalk := func(v *Vertex[T]) error {
		return nil
	}

	if err := WalkPreOrderDFS(g, source, dummyDfsWalk); err != nil {
		return err
	}

	for _, v := range g.GetVertices() {
		if v.Color == White {
			err := walkFunc(v)
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
