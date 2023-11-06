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

func TestWalkPreOrderDFS(t *testing.T) {
	g := newUndirectedGraph()
	if g.Kind() != graph.KindUndirected {
		t.Fatal("graph is not undirected")
	}

	collector := g.NewCollector()

	// Vertices we want during a DFS walk, starting from (1)
	wantVerticesFromSrc1 := []*graph.Vertex[int]{
		{
			Value:              1,
			DistanceFromSource: 0.0,
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
		{
			Value:              2,
			DistanceFromSource: 1.0,
			Color:              graph.Black,
		},
	}

	err := graph.WalkPreOrderDFS(g, 1, collector.WalkFunc)
	if err != nil {
		t.Fatal(err)
	}

	gotVerticesFromSrc1 := collector.Get()
	verifyVertices(t, wantVerticesFromSrc1, gotVerticesFromSrc1)

	// Vertices we want during a DFS walk, starting from (3)
	wantVerticesFromSrc3 := []*graph.Vertex[int]{
		{
			Value:              3,
			DistanceFromSource: 0.0,
			Color:              graph.Black,
		},
		{
			Value:              4,
			DistanceFromSource: 1.0,
			Color:              graph.Black,
		},
		{
			Value:              5,
			DistanceFromSource: 2.0,
			Color:              graph.Black,
		},
		{
			Value:              1,
			DistanceFromSource: 1.0,
			Color:              graph.Black,
		},
		{
			Value:              2,
			DistanceFromSource: 2.0,
			Color:              graph.Black,
		},
	}

	collector.Reset()
	if err := graph.WalkPreOrderDFS(g, 3, collector.WalkFunc); err != nil {
		t.Fatal(err)
	}
	gotVerticesFromSrc3 := collector.Get()
	verifyVertices(t, wantVerticesFromSrc3, gotVerticesFromSrc3)

	// No such vertex exists, error is expected
	collector.Reset()
	if err := graph.WalkPreOrderDFS(g, 42, collector.WalkFunc); err == nil {
		t.Fatal("expected an error during WalkPreOrderDFS with non-existing vertex")
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
	if err := graph.WalkPreOrderDFS(g, 1, shortCircuitWalker); err != nil {
		t.Fatal("WalkPreOrderDFS failed with a short-ciruit walker")
	}
	if !slices.Equal(shortCircuitVertexValues, []int{1, 3}) {
		t.Fatal("WalkPreOrderDFS with short-circuit walker yielded mismatched values")
	}

	// DFS with a walker which signals an error
	myErr := errors.New("My custom error")
	walkWithError := func(v *graph.Vertex[int]) error {
		return myErr
	}
	if err := graph.WalkPreOrderDFS(g, 1, walkWithError); err != myErr {
		t.Fatal("WalkPreOrderDFS is expected to return our custom error")
	}
}

func TestWalkPostOrderDFS(t *testing.T) {
	g := newUndirectedGraph()
	if g.Kind() != graph.KindUndirected {
		t.Fatal("graph is not undirected")
	}

	collector := g.NewCollector()

	// Vertices we want during a post-order DFS walk, starting
	// from (1)
	wantVerticesFromSrc1 := []*graph.Vertex[int]{
		{
			Value:              5,
			DistanceFromSource: 3.0,
			Color:              graph.Black,
		},
		{
			Value:              4,
			DistanceFromSource: 2.0,
			Color:              graph.Black,
		},
		{
			Value:              3,
			DistanceFromSource: 1.0,
			Color:              graph.Black,
		},
		{
			Value:              2,
			DistanceFromSource: 1.0,
			Color:              graph.Black,
		},
		{
			Value:              1,
			DistanceFromSource: 0.0,
			Color:              graph.Black,
		},
	}

	err := graph.WalkPostOrderDFS(g, 1, collector.WalkFunc)
	if err != nil {
		t.Fatal(err)
	}

	gotVerticesFromSrc1 := collector.Get()
	verifyVertices(t, wantVerticesFromSrc1, gotVerticesFromSrc1)

	// Vertices we want during a DFS walk, starting from (3)
	wantVerticesFromSrc3 := []*graph.Vertex[int]{
		{
			Value:              5,
			DistanceFromSource: 2.0,
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
			Value:              1,
			DistanceFromSource: 1.0,
			Color:              graph.Black,
		},
		{
			Value:              3,
			DistanceFromSource: 0.0,
			Color:              graph.Black,
		},
	}

	collector.Reset()
	if err := graph.WalkPostOrderDFS(g, 3, collector.WalkFunc); err != nil {
		t.Fatal(err)
	}
	gotVerticesFromSrc3 := collector.Get()
	verifyVertices(t, wantVerticesFromSrc3, gotVerticesFromSrc3)

	// No such vertex exists, error is expected
	collector.Reset()
	if err := graph.WalkPostOrderDFS(g, 42, collector.WalkFunc); err == nil {
		t.Fatal("expected an error during WalkPostOrderDFS with non-existing vertex")
	}

	// Stop walking the graph when we reach a given vertex
	shortCircuitVertexValues := make([]int, 0)
	shortCircuitWalker := func(v *graph.Vertex[int]) error {
		if v.Value == 3 {
			return graph.ErrStopWalking
		}
		shortCircuitVertexValues = append(shortCircuitVertexValues, v.Value)
		return nil
	}
	if err := graph.WalkPostOrderDFS(g, 1, shortCircuitWalker); err != nil {
		t.Fatal("WalkPostOrderDFS failed with a short-ciruit walker")
	}
	if !slices.Equal(shortCircuitVertexValues, []int{5, 4}) {
		t.Fatal("WalkPostOrderDFS with short-circuit walker yielded mismatched values")
	}

	// DFS with a walker which signals an error
	myErr := errors.New("My custom error")
	walkWithError := func(v *graph.Vertex[int]) error {
		return myErr
	}
	if err := graph.WalkPostOrderDFS(g, 1, walkWithError); err != myErr {
		t.Fatal("WalkPostOrderDFS is expected to return our custom error")
	}
}

func TestWalkUnreachableVertices(t *testing.T) {
	g := newUndirectedGraph()

	// Collect the unreachable vertex values
	got := make([]int, 0)
	walker := func(v *graph.Vertex[int]) error {
		got = append(got, v.Value)
		return nil
	}
	if err := graph.WalkUnreachableVertices(g, 1, walker); err != nil {
		t.Fatal(err)
	}

	// Vertices unreachable from (1)
	want := []int{10, 11, 12, 13}
	slices.Sort(got)

	if len(want) != len(got) {
		t.Fatalf("got %d items, want %d items", len(want), len(got))
	}

	for i, val := range want {
		if got[i] != val {
			t.Fatalf("got %v vertex, want %v", got[i], val)
		}
	}

	// Expecting an error when starting with a non-existing vertex
	dummyWalker := func(v *graph.Vertex[int]) error {
		return nil
	}
	if err := graph.WalkUnreachableVertices(g, 42, dummyWalker); err == nil {
		t.Fatal("expected and error during WalkUnreachableVertices with non-existing vertex")
	}

	// Stop walking when we reach a given vertex
	shortCircuitVertexValues := make([]int, 0)
	shortCircuitWalker := func(v *graph.Vertex[int]) error {
		if v.Value == 13 {
			return graph.ErrStopWalking
		}
		shortCircuitVertexValues = append(shortCircuitVertexValues, v.Value)
		return nil
	}
	if err := graph.WalkUnreachableVertices(g, 1, shortCircuitWalker); err != nil {
		t.Fatal("WalkUnreachableVertices failed with a short-ciruit walker")
	}

	// Since the order in which unreachable vertices is
	// non-deterministic, we should expect that the number of
	// unreachable vertices is 1 less than the total number of
	// unreachable vertices
	if len(shortCircuitVertexValues) > 3 {
		t.Fatal("WalkUnreachableVertices yielded too many vertices")
	}

	// Walk with a custom error being signalled
	myErr := errors.New("My custom error")
	walkWithError := func(v *graph.Vertex[int]) error {
		return myErr
	}
	if err := graph.WalkUnreachableVertices(g, 1, walkWithError); err != myErr {
		t.Fatal("WalkUnreachableVertices is expected to return our custom error")
	}
}
