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

package graph

import (
	"io"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// EchartsOpts provides options to be used when rendering the graph using
// [WriteEcharts] function.
type EchartsOpts struct {
	// SeriesName specifies the name to associate with the series data.
	SeriesName string

	// GlobalOpts specifies the global chart options
	GlobalOpts []charts.GlobalOpts

	// SeriesOpts specifies the options to apply for the series
	SeriesOpts []charts.SeriesOpts
}

// WriteEcharts renders the given [Graph] using go-echarts.
func WriteEcharts[T comparable](g Graph[T], w io.Writer, echartsOpts EchartsOpts) error {
	// Edge symbols depend on whether the graph is directed or undirected
	var edgeSymbols []string
	if g.Kind() == KindUndirected {
		edgeSymbols = []string{"none", "none"}
	} else {
		edgeSymbols = []string{"none", "arrow"}
	}

	// Create graph nodes and edges
	nodes := make([]opts.GraphNode, 0)
	links := make([]opts.GraphLink, 0)
	for _, u := range g.GetVertices() {
		// Vertices
		node := opts.GraphNode{
			Name:       u.Label,
			ItemStyle:  u.EchartsStyle,
			Symbol:     u.EchartsSymbol,
			SymbolSize: u.EchartsSymbolSize,
		}
		nodes = append(nodes, node)

		// Edges
		for _, v := range g.GetNeighbourVertices(u.Value) {
			edge := g.GetEdge(u.Value, v.Value)
			link := opts.GraphLink{
				Source:    u.Label,
				Target:    v.Label,
				LineStyle: edge.EchartsLineStyle,
				Label:     edge.EchartsEdgeLabel,
			}
			links = append(links, link)
		}
	}

	// Create graph and render it
	graph := charts.NewGraph()
	graph.AddSeries(echartsOpts.SeriesName, nodes, links)
	graph.SetGlobalOptions(echartsOpts.GlobalOpts...)

	// Default set of series opts
	seriesOpts := []charts.SeriesOpts{
		charts.WithGraphChartOpts(
			opts.GraphChart{
				Draggable:          opts.Bool(true),
				Force:              &opts.GraphForce{Repulsion: 500},
				Roam:               opts.Bool(true),
				FocusNodeAdjacency: opts.Bool(true),
				Layout:             "force",
				EdgeSymbol:         edgeSymbols,
			},
		),
		charts.WithLabelOpts(
			opts.Label{
				Show:     opts.Bool(true),
				Position: "inside",
			},
		),
	}
	seriesOpts = append(seriesOpts, echartsOpts.SeriesOpts...)
	graph.SetSeriesOptions(seriesOpts...)

	page := components.NewPage()
	page.AddCharts(graph)

	return page.Render(w)
}
