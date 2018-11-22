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
	List(vertexID domain.UUID) []domain.Vertex
	ConditionalList(vertexID domain.UUID, flagCondition bool) []domain.Vertex
	Insert(vertex domain.Node) error
}

// NewSimpleConnectedGraph simple implement with db
func NewSimpleConnectedGraph(db store.DB, reach helper.ReachHelper, list helper.ListHelper, insert helper.InsertHelper) ConnectedGraph {
	return &simpleImpl{
		Repo:         db,
		ReachHelper:  reach,
		ListHelper:   list,
		InsertHelper: insert,
	}
}

type simpleImpl struct {
	Repo         store.DB
	ReachHelper  helper.ReachHelper
	ListHelper   helper.ListHelper
	InsertHelper helper.InsertHelper
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
	condition := helper.NoneFlagCondition
	if flagCondition {
		condition = helper.TrueFlagCondition
	}

	if !flagCondition {
		condition = helper.FalseFlagCondition
	}

	rs, err := s.ReachHelper.GetReachCondition(vertexID, condition)
	if err != nil {
		pp.Print(err)
		return 0
	}

	return rs
}

func (s *simpleImpl) List(vertexID domain.UUID) []domain.Vertex {
	rs, err := s.ListHelper.GetReachList(vertexID)
	if err != nil {
		pp.Print(err)
		return nil
	}

	return rs
}

func (s *simpleImpl) ConditionalList(vertexID domain.UUID, flagCondition bool) []domain.Vertex {
	condition := helper.NoneFlagCondition
	if flagCondition {
		condition = helper.TrueFlagCondition
	}

	if !flagCondition {
		condition = helper.FalseFlagCondition
	}

	rs, err := s.ListHelper.GetReachListCondition(vertexID, condition)
	if err != nil {
		pp.Print(err)
		return nil
	}

	return rs
}

func (s *simpleImpl) Insert(vertex domain.Node) error {
	return s.InsertHelper.Insert(&vertex)
}
