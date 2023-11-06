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
	"testing"

	"gopkg.in/dnaeon/go-graph.v1"
)

// Creates a new test undirected graph
func newUndirectedGraph() graph.Graph[int] {
	g := graph.New[int](graph.KindUndirected)
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)
	g.AddEdge(3, 4)
	g.AddEdge(4, 5)

	// The following vertices form a cluster, which is unreachable
	// from the vertices above
	g.AddEdge(10, 11)
	g.AddEdge(11, 12)
	g.AddEdge(11, 13)

	return g
}

// Creates a new test undirected weighted graph
func newUndirectedWeightedGraph() graph.Graph[int] {
	g := graph.New[int](graph.KindUndirected)
	g.AddWeightedEdge(1, 2, 2)
	g.AddWeightedEdge(1, 3, 6)
	g.AddWeightedEdge(2, 3, 7)
	g.AddWeightedEdge(2, 4, 3)
	g.AddWeightedEdge(3, 4, 4)
	g.AddWeightedEdge(4, 5, 9)
	g.AddWeightedEdge(5, 6, 11)
	g.AddWeightedEdge(5, 7, 4)
	g.AddWeightedEdge(6, 7, 6)
	g.AddWeightedEdge(6, 8, 5)
	g.AddWeightedEdge(7, 8, 8)

	// Isolated vertices
	g.AddWeightedEdge(10, 11, 1)

	return g
}

// Creates a new test directed graph
func newDirectedGraph() graph.Graph[int] {
	g := graph.New[int](graph.KindDirected)
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)
	g.AddEdge(3, 4)
	g.AddEdge(4, 5)

	// The following vertices form a cluster, unreachable from the
	// previous vertices
	g.AddEdge(10, 11)
	g.AddEdge(11, 12)
	g.AddEdge(11, 13)

	return g
}

// A helper function to compare vertices we got after walking the
// graph against an expected set of vertices
func verifyVertices[T comparable](t *testing.T, want []*graph.Vertex[T], got []*graph.Vertex[T]) {
	if len(want) != len(got) {
		t.Fatalf("got %d number of vertices, want %d", len(got), len(want))
	}

	for i, v := range got {
		want := want[i]
		if want.Value != v.Value {
			t.Fatalf("want vertex value %v, got %v", want.Value, v.Value)
		}
		if want.DistanceFromSource != v.DistanceFromSource {
			t.Fatalf("want distance from source %.2f for vertex %v, got %.2f", want.DistanceFromSource, v.Value, v.DistanceFromSource)
		}
		if want.Color != v.Color {
			t.Fatalf("want %v color for vertex %v, got %v", want.Color, v.Value, v.Color)
		}
	}
}

func TestCreateUndirectedGraph(t *testing.T) {
	g := graph.New[int](graph.KindUndirected)
	if g.Kind() != graph.KindUndirected {
		t.Fatal("graph is expected to be undirected")
	}

	v := g.AddVertex(1)
	if v.Value != 1 {
		t.Fatalf("newly added vertex value must be 1")
	}

	e1 := g.AddEdge(1, 2)
	if e1 == nil {
		t.Fatal("got nil edge, expected valid edge")
	}

	if e1.From != 1 {
		t.Fatalf("outgoing edge vertex must be 1")
	}

	if e1.To != 2 {
		t.Fatalf("incoming edge vertex must be 2")
	}

	if e1.Weight != 0.0 {
		t.Fatal("edge weight must be 0")
	}

	if !g.VertexExists(1) {
		t.Fatal("vertex 1 must exist")
	}

	// The number of edges we have so far is 1
	if len(g.GetEdges()) != 1 {
		t.Fatal("expected number of edges in graph is 1")
	}

	// Adding same edge should not create a new edge
	e1Prime := g.AddEdge(1, 2)
	if e1.From != e1Prime.From {
		t.Fatal("e1.From and e1Prime.From are not equal")
	}
	if e1.To != e1Prime.To {
		t.Fatal("e1.To and e1Prime.To are not equal")
	}

	// The number of edges we have is still 1
	if len(g.GetEdges()) != 1 {
		t.Fatal("expected number of edges in graph is 1")
	}

	// Non-existing vertex
	if g.VertexExists(42) {
		t.Fatal("non-existing vertex found")
	}

	if !g.EdgeExists(1, 2) {
		t.Fatal("edge (1,2) does not exist")
	}

	// Non-existing edge
	if g.EdgeExists(1, 42) {
		t.Fatal("non-existing edge (1, 42) found")
	}

	e2 := g.GetEdge(1, 42)
	if e2 != nil {
		t.Fatal("non-existing edge (1, 42) retrieved")
	}

	e3 := g.AddEdge(1, 100)
	e3.DotAttributes["color"] = "red"

	wantNumVertices := 3
	wantNumEdges := 2

	if wantNumVertices != len(g.GetVertices()) {
		t.Fatalf("want %d number of vertices, got %d", wantNumVertices, len(g.GetVertices()))
	}

	if wantNumVertices != len(g.GetVertexValues()) {
		t.Fatalf("want %d number of vertices, got %d", wantNumVertices, len(g.GetVertexValues()))
	}

	if wantNumEdges != len(g.GetEdges()) {
		t.Fatalf("want %d number of edges, got %d", wantNumEdges, len(g.GetEdges()))
	}
}

