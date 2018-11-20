package graph

import (
	"github.com/hieuphq/completed-dag/domain"
	"github.com/hieuphq/completed-dag/helper"
	"github.com/hieuphq/completed-dag/store"
	"github.com/k0kubun/pp"
)

// ConnectedGraph DAG graph with connected implement
type ConnectedGraph interface {
	Reach(vertexID domain.UUID) int
	ConditionalReach(vertexID domain.UUID, flagCondition bool) int
	List(vertexID domain.UUID) []domain.Node
	ConditionalList(vertexID domain.UUID, flagCondition bool) []domain.Node
	Insert(vertex domain.Node) error
}

// NewSimpleConnectedGraph simple implement with db
func NewSimpleConnectedGraph(db store.DB, reach helper.ReachHelper) ConnectedGraph {
	return &simpleImpl{
		Repo:        db,
		ReachHelper: reach,
	}
}

type simpleImpl struct {
	Repo        store.DB
	ReachHelper helper.ReachHelper
}

func (s *simpleImpl) Reach(vertexID domain.UUID) int {
	rs, err := s.ReachHelper.GetReach(vertexID)
	if err != nil {
		pp.Print(err)
		return 0
	}

	return rs
}

func (s *simpleImpl) ConditionalReach(vertexID domain.UUID, flagCondition bool) int {
	return 0
}

func (s *simpleImpl) List(vertexID domain.UUID) []domain.Node {
	return nil
}

func (s *simpleImpl) ConditionalList(vertexID domain.UUID, flagCondition bool) []domain.Node {
	return nil
}

func (s *simpleImpl) Insert(vertex domain.Node) error {
	return nil
}
