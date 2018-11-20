package generator

import (
	"github.com/hieuphq/completed-dag/domain"
	"github.com/hieuphq/completed-dag/util"
)

// Generator a graph
type Generator interface {
	Generate(sizeVertices int, sizeEgdes int) []domain.Node
}

// NewSimpleGenerator simple generator
func NewSimpleGenerator() Generator {
	return &simpleImpl{}
}

type simpleImpl struct {
}

func (s *simpleImpl) Generate(sizeVertices int, sizeEgdes int) []domain.Node {
	ns := []domain.Node{}
	type key struct{ a, b domain.UUID }
	indexMap := map[int]domain.UUID{}
	vertexMap := map[domain.UUID]int{}
	parentMap := map[domain.UUID][]domain.UUID{}
	childrenMap := map[domain.UUID][]domain.UUID{}
	edges := make(map[key]struct{})

	nEgdes := sizeEgdes
	for idx := 0; idx < sizeVertices; idx++ {
		nID := domain.NewUUID()

		currNode := domain.Node{
			ID:   nID,
			Flag: util.RandomBool(),
		}
		indexMap[idx] = nID
		vertexMap[nID] = idx

		var nodeParents []domain.UUID
		if idx != 0 {
			parentIdx := util.RandomNumber(idx)
			nodeParents = append(nodeParents, indexMap[parentIdx])
			nEgdes = nEgdes - 1
			parentMap[nID] = nodeParents
			childrenMap[indexMap[parentIdx]] = append(childrenMap[indexMap[parentIdx]], nID)
			edges[key{indexMap[parentIdx], nID}] = struct{}{}
		}

		ns = append(ns, currNode)
	}

	for nEgdes > 0 {
		childIdx := util.RandomNumber(sizeVertices)
		child := indexMap[childIdx]
		parent := indexMap[util.RandomNumber(childIdx)]

		if _, exists := edges[key{parent, child}]; !exists {
			edges[key{parent, child}] = struct{}{}

			parentMap[child] = append(parentMap[child], parent)
			childrenMap[parent] = append(childrenMap[parent], child)
			nEgdes = nEgdes - 1
		}
	}

	rs := []domain.Node{}

	for idx := range ns {
		curr := ns[idx]

		itm := curr.Clone()
		itm.Children = childrenMap[curr.ID]
		itm.Parents = parentMap[curr.ID]

		rs = append(rs, *itm)
	}

	return rs
}
