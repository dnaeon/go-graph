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
	"fmt"
	"io"
	"strconv"
	"strings"
)

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
func WriteDot[T comparable](g Graph[T], w io.Writer) error {
	var graphKind string
	var edgeArrow string
	if g.Kind() == KindUndirected {
		graphKind = "graph"
		edgeArrow = "--"
	} else {
		graphKind = "digraph"
		edgeArrow = "->"
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

	for _, v := range g.GetVertices() {
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
			e := g.GetEdge(v.Value, u.Value)
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