func TestCreateDirectedGraph(t *testing.T) {
	g := graph.New[int](graph.KindDirected)
	if g.Kind() != graph.KindDirected {
		t.Fatal("graph is expected to be undirected")
	}

	v := g.AddVertex(1)
	if v.Value != 1 {
		t.Fatalf("newly added vertex value must be 1")
	}

	e1 := g.AddEdge(1, 2)
	if e1 == nil {
		t.Fatal("got nil edge, expected valid edge")
	}

	if e1.From != 1 {
		t.Fatalf("outgoing edge vertex must be 1")
	}

	if e1.To != 2 {
		t.Fatalf("incoming edge vertex must be 2")
	}

	if e1.Weight != 0.0 {
		t.Fatal("edge weight must be 0")
	}

	if !g.VertexExists(1) {
		t.Fatal("vertex 1 must exist")
	}

	// The number of edges we have so far is 1
	if len(g.GetEdges()) != 1 {
		t.Fatal("expected number of edges in graph is 1")
	}

	// Adding same edge should not create a new edge
	e1Prime := g.AddEdge(1, 2)
	if e1.From != e1Prime.From {
		t.Fatal("e1.From and e1Prime.From are not equal")
	}
	if e1.To != e1Prime.To {
		t.Fatal("e1.To and e1Prime.To are not equal")
	}

	// The number of edges we have so far is 1
	if len(g.GetEdges()) != 1 {
		t.Fatal("expected number of edges in graph is 1")
	}

	// Non-existing vertex
	if g.VertexExists(42) {
		t.Fatal("non-existing vertex found")
	}

	if !g.EdgeExists(1, 2) {
		t.Fatal("edge (1,2) does not exist")
	}

	// Non-existing edge
	if g.EdgeExists(1, 42) {
		t.Fatal("non-existing edge (1, 42) found")
	}

	e2 := g.GetEdge(1, 42)
	if e2 != nil {
		t.Fatal("non-existing edge (1, 42) retrieved")
	}

	e3 := g.AddEdge(1, 100)
	e3.DotAttributes["color"] = "red"

	wantNumVertices := 3
	wantNumEdges := 2

	if wantNumVertices != len(g.GetVertices()) {
		t.Fatalf("want %d number of vertices, got %d", wantNumVertices, len(g.GetVertices()))
	}

	if wantNumVertices != len(g.GetVertexValues()) {
		t.Fatalf("want %d number of vertices, got %d", wantNumVertices, len(g.GetVertexValues()))
	}

	if wantNumEdges != len(g.GetEdges()) {
		t.Fatalf("want %d number of edges, got %d", wantNumEdges, len(g.GetEdges()))
	}
}

func TestDegreeOfUndirectedGraph(t *testing.T) {
	g := graph.New[int](graph.KindUndirected)
	v1 := g.AddVertex(1)
	if v1.Degree.In != 0 {
		t.Fatal("v1 in-degree must be 0")
	}
	if v1.Degree.Out != 0 {
		t.Fatal("v1 out-degree must be 0")
	}

	// Make an edge, this should result in degree changes
	g.AddEdge(1, 2)
	if v1.Degree.In != 1 {
		t.Fatal("v1 in-degree must be 1")
	}
	if v1.Degree.Out != 1 {
		t.Fatal("v1 out-degree must be 1")
	}

	// Make another edge from (2) to (3). The degree of (2) must
	// change as well.
	g.AddEdge(2, 3)
	v2 := g.GetVertex(2)
	if v2.Degree.In != 2 {
		t.Fatal("v2 in-degree must be 2")
	}
	if v2.Degree.Out != 2 {
		t.Fatal("v2 out-degree must be 2")
	}
}

