// Copyright (c) 2026 Marin Atanasov Nikolov <dnaeon@gmail.com>
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

package main

import (
	"fmt"
	"os"

	"gopkg.in/dnaeon/go-graph.v1"
)

func main() {
	// A directed graph that is mostly a DAG, but contains two independent
	// cycles:
	//
	//   * 1 -> 2 -> 4, 1 -> 3 -> 4, 4 -> 5 is an acyclic spine with a
	//     shared successor (a diamond).
	//   * 10 <-> 11 is a two-node mutual dependency cycle.
	//   * 20 -> 21 -> 22 -> 20 is a three-node cycle.
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

	// WalkCycles walks over the cycles in the graph, invoking the walker
	// once per cyclic cluster with a representative cycle. For each cycle we
	// paint its vertices and the edges connecting them in red, so that the
	// cycles stand out in the rendered graph.
	cycleWalker := func(cycle []*graph.Vertex[int]) error {
		for i, v := range cycle {
			// Paint the vertex red.
			v.DotAttributes["color"] = "red"
			v.DotAttributes["fillcolor"] = "red"

			// Paint the edge from the previous vertex to this one red.
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

	// Emit the Dot representation, with the cycles painted red. Render it
	// with e.g. `go run main.go | dot -T svg -o cycles.svg`.
	if err := graph.WriteDot(g, os.Stdout); err != nil {
		fmt.Println(err)
	}
}
