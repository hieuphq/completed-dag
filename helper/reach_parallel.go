package helper

import (
	"runtime"
	"sync"

	"github.com/hieuphq/completed-dag/domain"
	"github.com/hieuphq/completed-dag/errors"
	"github.com/hieuphq/completed-dag/store"
)

type parallelReachImpl struct {
	db store.DB
}

// NewParallelReachHelper ...
func NewParallelReachHelper(db store.DB) ReachHelper {
	return &parallelReachImpl{
		db: db,
	}
}

func (s *parallelReachImpl) getNode(ID domain.UUID) (*domain.Node, error) {
	dt, err := s.db.Get(ID.ToBytes())
	if err != nil {
		return nil, err
	}

	return domain.NewNodeFromBytes(dt)
}

func (s *parallelReachImpl) putNode(node *domain.Node) error {
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

func (s *parallelReachImpl) GetReach(ID domain.UUID) (int, error) {
	var parentWG sync.WaitGroup
	childCount := 0
	parentCount := 0
	numCPU := runtime.NumCPU() / 2
	parentWG.Add(2)
	semaphoreChan := make(chan struct{}, numCPU)

	go func(cCount *int, pWG *sync.WaitGroup, sem chan struct{}) {
		var wg sync.WaitGroup
		childCount := make(chan int)
		doneCh := make(chan bool)

		result := 0
		go func(rs *int, cCh <-chan int, dCh <-chan bool) {
			for {
				select {
				case cVal := <-cCh:
					*rs = *rs + cVal

				case dVal := <-dCh:
					if dVal {
						break
					}
				}
			}

		}(&result, childCount, doneCh)

		wg.Add(1)
		sem <- struct{}{}
		s.getReachChildrenCount(ID, NoneFlagCondition, 0, childCount, &wg, sem)

		wg.Wait()

		doneCh <- true
		*cCount = result
		pWG.Done()

	}(&childCount, &parentWG, semaphoreChan)

	go func(pCount *int, pWG *sync.WaitGroup, sem chan struct{}) {
		var wg sync.WaitGroup
		parentCount := make(chan int)
		doneCh := make(chan bool)

		result := 0
		go func(rs *int, cCh <-chan int, dCh <-chan bool) {
			for {
				select {
				case cVal := <-cCh:
					*rs = *rs + cVal

				case dVal := <-dCh:
					if dVal {
						break
					}
				}
			}

		}(&result, parentCount, doneCh)

		wg.Add(1)
		sem <- struct{}{}
		s.getReachParentsCount(ID, NoneFlagCondition, 0, parentCount, &wg, sem)

		wg.Wait()

		doneCh <- true
		*pCount = result
		pWG.Done()

	}(&parentCount, &parentWG, semaphoreChan)

	parentWG.Wait()

	return childCount + parentCount, nil
}

func (s *parallelReachImpl) GetReachCondition(ID domain.UUID, condition FlagCondition) (int, error) {
	return 0, nil
}

func (s *parallelReachImpl) getReachParentsCount(ID domain.UUID, condition FlagCondition, level int, parentCh chan<- int, wg *sync.WaitGroup, sem chan struct{}) {
	defer func(w *sync.WaitGroup, cSem chan struct{}) {
		w.Done()
		<-cSem
	}(wg, sem)
	nd, err := s.getNode(ID)
	if err != nil {
		return
	}

	if len(nd.Parents) <= 0 {

		if level == 0 {
			parentCh <- 0
			return
		}

		switch condition {
		case NoneFlagCondition:
			parentCh <- 1

		case TrueFlagCondition:
			if nd.Flag {
				parentCh <- 1
			}

		case FalseFlagCondition:
			if !nd.Flag {
				parentCh <- 1
			}
		}

		return
	}

	for idx := range nd.Parents {
		cID := nd.Parents[idx]

		wg.Add(1)
		go func(cS *parallelReachImpl, cCID domain.UUID, cCondition FlagCondition, cLevel int, cParentCh chan<- int, cWg *sync.WaitGroup, cSem chan struct{}) {
			cSem <- struct{}{}
			cS.getReachParentsCount(cCID, cCondition, cLevel, cParentCh, cWg, cSem)
		}(s, cID, condition, level+1, parentCh, wg, sem)
	}

	if level > 0 {

		switch condition {
		case NoneFlagCondition:
			parentCh <- 1

		case TrueFlagCondition:
			if nd.Flag {
				parentCh <- 1
			}

		case FalseFlagCondition:
			if !nd.Flag {
				parentCh <- 1
			}
		}
	}

	return
}

func (s *parallelReachImpl) getReachChildrenCount(ID domain.UUID, condition FlagCondition, level int, childCh chan<- int, wg *sync.WaitGroup, sem chan struct{}) {
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
			childCh <- 0
			return
		}

		switch condition {
		case NoneFlagCondition:
			childCh <- 1

		case TrueFlagCondition:
			if nd.Flag {
				childCh <- 1
			}

		case FalseFlagCondition:
			if !nd.Flag {
				childCh <- 1
			}
		}

		return
	}

	for idx := range nd.Children {
		cID := nd.Children[idx]

		wg.Add(1)
		go func(cS *parallelReachImpl, cCID domain.UUID, cCondition FlagCondition, cLevel int, cChildCh chan<- int, cWg *sync.WaitGroup, cSem chan struct{}) {
			cSem <- struct{}{}
			cS.getReachChildrenCount(cCID, cCondition, cLevel, cChildCh, cWg, cSem)
		}(s, cID, condition, level+1, childCh, wg, sem)
	}

	if level > 0 {

		switch condition {
		case NoneFlagCondition:
			childCh <- 1

		case TrueFlagCondition:
			if nd.Flag {
				childCh <- 1
			}

		case FalseFlagCondition:
			if !nd.Flag {
				childCh <- 1
			}
		}
	}

	return
}
