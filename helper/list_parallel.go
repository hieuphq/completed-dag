package helper

import (
	"runtime"
	"sync"

	"github.com/k0kubun/pp"

	"github.com/hieuphq/completed-dag/domain"
	"github.com/hieuphq/completed-dag/errors"
	"github.com/hieuphq/completed-dag/store"
)

type parallelListImpl struct {
	db store.DB
}

// NewParallelListHelper ...
func NewParallelListHelper(db store.DB) ListHelper {
	return &parallelListImpl{
		db: db,
	}
}

func (s *parallelListImpl) getNode(ID domain.UUID) (*domain.Node, error) {
	dt, err := s.db.Get(ID.ToBytes())
	if err != nil {
		return nil, err
	}

	return domain.NewNodeFromBytes(dt)
}

func (s *parallelListImpl) putNode(node *domain.Node) error {
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

func (s *parallelListImpl) GetReachList(ID domain.UUID) ([]domain.Vertex, error) {
	return s.getReachWithCondition(ID, NoneFlagCondition)
}

func (s *parallelListImpl) GetReachListCondition(ID domain.UUID, condition FlagCondition) ([]domain.Vertex, error) {
	return s.getReachWithCondition(ID, condition)
}

func (s *parallelListImpl) getReachWithCondition(ID domain.UUID, condition FlagCondition) ([]domain.Vertex, error) {
	var parentWG sync.WaitGroup
	childList := []domain.Vertex{}
	parentList := []domain.Vertex{}
	numCPU := runtime.NumCPU()
	parentWG.Add(2)
	semaphoreChan := make(chan struct{}, numCPU)

	go func(children []domain.Vertex, pWG *sync.WaitGroup, sem chan struct{}) {
		var wg sync.WaitGroup
		childItm := make(chan vertexWithParentID)
		doneCh := make(chan bool)
		mapChildren := make(map[domain.UUID][]domain.Vertex)

		go func(iChildren []domain.Vertex, cCh <-chan vertexWithParentID, dCh <-chan bool) {
			childrenIDs := []domain.UUID{}
			mapParent := map[keyPC]struct{}{}
			for {
				select {
				case cVal := <-cCh:
					for idx := range cVal.Node.Parents {
						pID := cVal.Node.Parents[idx]

						currV := domain.Vertex{
							Node: &cVal.Node,
						}
						pParent, ok1 := mapChildren[cVal.Node.ID]

						if ok1 {
							currV.ChildrenVertices = domain.Vertices(pParent).Clone()
						}

						pChildren, ok := mapChildren[pID]
						if !ok {
							mapChildren[pID] = []domain.Vertex{currV}
						} else {
							mapChildren[pID] = append(pChildren, currV)
						}

						mapParent[keyPC{Parent: pID, Child: cVal.Node.ID}] = struct{}{}

						if pID == ID {
							childrenIDs = append(childrenIDs, cVal.Node.ID)
						}
					}

				case dVal := <-dCh:
					if dVal {
						// fmt.Println(mapChildren)
						// pp.Println("DONE")
						// pp.Println(childrenIDs)
						// pp.Println("DONE")
						// pp.Println("DONE")
						for idx := range childrenIDs {
							pp.Println(mapChildren[childrenIDs[idx]])
							iChildren = append(iChildren, mapChildren[childrenIDs[idx]]...)
						}

						pp.Println("iChild", iChildren)
						break
					}
				}
			}

		}(children, childItm, doneCh)

		wg.Add(1)
		sem <- struct{}{}
		s.getReachChildrenList(ID, condition, 0, childItm, &wg, sem)

		wg.Wait()

		doneCh <- true
		pWG.Done()

	}(childList, &parentWG, semaphoreChan)

	go func(iParent []domain.Vertex, pWG *sync.WaitGroup, sem chan struct{}) {
		// var wg sync.WaitGroup
		// parentCount := make(chan int)
		// doneCh := make(chan bool)

		// result := 0
		// go func(rs *int, cCh <-chan int, dCh <-chan bool) {
		// 	for {
		// 		select {
		// 		case cVal := <-cCh:
		// 			*rs = *rs + cVal

		// 		case dVal := <-dCh:
		// 			if dVal {
		// 				break
		// 			}
		// 		}
		// 	}

		// }(&result, parentCount, doneCh)

		// wg.Add(1)
		// sem <- struct{}{}
		// s.getReachParentsCount(ID, condition, 0, parentCount, &wg, sem)

		// wg.Wait()

		// doneCh <- true
		// *pCount = result
		pWG.Done()

	}(parentList, &parentWG, semaphoreChan)

	parentWG.Wait()

	pp.Println(childList)
	pp.Println(parentList)

	return domain.Vertices([]domain.Vertex{*combineFinalData(childList, parentList)}), nil
}

type vertexWithParentID struct {
	Node domain.Node
}

type keyPC struct {
	Parent domain.UUID
	Child  domain.UUID
}

func (s *parallelListImpl) getReachChildrenList(ID domain.UUID, condition FlagCondition, level int, childCh chan<- vertexWithParentID, wg *sync.WaitGroup, sem chan struct{}) {
	defer func(w *sync.WaitGroup, cSem chan struct{}) {
		w.Done()
		<-cSem
	}(wg, sem)

	nd, err := s.getNode(ID)
	if err != nil {
		return
	}

	if len(nd.Children) <= 0 {

		if level == 0 {
			return
		}

		switch condition {
		case NoneFlagCondition:
			childCh <- vertexWithParentID{
				Node: *nd,
			}

		case TrueFlagCondition:
			if nd.Flag {
				childCh <- vertexWithParentID{
					Node: *nd,
				}
			}

		case FalseFlagCondition:
			if !nd.Flag {
				childCh <- vertexWithParentID{
					Node: *nd,
				}
			}
		}

		return
	}

	for idx := range nd.Children {
		cID := nd.Children[idx]

		wg.Add(1)
		go func(cS *parallelListImpl, cCID domain.UUID, cCondition FlagCondition, cLevel int, cChildCh chan<- vertexWithParentID, cWg *sync.WaitGroup, cSem chan struct{}) {
			cSem <- struct{}{}
			cS.getReachChildrenList(cCID, cCondition, cLevel, cChildCh, cWg, cSem)
		}(s, cID, condition, level+1, childCh, wg, sem)
	}

	if level > 0 {

		switch condition {
		case NoneFlagCondition:
			childCh <- vertexWithParentID{
				Node: *nd,
			}

		case TrueFlagCondition:
			if nd.Flag {
				childCh <- vertexWithParentID{
					Node: *nd,
				}
			}

		case FalseFlagCondition:
			if !nd.Flag {
				childCh <- vertexWithParentID{
					Node: *nd,
				}
			}
		}
	}

	return
}
