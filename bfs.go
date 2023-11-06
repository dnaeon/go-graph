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

// WalkBFS performs Breadth-first Search (BFS) traversal of the graph,
// starting from the given source vertex.
func WalkBFS[T comparable](g Graph[T], source T, walkFunc WalkFunc[T]) error {
	if !g.VertexExists(source) {
		return fmt.Errorf("Source vertex %v not found in the graph", source)
	}

	// Make sure to reset all vertex attributes
	g.ResetVertexAttributes()

	// Push the source vertex to the queue and paint it
	srcVertex := g.GetVertex(source)
	srcVertex.Color = Gray
	queue := deque.New[*Vertex[T]]()
	queue.PushBack(srcVertex)

	for !queue.IsEmpty() {
		// Pop an item from the queue
		v, err := queue.PopFront()
		if err != nil {
			panic(err)
		}

		// Visit neighbours of V
		neighbours := g.GetNeighbourVertices(v.Value)
		for _, u := range neighbours {
			// First time seeing this vertex
			if u.Color == White {
				u.Color = Gray
				u.DistanceFromSource = v.DistanceFromSource + 1
				u.Parent = v
				queue.PushBack(u)
			}
		}

		walkErr := walkFunc(v)
		if walkErr == ErrStopWalking {
			return nil
		}
		if walkErr != nil {
			return walkErr
		}

		// We are done with V
		v.Color = Black
	}

	return nil
}
