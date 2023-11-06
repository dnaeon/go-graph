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
	"errors"
	"slices"
)

// Color represents the color with which a vertex is painted
type Color int

const (
	// White color means that the vertex has not been seen
	White Color = iota

	// Gray color means that the vertex is seen for the first time
	Gray

	// Black color means that the vertex has already been explored
	Black
)

// GraphKind represents the kind of the graph
type GraphKind int

const (
	// A kind which represents a directed graph
	KindDirected GraphKind = iota

	// A kind which represents an undirected graph
	KindUndirected
)

// DotAttributes contains the map of key/value pairs, which can be
// associated with vertices and edges.
type DotAttributes map[string]string

// Degree represents the number of incoming and outgoing edges of a
// vertex
type Degree struct {
	// The number of incoming edges
	In int

	// The number of outgoing edges
	Out int
}

// Vertex represents a vertex in the graph
type Vertex[T comparable] struct {
	// Value contains the value for the vertex
	Value T

	// Color represents the color the vertex is painted with
	Color Color

	// DistanceFromSource returns the distance of this vertex from
	// the source vertex.  This field is calculated after
	// performing a DFS or BFS.
	DistanceFromSource float64

	// Parent represents the parent vertex, which is calculated
	// during walking of the graph.  The resulting relationships
	// established by this field construct the DFS-tree, BFS-tree
	// or shortest-path tree, depending on how we walked the
	// graph.
	Parent *Vertex[T]

	// DotAttributes represents the list of attributes associated
	// with the vertex. The attributes will be used when
	// generating the Dot representation of the graph.
	DotAttributes DotAttributes

	// Degree represents the degree of the vertex
	Degree Degree
}

// NewVertex creates a new vertex with the given value
func NewVertex[T comparable](value T) *Vertex[T] {
	v := &Vertex[T]{
		Value:              value,
		Color:              White,
		DistanceFromSource: 0.0,
		Parent:             nil,
		DotAttributes:      make(DotAttributes),
		Degree:             Degree{In: 0, Out: 0},
	}

	return v
}

// Edge represents an edge connecting two vertices in the graph
type Edge[T comparable] struct {
	// From represents the source vertex of the edge
	From T

	// To represents the destination vertex of the edge
	To T

	// Weight represents the edge weight
	Weight float64

	// DotAttributes represents the list of attributes associated
	// with the edge. The attributes will be used when
	// generating the Dot representation of the graph.
	DotAttributes DotAttributes
}

// NewEdge creates an edge, which connects the given vertices
func NewEdge[T comparable](from, to T) *Edge[T] {
	e := &Edge[T]{
		From:          from,
		To:            to,
		Weight:        0.0,
		DotAttributes: make(DotAttributes),
	}

	return e
}

// WalkFunc is a function which receives a vertex while traversing the
// graph
type WalkFunc[T comparable] func(v *Vertex[T]) error

// ErrStopWalking is returned by WalkFunc to signal that further
// walking of the graph should be stopped.
var ErrStopWalking = errors.New("walking stopped")

// Collector provides an easy way to collect vertices while walking a
// graph
type Collector[T comparable] struct {
	items []*Vertex[T]
}

// NewCollector creates a new collector
func NewCollector[T comparable]() *Collector[T] {
	c := &Collector[T]{
		items: make([]*Vertex[T], 0),
	}

	return c
}

// WalkFunc collects the vertices it visits
func (c *Collector[T]) WalkFunc(v *Vertex[T]) error {
	c.items = append(c.items, v)

	return nil
}

// Get returns the collected vertices
func (c *Collector[T]) Get() []*Vertex[T] {
	return c.items
}

// Reset resets the set of collected vertices
func (c *Collector[T]) Reset() {
	c.items = make([]*Vertex[T], 0)
}

