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

package graph_test

import (
	"bytes"
	"strings"
	"testing"

	"gopkg.in/dnaeon/go-graph.v1"
)

func TestWriteDot(t *testing.T) {
	// Undirected graph
	g1 := graph.New[int](graph.KindUndirected)
	g1.AddEdge(1, 2)
	g1.AddEdge(1, 3)
	g1.AddEdge(3, 4)

	var buf1 bytes.Buffer
	if err := graph.WriteDot(g1, &buf1); err != nil {
		t.Fatal(err)
	}
	dotOutput1 := buf1.String()

	if !strings.Contains(dotOutput1, "strict graph") {
		t.Fatal("expected strict graph in Dot representation")
	}

	// Directed graph
	g2 := graph.New[int](graph.KindDirected)
	g2.AddEdge(1, 2)
	g2.AddEdge(1, 3)
	g2.AddEdge(3, 4)

	var buf2 bytes.Buffer
	if err := graph.WriteDot(g2, &buf2); err != nil {
		t.Fatal(err)
	}
	dotOutput2 := buf2.String()

	if !strings.Contains(dotOutput2, "strict digraph") {
		t.Fatal("expected strict digraph in Dot representation")
	}
}
