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

package graph_test

import (
	"errors"
	"slices"
	"testing"

	"gopkg.in/dnaeon/go-graph.v1"
)

func TestWalkShortestPath(t *testing.T) {
	g := newUndirectedWeightedGraph()
	if g.Kind() != graph.KindUndirected {
		t.Fatal("graph is not undirected")
	}

	// Vertices which form the shortest path from (1) to (8)
	wantShortestPath := []*graph.Vertex[int]{
		{
			Value:              1,
			DistanceFromSource: 0,
			Color:              graph.White,
		},
		{
			Value:              2,
			DistanceFromSource: 2,
			Color:              graph.White,
		},
		{
			Value:              4,
			DistanceFromSource: 5,
			Color:              graph.White,
		},
		{
			Value:              5,
			DistanceFromSource: 14,
			Color:              graph.White,
		},
		{
			Value:              7,
			DistanceFromSource: 18,
			Color:              graph.White,
		},
		{
			Value:              8,
			DistanceFromSource: 26,
			Color:              graph.White,
		},
	}

	collector := g.NewCollector()
	err := graph.WalkShortestPath(g, 1, 8, collector.WalkFunc)
	if err != nil {
		t.Fatal(err)
	}

	gotShortestPath := collector.Get()
	verifyVertices(t, wantShortestPath, gotShortestPath)

	// No path should exist between (1) and (10)
	err = graph.WalkShortestPath(g, 1, 10, collector.WalkFunc)
	if err == nil {
		t.Fatalf("no path should exist between 1 and 10")
	}

	// Source vertex does not exist
	dummyWalker := func(v *graph.Vertex[int]) error {
		return nil
	}
	if err := graph.WalkShortestPath(g, 42, 1, dummyWalker); err == nil {
		t.Fatal("WalkShortestPath should fail with non-existing source vertex")
	}

	// Destination vertex does not exist
	if err := graph.WalkShortestPath(g, 1, 42, dummyWalker); err == nil {
		t.Fatal("WalkShortestPath should fail with non-existing destination vertex")
	}

	// Short-circuit the walking
	result := make([]int, 0)
	shortCircuitWalker := func(v *graph.Vertex[int]) error {
		if v.Value == 4 {
			return graph.ErrStopWalking
		}
		result = append(result, v.Value)
		return nil
	}
	if err := graph.WalkShortestPath(g, 1, 8, shortCircuitWalker); err != nil {
		t.Fatal(err)
	}
	// Collected path should include []int{1, 2} only
	if !slices.Equal([]int{1, 2}, result) {
		t.Fatal("WalkShortestPath: short circuit walker values mismatch")
	}

	// Signal custom error while walking
	myErr := errors.New("my custom error")
	errWalker := func(v *graph.Vertex[int]) error {
		return myErr
	}
	if err := graph.WalkShortestPath(g, 1, 8, errWalker); err != myErr {
		t.Fatal("WalkShortestPath should fail with custom error")
	}

	// Invoke WalkDijkstra with non-existing vertex
	if err := graph.WalkDijkstra(g, 42, dummyWalker); err == nil {
		t.Fatal("WalkDijkstra should fail with non-existing vertex")
	}
}