// Graph represents a graph which establishes relationships between
// various objects.
type Graph[T comparable] interface {
	// Kind returns the kind of the graph
	Kind() GraphKind

	// AddVertex adds a new vertex to the graph
	AddVertex(v T) *Vertex[T]

	// GetVertex returns the vertex associated with the given value
	GetVertex(v T) *Vertex[T]

	// DeleteVertex deletes the vertex associated with the given value
	DeleteVertex(v T)

	// VertexExists is a predicate for testing whether a vertex
	// associated with the value exists.
	VertexExists(v T) bool

	// GetVertices returns all vertices in the graph
	GetVertices() []*Vertex[T]

	// GetVertexValues returns the values associated with all
	// vertices from the graph
	GetVertexValues() []T

	// AddEdge creates a new edge connecting `from` and `to`
	// vertices
	AddEdge(from, to T) *Edge[T]

	// AddWeightedEdge creates a new edge with the given weight
	AddWeightedEdge(from, to T, weight float64) *Edge[T]

	// GetEdge returns the edge, which connects `from` and `to`
	// vertices
	GetEdge(from, to T) *Edge[T]

	// DeleteEdge deletes the edge which connects `from` and `to`
	// vertices
	DeleteEdge(from, to T)

	// EdgeExists is a predicate for testing whether an edge
	// between `from` and `to` exists
	EdgeExists(from, to T) bool

	// GetEdges returns all edges from the graph
	GetEdges() []*Edge[T]

	// GetNeighbours returns the neighbours of V as values
	GetNeighbours(v T) []T

	// GetNeighbourVertices returns the neighbours of V as
	// vertices
	GetNeighbourVertices(v T) []*Vertex[T]

	// ResetVertexAttributes resets the attributes for all
	// vertices in the graph
	ResetVertexAttributes()

	// NewCollector creates and returns a new collector
	NewCollector() *Collector[T]

	// Clone creates a new copy of the graph
	Clone() Graph[T]
}

// UndirectedGraph represents an undirected graph
type UndirectedGraph[T comparable] struct {
	// The set of vertices in the graph
	vertices map[T]*Vertex[T]

	// The set of edges in the graph
	edges []*Edge[T]

	// The adjacency lists for our vertices
	adjacencyLists map[T][]T

	// The kind of the graph
	kind GraphKind
}

// NewGraph creates a new graph
func New[T comparable](kind GraphKind) Graph[T] {
	g := UndirectedGraph[T]{
		vertices:       make(map[T]*Vertex[T]),
		edges:          make([]*Edge[T], 0),
		adjacencyLists: make(map[T][]T),
		kind:           kind,
	}

	if kind == KindDirected {
		return &DirectedGraph[T]{
			UndirectedGraph: g,
		}
	}

	return &g
}

// Clone creates a new copy of the graph.
func (g *UndirectedGraph[T]) Clone() Graph[T] {
	newVertices := make(map[T]*Vertex[T])
	newEdges := make([]*Edge[T], 0)
	newAdjacencyLists := make(map[T][]T)

	// Clone vertices
	for k, v := range g.vertices {
		dotAttributes := make(DotAttributes)
		for dotK, dotV := range v.DotAttributes {
			dotAttributes[dotK] = dotV
		}

		newV := &Vertex[T]{
			Value:              v.Value,
			Color:              v.Color,
			DistanceFromSource: v.DistanceFromSource,
			Parent:             nil, // Parent will be populated a bit later
			DotAttributes:      dotAttributes,
			Degree:             Degree{In: v.Degree.In, Out: v.Degree.Out},
		}
		newVertices[k] = newV
	}

	// Populate parent field, now that we have all vertices created
	for k, v := range g.vertices {
		if v.Parent == nil {
			continue
		}
		newParent := newVertices[v.Parent.Value]
		newVertices[k].Parent = newParent
	}

	// Clone edges
	for _, e := range g.edges {
		dotAttributes := make(DotAttributes)
		for dotK, dotV := range e.DotAttributes {
			dotAttributes[dotK] = dotV
		}
		newE := &Edge[T]{
			From:          e.From,
			To:            e.To,
			Weight:        e.Weight,
			DotAttributes: dotAttributes,
		}
		newEdges = append(newEdges, newE)
	}

	// Clone adjacency lists
	for v, adjList := range g.adjacencyLists {
		newAdjList := make([]T, 0)
		for _, u := range adjList {
			newAdjList = append(newAdjList, u)
		}
		newAdjacencyLists[v] = newAdjList
	}

	// Create the new graph
	g1 := UndirectedGraph[T]{
		vertices:       newVertices,
		edges:          newEdges,
		adjacencyLists: newAdjacencyLists,
		kind:           g.kind,
	}

	if g.kind == KindDirected {
		return &DirectedGraph[T]{
			UndirectedGraph: g1,
		}
	}

	return &g1
}

