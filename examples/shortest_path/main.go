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
