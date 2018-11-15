package main

import (
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Node is data in db
type Node struct {
	ID      uuid.UUID
	Parents []uuid.UUID
	Flag    bool
}

// Clone node
func (n *Node) Clone() *Node {
	var ps []uuid.UUID

	for idx := range n.Parents {
		ps = append(ps, n.Parents[idx])

	}

	return &Node{
		ID:      n.ID,
		Parents: ps,
		Flag:    n.Flag,
	}
}

// AddParent to a node
func (n *Node) AddParent(ID uuid.UUID) {
	if n.Parents == nil || len(n.Parents) <= 0 {
		n.Parents = []uuid.UUID{ID}

		return
	}

	n.Parents = append(n.Parents, ID)
}

// Nodes node list
type Nodes []Node

func randomBool() bool {
	var src = rand.NewSource(time.Now().UnixNano())
	var r = rand.New(src)

	return r.Intn(2) != 0
}

func randomNumber(length int) int {
	var src = rand.NewSource(time.Now().UnixNano())
	var r = rand.New(src)

	return r.Intn(length)
}

// GenerateNodes generate note
func GenerateNodes(length int) []Node {
	rs := []Node{}

	for idx := 0; idx < length; idx++ {
		rs = append(rs, Node{
			ID:   uuid.NewV4(),
			Flag: randomBool(),
		})
	}

	return rs
}

// GenerateEdges ...
func (ns Nodes) GenerateEdges() Nodes {
	rs := []Node{}
	length := len(ns)

	for idx := range ns {
		rs = append(rs, *ns[idx].Clone())
	}

	for idx := 0; idx < length-1; idx++ {
		firstIdx := randomNumber(length)
		secondIdx := randomNumber(length)

		rs[secondIdx].AddParent(ns[firstIdx].ID)
	}
	return Nodes(rs)
}
