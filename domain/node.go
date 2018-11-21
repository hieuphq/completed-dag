package domain

import (
	"bytes"
	"encoding/gob"

	"github.com/k0kubun/pp"
)

// Node is data in db
type Node struct {
	ID       UUID
	Parents  []UUID
	Children []UUID
	Flag     bool
}

// Clone node
func (n *Node) Clone() *Node {
	var ps []UUID

	for idx := range n.Parents {
		ps = append(ps, n.Parents[idx])

	}

	var cs []UUID
	for idx := range n.Children {
		cs = append(cs, n.Children[idx])

	}

	return &Node{
		ID:       n.ID,
		Parents:  ps,
		Children: cs,
		Flag:     n.Flag,
	}
}

// AddParent to a node
func (n *Node) AddParent(ID UUID) {
	if n.Parents == nil || len(n.Parents) <= 0 {
		n.Parents = []UUID{ID}

		return
	}

	n.Parents = append(n.Parents, ID)
}

// ToBytes convert to bytes
func (n *Node) ToBytes() ([]byte, error) {
	encBuf := new(bytes.Buffer)
	err := gob.NewEncoder(encBuf).Encode(*n)
	if err != nil {
		return nil, err
	}

	return encBuf.Bytes(), nil
}

// NewNodeFromBytes create node from bytes
func NewNodeFromBytes(b []byte) (*Node, error) {
	decBuf := bytes.NewBuffer(b)
	dt := Node{}
	err := gob.NewDecoder(decBuf).Decode(&dt)

	if err != nil {
		return nil, err
	}

	return &dt, nil
}

// Nodes node list
type Nodes []Node

// Print string
func (ns Nodes) Print() {
	for idx := range ns {
		curr := ns[idx]

		pp.Println("ID: " + curr.ID.String())
		pp.Println("Parents: ")
		for jdx := range curr.Parents {
			pp.Println("  " + curr.Parents[jdx].String())
		}
		pp.Println("Children: ")
		for jdx := range curr.Children {
			pp.Println("  " + curr.Children[jdx].String())
		}
		pp.Println("Flag: ", curr.Flag)
		pp.Println()
	}
}
