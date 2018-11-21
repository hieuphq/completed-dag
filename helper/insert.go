package helper

import (
	"github.com/hieuphq/completed-dag/domain"
	"github.com/hieuphq/completed-dag/errors"
	"github.com/hieuphq/completed-dag/store"
	"github.com/k0kubun/pp"
)

// InsertHelper reach helper
type InsertHelper interface {
	Insert(n *domain.Node) error
}

// NewSimpleInsertHelper ...
func NewSimpleInsertHelper(db store.DB) InsertHelper {
	return &simpleInsertImpl{
		db: db,
	}
}

type simpleInsertImpl struct {
	db store.DB
}

func getNodeFromDB(ID domain.UUID, db store.DB) (*domain.Node, error) {
	dt, err := db.Get(ID.ToBytes())
	if err != nil {
		return nil, err
	}

	return domain.NewNodeFromBytes(dt)
}
func putNodeFromDB(node *domain.Node, db store.DB) error {
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

	return db.Put(uuid.ToBytes(), nodeDt)
}

func (s *simpleInsertImpl) getNode(ID domain.UUID) (*domain.Node, error) {
	return getNodeFromDB(ID, s.db)
}

func (s *simpleInsertImpl) putNode(node *domain.Node) error {
	return putNodeFromDB(node, s.db)
}

func (s *simpleInsertImpl) Insert(n *domain.Node) error {

	currNode, err := s.getNode(n.ID)
	if err != nil && err != errors.ErrNil {
		return err
	}

	virtualDB := store.NewMemory()

	// isNew
	if err != nil && err == errors.ErrNil {

		if len(n.Parents) > 0 {
			parentMap := map[domain.UUID]domain.Node{}
			for idx := range n.Parents {
				pID := n.Parents[idx]
				pNode, err := s.getNode(pID)
				if err != nil {
					pp.Print(err)
					continue
				}

				// if contains parent -> remove a connection
				if pNode.ContainsParentID(n.ID) {
					if pNode.Greater(n) {
						// parent > currN --> remove parentID in parent
						pNode.RemoveParent(n.ID)

					} else {
						// parent < currN --> remove parentID in currN
						n.RemoveParent(pNode.ID)
					}
				}

				pNode.AddChild(n.ID)
				parentMap[pID] = *pNode
			}

			// save into virtual db
			putNodeFromDB(n, virtualDB)
			for key := range parentMap {
				curV := parentMap[key]
				putNodeFromDB(&curV, virtualDB)
			}

		}

	}

	// not NEW --> NO NEED
	if err == nil && currNode != nil {
	}

	virtualDB.Transfer(s.db)

	return nil
}
