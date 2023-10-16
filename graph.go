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