func TestDegreeOfDirectedGraph(t *testing.T) {
	g := graph.New[int](graph.KindDirected)

	// Create our vertices, so we can test them later on
	v1 := g.AddVertex(1)
	v2 := g.AddVertex(2)
	v3 := g.AddVertex(3)

	if v1.Degree.In != 0 {
		t.Fatal("v1 in-degree must be 0")
	}
	if v1.Degree.Out != 0 {
		t.Fatal("v1 out-degree must be 0")
	}

	if v2.Degree.In != 0 {
		t.Fatal("v2 in-degree must be 0")
	}
	if v2.Degree.Out != 0 {
		t.Fatal("v2 out-degree must be 0")
	}

	if v3.Degree.In != 0 {
		t.Fatal("v3 in-degree must be 0")
	}
	if v3.Degree.Out != 0 {
		t.Fatal("v3 out-degree must be 0")
	}

	// Make the (v1, v2) edge
	g.AddEdge(1, 2)

	if v1.Degree.In != 0 {
		t.Fatal("v1 in-degree must be 0")
	}
	if v1.Degree.Out != 1 {
		t.Fatal("v1 out-degree must be 1")
	}

	if v2.Degree.In != 1 {
		t.Fatal("v2 in-degree must be 1")
	}
	if v2.Degree.Out != 0 {
		t.Fatal("v2 out-degree must be 0")
	}

	// Make the (v2, v3) edge
	g.AddEdge(2, 3)

	if v2.Degree.Out != 1 {
		t.Fatal("v2 out-degree must be 1")
	}
	if v2.Degree.In != 1 {
		t.Fatal("v2 in-degree must be 1")
	}

	if v3.Degree.In != 1 {
		t.Fatal("v3 in-degree must be 1")
	}
	if v3.Degree.Out != 0 {
		t.Fatal("v3 out-degree must be 0")
	}
}

func TestDeleteEdgeUndirectedGraph(t *testing.T) {
	g := graph.New[int](graph.KindUndirected)
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)
	g.AddEdge(3, 4)

	if len(g.GetVertexValues()) != 4 {
		t.Fatal("graph must have 4 vertices")
	}

	if len(g.GetEdges()) != 3 {
		t.Fatal("graph must have 3 edges")
	}

	// Delete an edge, the total number of vertices should remain
	// the same
	g.DeleteEdge(3, 4)
	if len(g.GetVertexValues()) != 4 {
		t.Fatal("graph must have 4 vertices")
	}

	if len(g.GetEdges()) != 2 {
		t.Fatal("graph must have 2 edges")
	}

	// Deleting a vertex should delete the vertex as well as the
	// edges of the vertex
	g.DeleteVertex(1)
	if len(g.GetVertexValues()) != 3 {
		t.Fatal("graph must have 3 vertices")
	}

	if len(g.GetEdges()) != 1 {
		t.Fatal("graph must have 1 edge")
	}

	// Delete a non-existing edge
	g.DeleteEdge(42, 42)
	if len(g.GetEdges()) != 1 {
		t.Fatal("graph must have 1 edge")
	}

	// Delete a non-existing vertex
	g.DeleteVertex(42)
	if len(g.GetVertexValues()) != 3 {
		t.Fatal("graph must have 3 vertices")
	}
}

func TestDeleteEdgeDirectedGraph(t *testing.T) {
	g := graph.New[int](graph.KindDirected)
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)
	g.AddEdge(3, 4)

	if len(g.GetVertexValues()) != 4 {
		t.Fatal("g2 must have 4 vertices")
	}

	if len(g.GetEdges()) != 3 {
		t.Fatal("g1 must have 3 edges")
	}

	// Delete an edge, the total number of vertices should remain
	// the same
	g.DeleteEdge(3, 4)
	if len(g.GetVertexValues()) != 4 {
		t.Fatal("graph must have 4 vertices")
	}

	if len(g.GetEdges()) != 2 {
		t.Fatal("graph must have 2 edges")
	}

	// Deleting a vertex should delete the vertex as well as the
	// edges of the vertex
	g.DeleteVertex(1)
	if len(g.GetVertexValues()) != 3 {
		t.Fatal("graph must have 3 vertices")
	}

	if len(g.GetEdges()) != 1 {
		t.Fatal("graph must have 1 edge")
	}

	// Delete a non-existing edge
	g.DeleteEdge(42, 42)
	if len(g.GetEdges()) != 1 {
		t.Fatal("graph must have 1 edge")
	}

	// Delete a non-existing vertex
	g.DeleteVertex(42)
	if len(g.GetVertexValues()) != 3 {
		t.Fatal("graph must have 3 vertices")
	}
}

