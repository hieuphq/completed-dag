package domain

import (
	"fmt"
	"strings"
)

// Vertex ...
type Vertex struct {
	Node             *Node
	ChildrenVertices Vertices
	ParentVertices   Vertices
}

// Vertices vertex list
type Vertices []Vertex

// Join 2 arrays
func (a Vertices) Join(b Vertices) Vertices {
	var rs []Vertex
	for idx := range a {
		rs = append(rs, a[idx])
	}

	for idx := range b {
		rs = append(rs, b[idx])
	}

	return rs
}

// Append a item
func (a Vertices) Append(v Vertex) Vertices {
	return append(a, v)
}

// ToString vertex string
func (v *Vertex) ToString(level int) string {
	lvStr := ""
	for idx := 0; idx < level; idx++ {
		lvStr = lvStr + "    "
	}
	rs := []string{}
	rs = append(rs, v.Node.ToString(level))

	if len(v.ChildrenVertices) > 0 {
		rs = append(rs, lvStr+fmt.Sprintf("ChildrenVerticess Vertex %v", len(v.ChildrenVertices)))
		for idx := range v.ChildrenVertices {
			rs = append(rs, v.ChildrenVertices[idx].ToString(level))
		}
	}

	if len(v.ParentVertices) > 0 {
		rs = append(rs, lvStr+"Parent Vertices")
		for idx := range v.ParentVertices {
			rs = append(rs, v.ParentVertices[idx].ToString(level))
		}
	}

	rs = append(rs, lvStr+"END vertex")
	return strings.Join(rs, "\n")
}

// ToString vertices string
func (a Vertices) ToString(level int) string {
	lvStr := ""
	for idx := 0; idx < level; idx++ {
		lvStr = lvStr + "    "
	}
	rs := []string{}
	rs = append(rs, lvStr+"===============================")
	for idx := range a {
		rs = append(rs, lvStr+"++++++++++++++++++++++++++++++")
		rs = append(rs, lvStr+a[idx].ToString(level+1))
		rs = append(rs, lvStr+"\n")
	}
	rs = append(rs, lvStr+"===============================END")
	rs = append(rs, "\n")

	return strings.Join(rs, "\n")
}