// NewCollector creates a new collector
func (g *UndirectedGraph[T]) NewCollector() *Collector[T] {
	c := NewCollector[T]()

	return c
}

// Kind returns the kind of the graph
func (g *UndirectedGraph[T]) Kind() GraphKind {
	return g.kind
}

// ResetVertexAttributes resets the attributes for each vertex in the
// graph
func (g *UndirectedGraph[T]) ResetVertexAttributes() {
	for _, v := range g.vertices {
		v.Color = White
		v.DistanceFromSource = 0.0
		v.Parent = nil
	}
}

// GetVertex returns the vertex associated with the given value
func (g *UndirectedGraph[T]) GetVertex(value T) *Vertex[T] {
	return g.vertices[value]
}

// VertexExists returns a boolean indicating whether a vertex with the
// given value exists
func (g *UndirectedGraph[T]) VertexExists(value T) bool {
	_, exists := g.vertices[value]
	return exists
}

// GetVertices returns the set of vertices in the graph
func (g *UndirectedGraph[T]) GetVertices() []*Vertex[T] {
	result := make([]*Vertex[T], 0)
	for _, v := range g.vertices {
		result = append(result, v)
	}

	return result
}

// GetVertexValues returns the set of vertex values
func (g *UndirectedGraph[T]) GetVertexValues() []T {
	result := make([]T, 0)
	for k := range g.vertices {
		result = append(result, k)
	}

	return result
}

// GetEdges returns the set of edges in the graph
func (g *UndirectedGraph[T]) GetEdges() []*Edge[T] {
	return g.edges
}

// GetNeighbours returns the list of direct neighbours of V
func (g *UndirectedGraph[T]) GetNeighbours(v T) []T {
	return g.adjacencyLists[v]
}

// GetNeighbourVertices returns the list of neighbour vertices of V
func (g *UndirectedGraph[T]) GetNeighbourVertices(v T) []*Vertex[T] {
	neighbours := g.GetNeighbours(v)
	result := make([]*Vertex[T], 0)
	for _, u := range neighbours {
		result = append(result, g.GetVertex(u))
	}

	return result
}

// AddVertex adds a vertex to the graph
func (g *UndirectedGraph[T]) AddVertex(value T) *Vertex[T] {
	if g.VertexExists(value) {
		return g.GetVertex(value)
	}

	vertex := NewVertex(value)
	g.vertices[value] = vertex

	return vertex
}

// DeleteVertex removes a vertex from the graph
func (g *UndirectedGraph[T]) DeleteVertex(v T) {
	if !g.VertexExists(v) {
		return
	}

	// Delete edges in the graph, which connect V with any other
	// vertex in the graph
	for _, e := range g.GetEdges() {
		if e.From == v || e.To == v {
			g.DeleteEdge(e.From, e.To)
		}
	}

	// Delete the vertex itself
	delete(g.vertices, v)
}

// GetEdge returns the edge connecting the two vertices
func (g *UndirectedGraph[T]) GetEdge(from, to T) *Edge[T] {
	for _, e := range g.edges {
		if (e.From == from && e.To == to) || (e.From == to && e.To == from) {
			return e
		}
	}

	return nil
}

