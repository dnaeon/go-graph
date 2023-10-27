package graph

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"gopkg.in/dnaeon/go-deque.v1"
	"gopkg.in/dnaeon/go-priorityqueue.v1"
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
	Directed GraphKind = iota

	// A kind which represents an undirected graph
	Undirected
)

// DotAttributes contains the map of key/value pairs, which can be
// associated with vertices and edges.
type DotAttributes map[string]string

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
}

// NewVertex creates a new vertex with the given value
func NewVertex[T comparable](value T) *Vertex[T] {
	v := &Vertex[T]{
		Value:              value,
		Color:              White,
		DistanceFromSource: 0.0,
		Parent:             nil,
		DotAttributes:      make(DotAttributes),
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

// Graph represents a graph which establishes relationships between
// various objects.
type Graph[T comparable] interface {
	// Kind returns the kind of the graph
	Kind() GraphKind

	// AddVertex adds a new vertex to the graph
	AddVertex(v T) *Vertex[T]

	// GetVertex returns the vertex associated with the given value
	GetVertex(v T) *Vertex[T]

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

	// GetUnreachableVertices returns the vertices, which are
	// unreachable from the given source vertex.
	GetUnreachableVertices(source T) []*Vertex[T]

	// WalkDFS performs a Depth-first Search (DFS) traversal of
	// the graph, starting from the given source vertex
	WalkDFS(source T, walkFunc WalkFunc[T]) error

	// WalkBFS performs a Breadth-first Search (DFS) traversal of
	// the graph, starting from the given source vertex
	WalkBFS(source T, walkFunc WalkFunc[T]) error

	// WalkShortestPath walks over the shortest path, which
	// connects the given source and destination vertices.
	WalkShortestPath(source T, dest T, walkFunc WalkFunc[T]) error

	// WalkUnreachableVertices walks over the vertices, which are
	// unreachable from the given source vertex
	WalkUnreachableVertices(source T, walkFunc WalkFunc[T]) error

	// ShortestPath returns the list of vertices, which form the
	// shortest path connecting the given source and destination
	// vertices.
	ShortestPath(source T, dest T) ([]*Vertex[T], error)

	// WriteDot formats the graph in Dot representation
	WriteDot(w io.Writer) error
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
	g := &UndirectedGraph[T]{
		vertices:       make(map[T]*Vertex[T]),
		edges:          make([]*Edge[T], 0),
		adjacencyLists: make(map[T][]T),
		kind:           kind,
	}

	return g
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

// GetEdge returns the edge connecting the two vertices
func (g *UndirectedGraph[T]) GetEdge(from, to T) *Edge[T] {
	for _, e := range g.edges {
		if (e.From == from && e.To == to) || (e.From == to && e.To == from) {
			return e
		}
	}

	return nil
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

// AddWeightedEdge adds an edge between two vertices and sets weight
// for the edge
func (g *UndirectedGraph[T]) AddWeightedEdge(from, to T, weight float64) *Edge[T] {
	e := g.AddEdge(from, to)
	e.Weight = weight

	return e
}

// WalkDFS performs Depth-first Search (DFS) traversal of the graph,
// starting from the given source vertex.
func (g *UndirectedGraph[T]) WalkDFS(source T, walkFunc WalkFunc[T]) error {
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
			return err
		}

		// Visit the neigbours of V
		neighbours := g.GetNeighbourVertices(v.Value)
		for _, u := range neighbours {
			// First time seeing this neighbour vertex, push it to the stack
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

// WalkBFS performs Breadth-first Search (BFS) traversal of the graph,
// starting from the given source vertex.
func (g *UndirectedGraph[T]) WalkBFS(source T, walkFunc WalkFunc[T]) error {
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
			return err
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

// WalkUnreachableVertices walks over the vertices which are
// unreachable from the given source vertex
func (g *UndirectedGraph[T]) WalkUnreachableVertices(source T, walkFunc WalkFunc[T]) error {
	// In order to find all unreachable vertices we will first DFS
	// traverse the graph.  The vertices which remain White after
	// we've walked the graph are unreachable from the source
	// vertex.
	dummyDfsWalk := func(v *Vertex[T]) error {
		return nil
	}

	if err := g.WalkDFS(source, dummyDfsWalk); err != nil {
		return err
	}

	for _, v := range g.vertices {
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

// GetUnreachableVertices returns the list of vertices, which are
// unreachable from a given source vertex
func (g *UndirectedGraph[T]) GetUnreachableVertices(source T) []*Vertex[T] {
	result := make([]*Vertex[T], 0)
	walkFunc := func(v *Vertex[T]) error {
		result = append(result, v)
		return nil
	}

	if err := g.WalkUnreachableVertices(source, walkFunc); err != nil {
		return nil
	}

	return result
}

// Initializes the source vertex as part of Dijkstra's algorithm
func (g *UndirectedGraph[T]) initializeSourceVertex(source T) error {
	if !g.VertexExists(source) {
		return fmt.Errorf("Source vertex %v not found in graph", source)
	}

	// Set tentative distance for all vertices
	g.ResetVertexAttributes()
	for _, v := range g.vertices {
		v.DistanceFromSource = math.Inf(1)
	}

	// Initialize source vertex
	srcV := g.GetVertex(source)
	srcV.DistanceFromSource = 0.0

	return nil
}

// Relaxes the edge as part of Dijkstra's algorithm
func (g *UndirectedGraph[T]) relaxEdge(from, to T) error {
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
func (g *UndirectedGraph[T]) WalkDijkstra(source T, walkFunc WalkFunc[T]) error {
	if err := g.initializeSourceVertex(source); err != nil {
		return err
	}

	// Enqueue all vertices
	queue := priorityqueue.New[*Vertex[T], float64](priorityqueue.MinHeap)
	for _, v := range g.vertices {
		queue.Put(v, v.DistanceFromSource)
	}

	for !queue.IsEmpty() {
		item := queue.Get()
		v := item.Value
		// Relax edges connecting V and it's neighbours
		for _, u := range g.GetNeighbourVertices(v.Value) {
			oldDist := u.DistanceFromSource
			if err := g.relaxEdge(v.Value, u.Value); err != nil {
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

// Reverses the given slice of items
func Reverse[T any](items []T) {
	size := len(items)
	for i := 0; i < size/2; i++ {
		k := size - i - 1
		items[i], items[k] = items[k], items[i]
	}
}

// WalkShortestPath yields the vertices which represent the shortest
// path between SOURCE and DEST.
func (g *UndirectedGraph[T]) WalkShortestPath(source T, dest T, walkFunc WalkFunc[T]) error {
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

	if err := g.WalkDijkstra(source, walker); err != nil {
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

	Reverse(result)
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

// ShortestPath returns the list of vertices, which represent the
// shortest path between a given SOURCE and DEST vertex using
// Dijkstra's algorithim.
func (g *UndirectedGraph[T]) ShortestPath(source T, dest T) ([]*Vertex[T], error) {
	result := make([]*Vertex[T], 0)
	walker := func(v *Vertex[T]) error {
		result = append(result, v)
		return nil
	}

	if err := g.WalkShortestPath(source, dest, walker); err != nil {
		return nil, err
	}

	return result, nil
}

// DefaultNodeAttributes represents the map of default attributes to
// be applied for al nodes when representing the graph in Dot format.
var DotDefaultNodeAttributes = DotAttributes{
	"color":     "lightblue",
	"fillcolor": "lightblue",
	"fontcolor": "black",
	"shape":     "record",
	"style":     "filled, rounded",
}

// DotDefaultEdgeAttributes represents the map of default attributes
// to be applied for all edges when representing the graph in Dot
// format.
var DotDefaultEdgeAttributes = DotAttributes{
	"color": "black",
}

// dotId returns the unique node id, which is used when generating the
// graph representation in Dot.
func dotId(v any) int64 {
	addr := fmt.Sprintf("%p", v)
	id, err := strconv.ParseInt(addr[2:], 16, 64)
	if err != nil {
		panic(err)
	}

	return id
}

// formatDotAttributes formats the given map of attributes in Dot
// format
func formatDotAttributes(items DotAttributes) string {
	attrs := ""
	for k, v := range items {
		attrs += fmt.Sprintf("%s=%q ", k, v)
	}

	return strings.TrimRight(attrs, " ")
}

// WriteDot generates the Dot representation of the graph
func (g *UndirectedGraph[T]) WriteDot(w io.Writer) error {
	var graphKind string
	if g.kind == Undirected {
		graphKind = "graph"
	} else {
		graphKind = "digraph"
	}
	if _, err := fmt.Fprintf(w, "strict %s {\n", graphKind); err != nil {
		return err
	}

	// Default node attributes
	if _, err := fmt.Fprintf(w, "\tnode [%s]\n", formatDotAttributes(DotDefaultNodeAttributes)); err != nil {
		return err
	}

	// Default edge attributes
	if _, err := fmt.Fprintf(w, "\tedge [%s]\n", formatDotAttributes(DotDefaultEdgeAttributes)); err != nil {
		return err
	}

	for _, v := range g.vertices {
		// Do we have a label?
		_, ok := v.DotAttributes["label"]
		if !ok {
			v.DotAttributes["label"] = fmt.Sprintf("%v", v.Value)
		}

		_, err := fmt.Fprintf(w, "\t%d [%s]\n", dotId(v), formatDotAttributes(v.DotAttributes))
		if err != nil {
			return err
		}

		for _, u := range g.GetNeighbourVertices(v.Value) {
			e := g.GetEdge(u.Value, v.Value)
			var edgeArrow string
			if g.kind == Undirected {
				edgeArrow = "--"
			} else {
				edgeArrow = "->"
			}

			if _, err := fmt.Fprintf(w, "\t%d %s %d [%s]\n", dotId(v), edgeArrow, dotId(u), formatDotAttributes(e.DotAttributes)); err != nil {
				return err
			}
		}
	}

	if _, err := fmt.Fprintln(w, "}"); err != nil {
		return err
	}
	return nil
}