func TestCloneUndirectedGraph(t *testing.T) {
	g1 := graph.New[int](graph.KindUndirected)
	g1.AddEdge(1, 2)
	g1.AddEdge(2, 3)
	g1.AddEdge(3, 4)

	v1 := g1.GetVertex(1)
	v1.DotAttributes["label"] = "v1"

	e1 := g1.GetEdge(1, 2)
	e1.DotAttributes["label"] = "e1"

	// Walk the graph, so that we build a depth-first tree
	dummyWalker := func(v *graph.Vertex[int]) error {
		return nil
	}
	if err := graph.WalkPreOrderDFS(g1, 1, dummyWalker); err != nil {
		t.Fatal(err)
	}

	// Clone the graph
	g2 := g1.Clone()

	// We expect the same number of vertices and edges
	if len(g1.GetEdges()) != len(g2.GetEdges()) {
		t.Fatal("g1 and g2 must have the same number of edges")
	}

	if len(g1.GetVertices()) != len(g2.GetVertices()) {
		t.Fatal("g1 and g2 must have the same number of vertices")
	}

	// Vertex Dot attributes should be carried over
	v1Prime := g2.GetVertex(1)
	if v1.DotAttributes["label"] != v1Prime.DotAttributes["label"] {
		t.Fatal("vertex dot attributes do not match")
	}

	v1Prime.DotAttributes["label"] = "v1Prime"
	if v1.DotAttributes["label"] == v1Prime.DotAttributes["label"] {
		t.Fatal("vertex dot attributes match")
	}

	// Edge Dot attributes should be carried over
	e1Prime := g2.GetEdge(1, 2)
	if e1.DotAttributes["label"] != e1Prime.DotAttributes["label"] {
		t.Fatal("edge dot attributes do not match")
	}

	e1Prime.DotAttributes["label"] = "e1Prime"
	if e1.DotAttributes["label"] == e1Prime.DotAttributes["label"] {
		t.Fatal("edge attributes match")
	}

	// Parents for (2) should point to different instances, but
	// they should contain the same values
	v2 := g1.GetVertex(2)
	v2Prime := g2.GetVertex(2)
	if v2.Parent == v2Prime.Parent {
		t.Fatal("v2 and v2Prime have same parent")
	}

	if v2.Parent.Value != v2Prime.Parent.Value {
		t.Fatal("v2 and v2Prime parent values mismatch")
	}
}

func TestCloneDirectedGraph(t *testing.T) {
	g1 := graph.New[int](graph.KindDirected)
	g1.AddEdge(1, 2)
	g1.AddEdge(2, 3)
	g1.AddEdge(3, 4)

	v1 := g1.GetVertex(1)
	v1.DotAttributes["label"] = "v1"

	e1 := g1.GetEdge(1, 2)
	e1.DotAttributes["label"] = "e1"

	// Walk the graph, so that we build a depth-first tree
	dummyWalker := func(v *graph.Vertex[int]) error {
		return nil
	}
	if err := graph.WalkPreOrderDFS(g1, 1, dummyWalker); err != nil {
		t.Fatal(err)
	}

	// Clone the graph
	g2 := g1.Clone()

	// We expect the same number of vertices and edges
	if len(g1.GetEdges()) != len(g2.GetEdges()) {
		t.Fatal("g1 and g2 must have the same number of edges")
	}

	if len(g1.GetVertices()) != len(g2.GetVertices()) {
		t.Fatal("g1 and g2 must have the same number of vertices")
	}

	// Vertex Dot attributes should be carried over
	v1Prime := g2.GetVertex(1)
	if v1.DotAttributes["label"] != v1Prime.DotAttributes["label"] {
		t.Fatal("vertex dot attributes do not match")
	}

	v1Prime.DotAttributes["label"] = "v1Prime"
	if v1.DotAttributes["label"] == v1Prime.DotAttributes["label"] {
		t.Fatal("vertex dot attributes match")
	}

	// Edge Dot attributes should be carried over
	e1Prime := g2.GetEdge(1, 2)
	if e1.DotAttributes["label"] != e1Prime.DotAttributes["label"] {
		t.Fatal("edge dot attributes do not match")
	}

	e1Prime.DotAttributes["label"] = "e1Prime"
	if e1.DotAttributes["label"] == e1Prime.DotAttributes["label"] {
		t.Fatal("edge attributes match")
	}

	// Parents for (2) should point to different instances, but
	// they should contain the same values
	v2 := g1.GetVertex(2)
	v2Prime := g2.GetVertex(2)
	if v2.Parent == v2Prime.Parent {
		t.Fatal("v2 and v2Prime have same parent")
	}

	if v2.Parent.Value != v2Prime.Parent.Value {
		t.Fatal("v2 and v2Prime parent values mismatch")
	}
}