// DeleteEdge deletes the edge, which connects the `from` and `to`
// vertices
func (g *UndirectedGraph[T]) DeleteEdge(from, to T) {
	if !g.EdgeExists(from, to) {
		return
	}

	for idx, e := range g.edges {
		if (e.From == from && e.To == to) || (e.From == to && e.To == from) {
			g.edges = slices.Delete(g.edges, idx, idx+1)
		}
	}

	// Update the adjacency lists
	for idx, v := range g.adjacencyLists[from] {
		if v == to {
			g.adjacencyLists[from] = slices.Delete(g.adjacencyLists[from], idx, idx+1)
		}
	}

	for idx, v := range g.adjacencyLists[to] {
		if v == from {
			g.adjacencyLists[to] = slices.Delete(g.adjacencyLists[to], idx, idx+1)
		}
	}

	// Update degree
	fromV := g.GetVertex(from)
	fromV.Degree.In -= 1
	fromV.Degree.Out -= 1

	toV := g.GetVertex(to)
	toV.Degree.In -= 1
	toV.Degree.Out -= 1
}

// EdgeExists returns a boolean indicating whether an edge between two
// vertices exists.
func (g *UndirectedGraph[T]) EdgeExists(from, to T) bool {
	e := g.GetEdge(from, to)
	if e != nil {
		return true
	}

	return false
}

// AddEdge adds an edge between two vertices in the graph
func (g *UndirectedGraph[T]) AddEdge(from, to T) *Edge[T] {
	if g.EdgeExists(from, to) {
		return g.GetEdge(from, to)
	}

	fromV := g.AddVertex(from)
	toV := g.AddVertex(to)

	// Create the edge
	e := NewEdge(from, to)
	g.edges = append(g.edges, e)

	// Update the adjacency lists
	g.adjacencyLists[from] = append(g.adjacencyLists[from], to)
	g.adjacencyLists[to] = append(g.adjacencyLists[to], from)

	// Update the vertices degree
	fromV.Degree.In += 1
	fromV.Degree.Out += 1
	toV.Degree.In += 1
	toV.Degree.Out += 1

	return e
}

// AddWeightedEdge adds an edge between two vertices and sets weight
// for the edge
func (g *UndirectedGraph[T]) AddWeightedEdge(from, to T, weight float64) *Edge[T] {
	e := g.AddEdge(from, to)
	e.Weight = weight

	return e
}

// DirectedGraph represents a directed graph
type DirectedGraph[T comparable] struct {
	UndirectedGraph[T]
}

// AddEdge adds an edge between two vertices in the graph
func (g *DirectedGraph[T]) AddEdge(from, to T) *Edge[T] {
	if g.EdgeExists(from, to) {
		return g.GetEdge(from, to)
	}

	fromV := g.AddVertex(from)
	toV := g.AddVertex(to)

	// Create the edge
	e := NewEdge(from, to)
	g.edges = append(g.edges, e)

	// Update the adjacency lists
	g.adjacencyLists[from] = append(g.adjacencyLists[from], to)

	// Update vertices degree
	fromV.Degree.Out += 1
	toV.Degree.In += 1

	return e
}

// EdgeExists returns a boolean indicating whether an edge between two
// vertices exists.
func (g *DirectedGraph[T]) EdgeExists(from, to T) bool {
	e := g.GetEdge(from, to)
	if e != nil {
		return true
	}

	return false
}

// GetEdge returns the edge connecting the two vertices
func (g *DirectedGraph[T]) GetEdge(from, to T) *Edge[T] {
	for _, e := range g.edges {
		if e.From == from && e.To == to {
			return e
		}
	}

	return nil
}

// DeleteEdge deletes the edge, which connects the `from` and `to`
// vertices
func (g *DirectedGraph[T]) DeleteEdge(from, to T) {
	if !g.EdgeExists(from, to) {
		return
	}

	// Remove the edge itself
	for idx, e := range g.edges {
		if e.From == from && e.To == to {
			g.edges = slices.Delete(g.edges, idx, idx+1)
		}
	}

	// Update the adjacency lists
	for idx, v := range g.adjacencyLists[from] {
		if v == to {
			g.adjacencyLists[from] = slices.Delete(g.adjacencyLists[from], idx, idx+1)
		}
	}

	fromV := g.GetVertex(from)
	fromV.Degree.Out -= 1

	toV := g.GetVertex(to)
	toV.Degree.In -= 1
}
