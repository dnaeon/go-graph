package graph

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

// Vertex represents a vertex in the graph
type Vertex[T comparable] struct {
	// Value contains the value for the vertex
	Value T

	// Color represents the color the vertex is painted with
	Color Color
}

// NewVertex creates a new vertex with the given value
func NewVertex[T comparable](value T) *Vertex[T] {
	v := &Vertex[T]{
		Value: value,
		Color: White,
	}

	return v
}

// Edge represents an edge connecting two vertices in the graph
type Edge[T comparable] struct {
	// From represents the source vertex of the edge
	From T

	// To represents the destination vertex of the edge
	To T
}

// NewEdge creates an edge, which connects the given vertices
func NewEdge[T comparable](from, to T) *Edge[T] {
	e := &Edge[T]{
		From: from,
		To:   to,
	}

	return e
}

// Graph represents an undirected graph
type Graph[T comparable] struct {
	// The set of vertices in the graph
	vertices map[T]*Vertex[T]

	// The set of edges in the graph
	edges []*Edge[T]

	// The adjacency lists for our vertices
	adjacencyLists map[T][]T
}

// NewGraph creates a new graph
func NewGraph[T comparable]() *Graph[T] {
	g := &Graph[T]{
		vertices:       make(map[T]*Vertex[T]),
		edges:          make([]*Edge[T], 0),
		adjacencyLists: make(map[T][]T),
	}

	return g
}

// GetVertex returns the vertex associated with the given value
func (g *Graph[T]) GetVertex(value T) *Vertex[T] {
	return g.vertices[value]
}

// VertexExists returns a boolean indicating whether a vertex with the
// given value exists
func (g *Graph[T]) VertexExists(value T) bool {
	_, exists := g.vertices[value]
	return exists
}

// GetVertices returns the set of vertices in the graph
func (g *Graph[T]) GetVertices() []*Vertex[T] {
	result := make([]*Vertex[T], 0)
	for _, v := range g.vertices {
		result = append(result, v)
	}

	return result
}

// GetVertexValues returns the set of vertex values
func (g *Graph[T]) GetVertexValues() []T {
	result := make([]T, 0)
	for k := range g.vertices {
		result = append(result, k)
	}

	return result
}

// GetEdges returns the set of edges in the graph
func (g *Graph[T]) GetEdges() []*Edge[T] {
	return g.edges
}

// GetNeighbours returns the list of direct neighbours of V
func (g *Graph[T]) GetNeighbours(v T) []T {
	return g.adjacencyLists[v]
}

// GetNeighbourVertices returns the list of neighbour vertices of V
func (g *Graph[T]) GetNeighbourVertices(v T) []*Vertex[T] {
	neighbours := g.GetNeighbours(v)
	result := make([]*Vertex[T], 0)
	for _, u := range neighbours {
		result = append(result, g.GetVertex(u))
	}

	return result
}

// AddVertex adds a vertex to the graph
func (g *Graph[T]) AddVertex(value T) *Vertex[T] {
	if g.VertexExists(value) {
		return g.GetVertex(value)
	}

	vertex := NewVertex(value)
	g.vertices[value] = vertex

	return vertex
}

// GetEdge returns the edge connecting the two vertices
func (g *Graph[T]) GetEdge(from, to T) *Edge[T] {
	for _, e := range g.edges {
		if e.From == from && e.To == to {
			return e
		}
	}

	return nil
}

// EdgeExists returns a boolean indicating whether an edge between two
// vertices exists.
func (g *Graph[T]) EdgeExists(from, to T) bool {
	for _, e := range g.edges {
		if e.From == from && e.To == to {
			return true
		}
	}

	return false
}

// AddEdge adds an edge between two vertices in the graph
func (g *Graph[T]) AddEdge(from, to T) *Edge[T] {
	if g.EdgeExists(from, to) {
		return g.GetEdge(from, to)
	}

	g.AddVertex(from)
	g.AddVertex(to)

	// Create the edge
	e := NewEdge(from, to)
	g.edges = append(g.edges, e)

	// Update the adjacency lists
	g.adjacencyLists[from] = append(g.adjacencyLists[from], to)
	g.adjacencyLists[to] = append(g.adjacencyLists[to], from)

	return e
}
