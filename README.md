# go-graph

[![Build Status](https://github.com/dnaeon/go-graph/actions/workflows/test.yaml/badge.svg)](https://github.com/dnaeon/go-graph/actions/workflows/test.yaml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/gopkg.in/dnaeon/go-graph.v1.svg)](https://pkg.go.dev/gopkg.in/dnaeon/go-graph.v1)
[![Go Report Card](https://goreportcard.com/badge/gopkg.in/dnaeon/go-graph.v1)](https://goreportcard.com/report/gopkg.in/dnaeon/go-graph.v1)
[![codecov](https://codecov.io/gh/dnaeon/go-graph/branch/v1/graph/badge.svg)](https://codecov.io/gh/dnaeon/go-graph)

A simple and generic library for working with
[Graphs](https://en.wikipedia.org/wiki/Graph_(discrete_mathematics))
in Go.

![Example Directed Graph](./images/directed-g.svg)

## Installation

Execute the following command.

``` shell
go get -v gopkg.in/dnaeon/go-graph.v1
```

## Usage

Consider the following undirected graph.

![Example Undirected Graph](./images/undirected-g.svg)

This snippet creates the graph and performs DFS traversal on it.

``` go
package main

import (
	"fmt"

	"gopkg.in/dnaeon/go-graph.v1"
)

func main() {
	g := graph.New[int](graph.KindUndirected)
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 4)
	g.AddEdge(3, 4)
	g.AddEdge(4, 5)

	walker := func(v *graph.Vertex[int]) error {
		fmt.Println(v.Value)
		return nil
	}

	fmt.Printf("DFS (pre-order) from (1):\n")
	if err := graph.WalkPreOrderDFS(g, 1, walker); err != nil {
		fmt.Printf("DFS (pre-order): %s\n", err)
	}

	fmt.Printf("\nDFS (pre-order) from (3):\n")
	if err := graph.WalkPreOrderDFS(g, 3, walker); err != nil {
		fmt.Printf("DFS (pre-order): %s\n", err)
	}

	fmt.Printf("\nDFS (post-order) from (1):\n")
	if err := graph.WalkPostOrderDFS(g, 1, walker); err != nil {
		fmt.Printf("DFS (pre-order): %s\n", err)
	}

	fmt.Printf("\nDFS (post-order) from (3):\n")
	if err := graph.WalkPostOrderDFS(g, 3, walker); err != nil {
		fmt.Printf("DFS (post-order): %s\n", err)
	}
}
```

This example creates a directed graph, and then walks over the
vertices in topological order.

![Example Directed Graph](./images/directed-g.svg)

``` go
package main

import (
	"fmt"

	"gopkg.in/dnaeon/go-graph.v1"
)

func main() {
	g := graph.New[int](graph.KindDirected)
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 4)
	g.AddEdge(3, 4)
	g.AddEdge(4, 5)

	walker := func(v *graph.Vertex[int]) error {
		fmt.Println(v.Value)
		return nil
	}

	fmt.Println("Topo order:")
	if err := graph.WalkTopoOrder(g, walker); err != nil {
		fmt.Printf("WalkTopoOrder: %s\n", err)
	}
}
```

Output from above code looks like this.

``` text
Topo order:
5
4
3
2
1
```

A topological sort is only possible on a Directed Acyclic Graph (DAG).  When the
graph contains a cycle, `WalkTopoOrder` fails with a `*CycleError`, which wraps
`ErrCycleDetected` and carries one cycle as a witness. This makes it a
convenient way to answer _"is this graph a DAG?"_.

To find __all__ the cycles in a graph, rather than a single one, use
`StronglyConnectedComponents` and `WalkCycles`. A strongly connected component
with more than one vertex (or a single vertex with a self-loop) is a cyclic
cluster, and `WalkCycles` walks over the cycles, invoking the walker once per
cluster with a representative cycle.

The following code creates a directed graph with two independent cycles, then
uses `WalkCycles` to paint the vertices and edges forming each cycle in red, so
they stand out when the graph is rendered. See the
[examples/cycles](./examples/cycles) example for the full source.

``` go
package main

import (
	"fmt"
	"os"

	"gopkg.in/dnaeon/go-graph.v1"
)

func main() {
	g := graph.New[int](graph.KindDirected)

	// Acyclic spine with a shared successor.
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 4)
	g.AddEdge(3, 4)
	g.AddEdge(4, 5)

	// Two-node cycle: 10 <-> 11.
	g.AddEdge(1, 10)
	g.AddEdge(10, 11)
	g.AddEdge(11, 10)

	// Three-node cycle: 20 -> 21 -> 22 -> 20.
	g.AddEdge(1, 20)
	g.AddEdge(20, 21)
	g.AddEdge(21, 22)
	g.AddEdge(22, 20)

	// Paint the vertices and edges forming each cycle in red.
	cycleWalker := func(cycle []*graph.Vertex[int]) error {
		for i, v := range cycle {
			v.DotAttributes["color"] = "red"
			v.DotAttributes["fillcolor"] = "red"

			if i > 0 {
				edge := g.GetEdge(cycle[i-1].Value, v.Value)
				edge.DotAttributes["color"] = "red"
			}
		}

		return nil
	}
	if err := graph.WalkCycles(g, cycleWalker); err != nil {
		fmt.Printf("WalkCycles: %s\n", err)
		return
	}

	// Emit the Dot representation, with the cycles painted red.
	if err := graph.WriteDot(g, os.Stdout); err != nil {
		fmt.Println(err)
	}
}
```

Render the output with `graphviz`, e.g.

``` shell
go run ./examples/cycles/main.go | dot -T svg -o cycles.svg
```

The cyclic vertices and edges (`10 <-> 11` and `20 -> 21 -> 22 -> 20`)
are painted red, while the acyclic part of the graph keeps its default
color.

![Example Directed Graph with Cycles Painted](./images/cycles.svg)

Generate the [Dot
representation](https://graphviz.org/doc/info/lang.html) for a graph.

The following code generates the Dot representation of the directed
graph used in the previous example.

``` go
package main

import (
	"fmt"
	"os"

	"gopkg.in/dnaeon/go-graph.v1"
)

func main() {
	g := graph.New[int](graph.KindDirected)
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 4)
	g.AddEdge(3, 4)
	g.AddEdge(4, 5)

	if err := graph.WriteDot(g, os.Stdout); err != nil {
		fmt.Println(err)
	}
}
```

Consider the following undirected weighted graph. The edges in this
graph represent the distance between vertices in the graph.

![Example Undirected Weighted Graph](./images/undirected-weighted-1.svg)

The following code will print the shortest path between the given
source and destination vertices, then paint the visited vertices in
green, and finally print the Dot representation of the graph.


``` go
package main

import (
	"fmt"
	"os"

	"gopkg.in/dnaeon/go-graph.v1"
)

func main() {
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

	var prev *graph.Vertex[int]
	walker := func(v *graph.Vertex[int]) error {
		// Paint vertices, which form the shortest path in
		// green
		v.DotAttributes["color"] = "green"
		v.DotAttributes["fillcolor"] = "green"

		if prev != nil {
			edge := g.GetEdge(prev.Value, v.Value)
			edge.DotAttributes["label"] = fmt.Sprintf("%d", int(v.DistanceFromSource))
		}

		prev = v

		fmt.Println(v.Value)
		return nil
	}

	fmt.Printf("Shortest path from (1) to (8):\n")
	if err := graph.WalkShortestPath(g, 1, 8, walker); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\nDot representation of graph:\n\n")
	if err := graph.WriteDot(g, os.Stdout); err != nil {
		fmt.Println(err)
	}
}
```

This is what the shortest path between vertices `(1)` and `(8)` looks
like in Dot representation.

![Example Undirected Weighted Graph - Painted](./images/undirected-weighted-2.svg)

Graphs can also be rendered using
[go-echarts](https://github.com/go-echarts/go-echarts). The
[examples/shortest-path-echarts](./examples/shortest-path-echarts) example is
similar to the shortest-path example above, but instead of rendering in Dot
representation it renders the graph using `echarts`.

![Shortest Path using go-echarts](./images/shortest-path-echarts.png)

Make sure to also check the included [test cases](./graph_test.go) and
[examples](./examples) directory from this repo.

## Tests

Run the tests.

``` shell
make test
```

## License

`go-graph` is Open Source and licensed under the [BSD
License](http://opensource.org/licenses/BSD-2-Clause).
