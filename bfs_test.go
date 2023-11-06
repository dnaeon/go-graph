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

func TestWalkBFS(t *testing.T) {
	g := newUndirectedGraph()
	if g.Kind() != graph.KindUndirected {
		t.Fatal("graph is not undirected")
	}

	collector := g.NewCollector()

	// Vertices we want during a BFS walk, starting from (1)
	wantVerticesFromSrc1 := []*graph.Vertex[int]{
		{
			Value:              1,
			DistanceFromSource: 0.0,
			Color:              graph.Black,
		},
		{
			Value:              2,
			DistanceFromSource: 1.0,
			Color:              graph.Black,
		},
		{
			Value:              3,
			DistanceFromSource: 1.0,
			Color:              graph.Black,
		},
		{
			Value:              4,
			DistanceFromSource: 2.0,
			Color:              graph.Black,
		},
		{
			Value:              5,
			DistanceFromSource: 3.0,
			Color:              graph.Black,
		},
	}

	err := graph.WalkBFS(g, 1, collector.WalkFunc)
	if err != nil {
		t.Fatal(err)
	}

	gotVerticesFromSrc1 := collector.Get()
	verifyVertices(t, wantVerticesFromSrc1, gotVerticesFromSrc1)

	// Vertices we want during a BFS walk, starting from (3)
	wantVerticesFromSrc3 := []*graph.Vertex[int]{
		{
			Value:              3,
			DistanceFromSource: 0.0,
			Color:              graph.Black,
		},
		{
			Value:              1,
			DistanceFromSource: 1.0,
			Color:              graph.Black,
		},
		{
			Value:              4,
			DistanceFromSource: 1.0,
			Color:              graph.Black,
		},
		{
			Value:              2,
			DistanceFromSource: 2.0,
			Color:              graph.Black,
		},
		{
			Value:              5,
			DistanceFromSource: 2.0,
			Color:              graph.Black,
		},
	}

	collector.Reset()
	if err := graph.WalkBFS(g, 3, collector.WalkFunc); err != nil {
		t.Fatal(err)
	}
	gotVerticesFromSrc3 := collector.Get()
	verifyVertices(t, wantVerticesFromSrc3, gotVerticesFromSrc3)

	// No such vertex exists, error is expected
	collector.Reset()
	if err := graph.WalkBFS(g, 42, collector.WalkFunc); err == nil {
		t.Fatal("expected an error during WalkBFS with non-existing vertex")
	}

	// Stop walking the graph when we reach a given vertex
	shortCircuitVertexValues := make([]int, 0)
	shortCircuitWalker := func(v *graph.Vertex[int]) error {
		if v.Value == 4 {
			return graph.ErrStopWalking
		}
		shortCircuitVertexValues = append(shortCircuitVertexValues, v.Value)
		return nil
	}
	if err := graph.WalkBFS(g, 1, shortCircuitWalker); err != nil {
		t.Fatal("WalkBFS failed with a short-ciruit walker")
	}
	if !slices.Equal(shortCircuitVertexValues, []int{1, 2, 3}) {
		t.Fatal("WalkBFS with short-circuit walker yielded mismatched values")
	}

	// DFS with a walker which signals an error
	myErr := errors.New("My custom error")
	walkWithError := func(v *graph.Vertex[int]) error {
		return myErr
	}
	if err := graph.WalkBFS(g, 1, walkWithError); err != myErr {
		t.Fatal("WalkBFS is expected to return our custom error")
	}
}
