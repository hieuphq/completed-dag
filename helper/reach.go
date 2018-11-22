package helper

import (
	"github.com/hieuphq/completed-dag/domain"
	"github.com/hieuphq/completed-dag/errors"
	"github.com/hieuphq/completed-dag/store"
)

// FlagCondition condition to filter
type FlagCondition uint8

const (

	// NoneFlagCondition flag is none
	NoneFlagCondition FlagCondition = iota

	// TrueFlagCondition flag is true
	TrueFlagCondition

	// FalseFlagCondition flag is false
	FalseFlagCondition
)

// ReachHelper reach helper
type ReachHelper interface {
	GetReach(ID domain.UUID) (int, error)
	GetReachCondition(ID domain.UUID, condition FlagCondition) (int, error)
}

// NewSimpleReachHelper ...
func NewSimpleReachHelper(db store.DB) ReachHelper {
	return &simpleReachImpl{
		db: db,
	}
}

type simpleReachImpl struct {
	db store.DB
}

func (s *simpleReachImpl) getNode(ID domain.UUID) (*domain.Node, error) {
	dt, err := s.db.Get(ID.ToBytes())
	if err != nil {
		return nil, err
	}

	return domain.NewNodeFromBytes(dt)
}

func (s *simpleReachImpl) putNode(node *domain.Node) error {
	uuid := node.ID

	if uuid.IsZero() {
		return errors.ErrInvalidID
	}

	if node == nil {
		return errors.ErrNil
	}

	nodeDt, err := node.ToBytes()
	if err != nil {
		return err
	}

	return s.db.Put(uuid.ToBytes(), nodeDt)
}

func (s *simpleReachImpl) GetReach(ID domain.UUID) (int, error) {
	childCount, err := s.getReachChildrenCount(ID, NoneFlagCondition, 0)
	if err != nil {
		return 0, err
	}

	parentCount, err := s.getReachParentsCount(ID, NoneFlagCondition, 0)
	if err != nil {
		return 0, err
	}

	return childCount + parentCount, nil
}

func (s *simpleReachImpl) GetReachCondition(ID domain.UUID, condition FlagCondition) (int, error) {
	childCount, err := s.getReachChildrenCount(ID, condition, 0)
	if err != nil {
		return 0, err
	}

	parentCount, err := s.getReachParentsCount(ID, condition, 0)
	if err != nil {
		return 0, err
	}

	return childCount + parentCount, nil
}

func (s *simpleReachImpl) getReachChildrenCount(ID domain.UUID, condition FlagCondition, level int) (int, error) {
	nd, err := s.getNode(ID)
	if err != nil {
		return 0, err
	}

	sum := 0

	if len(nd.Children) <= 0 {

		if level == 0 {
			return 0, nil
		}

		switch condition {
		case NoneFlagCondition:
			sum = sum + 1

		case TrueFlagCondition:
			if nd.Flag {
				sum = sum + 1
			}

		case FalseFlagCondition:
			if !nd.Flag {
				sum = sum + 1
			}
		}

		return sum, nil
	}

	for idx := range nd.Children {
		cID := nd.Children[idx]

		count, err := s.getReachChildrenCount(cID, condition, level+1)
		if err != nil {
			return sum, err
		}

		sum = sum + count
	}

	if level > 0 {

		switch condition {
		case NoneFlagCondition:
			sum = sum + 1

		case TrueFlagCondition:
			if nd.Flag {
				sum = sum + 1
			}

		case FalseFlagCondition:
			if !nd.Flag {
				sum = sum + 1
			}
		}
	}
	return sum, nil
}

func (s *simpleReachImpl) getReachParentsCount(ID domain.UUID, condition FlagCondition, level int) (int, error) {
	nd, err := s.getNode(ID)
	if err != nil {
		return 0, err
	}

	sum := 0

	if len(nd.Parents) <= 0 {

		if level == 0 {
			return 0, nil
		}

		switch condition {
		case NoneFlagCondition:
			sum = sum + 1

		case TrueFlagCondition:
			if nd.Flag {
				sum = sum + 1
			}

		case FalseFlagCondition:
			if !nd.Flag {
				sum = sum + 1
			}
		}

		return sum, nil
	}

	for idx := range nd.Parents {
		cID := nd.Parents[idx]

		count, err := s.getReachParentsCount(cID, condition, level+1)
		if err != nil {
			return sum, err
		}

		sum = sum + count
	}

	if level > 0 {

		switch condition {
		case NoneFlagCondition:
			sum = sum + 1

		case TrueFlagCondition:
			if nd.Flag {
				sum = sum + 1
			}

		case FalseFlagCondition:
			if !nd.Flag {
				sum = sum + 1
			}
		}
	}
	return sum, nil
}

func combineFinalData(child []domain.Vertex, parent []domain.Vertex) *domain.Vertex {
	if len(child) <= 0 && len(parent) <= 0 {
		return nil
	}

	var currNode *domain.Node
	var childVertices []domain.Vertex
	var parentVertices []domain.Vertex
	if len(child) > 0 {
		currNode = child[0].Node.Clone()
		childVertices = child[0].ChildrenVertices
	}

	if len(parent) > 0 {
		currNode = parent[0].Node.Clone()
		parentVertices = parent[0].ParentVertices
	}

	return &domain.Vertex{
		Node:             currNode,
		ChildrenVertices: childVertices,
		ParentVertices:   parentVertices,
	}

}
