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
	"math"
	"slices"

	"gopkg.in/dnaeon/go-priorityqueue.v1"
)

// Initializes the source vertex as part of Dijkstra's algorithm
func initializeSourceVertex[T comparable](g Graph[T], source T) error {
	if !g.VertexExists(source) {
		return fmt.Errorf("Source vertex %v not found in graph", source)
	}

	// Set tentative distance for all vertices
	g.ResetVertexAttributes()
	for _, v := range g.GetVertices() {
		v.DistanceFromSource = math.Inf(1)
	}

	// Initialize source vertex
	srcV := g.GetVertex(source)
	srcV.DistanceFromSource = 0.0

	return nil
}

// Relaxes the edge as part of Dijkstra's algorithm
func relaxEdge[T comparable](g Graph[T], from T, to T) error {
	if !g.EdgeExists(from, to) {
		return fmt.Errorf("No edge exists between %v and %v", from, to)
	}

	edge := g.GetEdge(from, to)
	fromV := g.GetVertex(from)
	toV := g.GetVertex(to)

	// Compute alt distance and compare against current distance
	alt := fromV.DistanceFromSource + edge.Weight
	if alt < toV.DistanceFromSource {
		toV.DistanceFromSource = alt
		toV.Parent = fromV
	}

	return nil
}

// WalkDijkstra implements Dijkstra's algorithm for finding the
// shortest-path from a given source vertex to all other vertices in
// the graph.
//
// This method will pass each visited vertex while traversing the
// shortest path from the source vertex to each other vertex in the
// graph.
//
// Note, that this method only constructs the shortest-path tree and
// yields each visited vertex. In order to stop walking the graph
// callers of this method should return ErrStopWalking error and refer
// to the shortest-path tree, or use the WalkShortestPath method.
func WalkDijkstra[T comparable](g Graph[T], source T, walkFunc WalkFunc[T]) error {
	if err := initializeSourceVertex(g, source); err != nil {
		return err
	}

	// Enqueue all vertices
	queue := priorityqueue.New[*Vertex[T], float64](priorityqueue.MinHeap)
	for _, v := range g.GetVertices() {
		queue.Put(v, v.DistanceFromSource)
	}

	for !queue.IsEmpty() {
		item := queue.Get()
		v := item.Value
		// Relax edges connecting V and it's neighbours
		for _, u := range g.GetNeighbourVertices(v.Value) {
			oldDist := u.DistanceFromSource
			if err := relaxEdge(g, v.Value, u.Value); err != nil {
				return err
			}
			// Update the priority, if needed
			if u.DistanceFromSource != oldDist {
				queue.Update(u, u.DistanceFromSource)
			}
		}

		err := walkFunc(v)
		if err == ErrStopWalking {
			return nil
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// WalkShortestPath yields the vertices which represent the shortest
// path between SOURCE and DEST.
func WalkShortestPath[T comparable](g Graph[T], source T, dest T, walkFunc WalkFunc[T]) error {
	if !g.VertexExists(source) {
		return fmt.Errorf("Source vertex %v not found in the graph", source)
	}

	if !g.VertexExists(dest) {
		return fmt.Errorf("Destination vertex %v not found in the graph", source)
	}

	// A walker which stops walking the graph, as soon as we reach
	// the destination vertex
	walker := func(v *Vertex[T]) error {
		if v.Value == dest {
			return ErrStopWalking
		}
		return nil
	}

	if err := WalkDijkstra(g, source, walker); err != nil {
		return err
	}

	// Make our way from the destination vertex back to the source
	// by following the relationships established by the
	// shortest-path tree.
	destV := g.GetVertex(dest)
	result := make([]*Vertex[T], 0)
	v := destV
	for {
		result = append(result, v)
		if v.Value == source {
			break
		}

		if v.Parent == nil {
			return fmt.Errorf("No path exists between %v and %v", source, dest)
		}
		v = v.Parent
	}

	slices.Reverse(result)
	for _, v := range result {
		err := walkFunc(v)
		if err == ErrStopWalking {
			return nil
		}
		if err != nil {
			return err
		}
	}

	return nil
}
