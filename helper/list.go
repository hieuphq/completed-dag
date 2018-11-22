package helper

import (
	"github.com/hieuphq/completed-dag/domain"
	"github.com/hieuphq/completed-dag/errors"
	"github.com/hieuphq/completed-dag/store"
)

// ListHelper reach helper
type ListHelper interface {
	GetReachList(ID domain.UUID) ([]domain.Vertex, error)
	GetReachListCondition(ID domain.UUID, condition FlagCondition) ([]domain.Vertex, error)
}

// NewSimpleListHelper ...
func NewSimpleListHelper(db store.DB) ListHelper {
	return &simpleListImpl{
		db: db,
	}
}

type simpleListImpl struct {
	db store.DB
}

func (s *simpleListImpl) getNode(ID domain.UUID) (*domain.Node, error) {
	dt, err := s.db.Get(ID.ToBytes())
	if err != nil {
		return nil, err
	}

	return domain.NewNodeFromBytes(dt)
}

func (s *simpleListImpl) putNode(node *domain.Node) error {
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

func (s *simpleListImpl) GetReachList(ID domain.UUID) ([]domain.Vertex, error) {
	child, err := s.getReachChildrenList(ID, NoneFlagCondition, 0)
	if err != nil {
		return nil, err
	}

	parent, err := s.getReachParentsList(ID, NoneFlagCondition, 0)
	if err != nil {
		return nil, err
	}

	return domain.Vertices([]domain.Vertex{*combineFinalData(child, parent)}), nil
}

func (s *simpleListImpl) GetReachListCondition(ID domain.UUID, condition FlagCondition) ([]domain.Vertex, error) {
	child, err := s.getReachChildrenList(ID, condition, 0)
	if err != nil {
		return nil, err
	}

	parent, err := s.getReachParentsList(ID, condition, 0)
	if err != nil {
		return nil, err
	}

	return domain.Vertices([]domain.Vertex{*combineFinalData(child, parent)}), nil
}

func (s *simpleListImpl) getReachChildrenList(ID domain.UUID, condition FlagCondition, level int) ([]domain.Vertex, error) {
	nd, err := s.getNode(ID)

	if err != nil {
		return nil, err
	}

	child := domain.Vertices{}
	currVertex := domain.Vertex{
		Node: nd,
	}

	if len(nd.Children) <= 0 {

		if level == 0 {
			return []domain.Vertex{}, nil
		}

		switch condition {
		case NoneFlagCondition:
			child = child.Append(currVertex)

		case TrueFlagCondition:
			if nd.Flag {
				child = child.Append(currVertex)
			}

		case FalseFlagCondition:
			if !nd.Flag {
				child = child.Append(currVertex)
			}
		}
		return child, nil
	}

	childVertices := domain.Vertices{}
	for idx := range nd.Children {
		cID := nd.Children[idx]

		currChildren, err := s.getReachChildrenList(cID, condition, level+1)
		if err != nil {
			return child, err
		}
		childVertices = childVertices.Join(domain.Vertices(currChildren))
	}

	switch condition {
	case NoneFlagCondition:
		currVertex.ChildrenVertices = currVertex.ChildrenVertices.Join(childVertices)

	case TrueFlagCondition:
		if nd.Flag {
			currVertex.ChildrenVertices = currVertex.ChildrenVertices.Join(childVertices)
		}

	case FalseFlagCondition:
		if !nd.Flag {
			currVertex.ChildrenVertices = currVertex.ChildrenVertices.Join(childVertices)
		}
	}
	return []domain.Vertex{currVertex}, nil
}

func (s *simpleListImpl) getReachParentsList(ID domain.UUID, condition FlagCondition, level int) ([]domain.Vertex, error) {
	nd, err := s.getNode(ID)

	if err != nil {
		return nil, err
	}

	parent := domain.Vertices{}
	currVertex := domain.Vertex{
		Node: nd,
	}

	if len(nd.Parents) <= 0 {

		if level == 0 {
			return []domain.Vertex{}, nil
		}

		switch condition {
		case NoneFlagCondition:
			parent = parent.Append(currVertex)

		case TrueFlagCondition:
			if nd.Flag {
				parent = parent.Append(currVertex)
			}

		case FalseFlagCondition:
			if !nd.Flag {
				parent = parent.Append(currVertex)
			}
		}
		return parent, nil
	}

	parentVertices := domain.Vertices{}
	for idx := range nd.Parents {
		cID := nd.Parents[idx]

		currParent, err := s.getReachParentsList(cID, condition, level+1)
		if err != nil {
			return parent, err
		}
		// currItm.ParentVertices = currParent
		parentVertices = parentVertices.Join(domain.Vertices(currParent))
	}

	switch condition {
	case NoneFlagCondition:
		currVertex.ParentVertices = currVertex.ParentVertices.Join(parentVertices)

	case TrueFlagCondition:
		if nd.Flag {
			currVertex.ParentVertices = currVertex.ParentVertices.Join(parentVertices)
		}

	case FalseFlagCondition:
		if !nd.Flag {
			currVertex.ParentVertices = currVertex.ParentVertices.Join(parentVertices)
		}
	}
	return []domain.Vertex{currVertex}, nil
}
