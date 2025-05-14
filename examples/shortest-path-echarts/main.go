// Copyright (c) 2025 Marin Atanasov Nikolov <dnaeon@gmail.com>
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

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

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
		// Paint vertices in green, which form the shortest path
		v.EchartsStyle = &opts.ItemStyle{
			Color: "green",
		}

		if prev != nil {
			edge := g.GetEdge(prev.Value, v.Value)
			edge.EchartsLineStyle = &opts.LineStyle{
				Width: 5.0,
				Color: "pink",
			}
		}
		prev = v

		return nil
	}

	if err := graph.WalkShortestPath(g, 1, 8, walker); err != nil {
		fmt.Println(err)
	}

	echartsOpts := graph.EchartsOpts{
		GlobalOpts: []charts.GlobalOpts{
			charts.WithInitializationOpts(opts.Initialization{Width: "100%", Height: "95vh"}),
			charts.WithTitleOpts(opts.Title{Title: "basic shortest path example"}),
			charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true)}),
		},
		SeriesName: "shortest path",
	}

	if err := graph.WriteEcharts(g, os.Stdout, echartsOpts); err != nil {
		fmt.Println(err)
	}
}
