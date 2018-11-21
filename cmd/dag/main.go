package main

import (
	"github.com/hieuphq/completed-dag/domain"
	"github.com/hieuphq/completed-dag/generator"
	"github.com/hieuphq/completed-dag/graph"
	"github.com/hieuphq/completed-dag/helper"
	"github.com/hieuphq/completed-dag/store"
	"github.com/k0kubun/pp"
)

func main() {
	gen := generator.NewSimpleGenerator()

	db := store.NewMemory()

	defer db.Close()

	size := 4
	ns := gen.Generate(size, size-1)
	err := saveToDB(ns, db)
	if err != nil {
		pp.Println(err)
	}

	g := graph.NewSimpleConnectedGraph(db, helper.NewSimpleReachHelper(db), helper.NewSimpleInsertHelper(db))

	domain.Nodes(ns).Print()

	// for idx := range ns {
	// curr := ns[1]
	// rs := g.List(curr.ID)
	// pp.Println("Finished")
	// pp.Println(rs)
	// fmt.Println(domain.Vertices(rs).ToString(0))
	// pp.Println(g.Reach(curr.ID))
	// }

	pp.Println(ns)

	nNode := domain.Node{
		ID:      domain.NewUUID(),
		Parents: []domain.UUID{ns[0].ID},
	}
	pp.Println(nNode)
	g.Insert(nNode)

	currAll, err := getFromDB(db)
	if err != nil {
		pp.Println(err)
	}
	pp.Println(currAll)
}

func getFromDB(db store.DB) ([]domain.Node, error) {
	currAll, err := db.All()
	if err != nil {
		pp.Println(err)
	}

	rs := []domain.Node{}
	for key := range currAll {
		n, err := domain.NewNodeFromBytes(currAll[key])
		if err != nil {
			return nil, err
		}
		rs = append(rs, *n)
	}

	return rs, nil
}

func saveToDB(ns []domain.Node, db store.DB) error {
	for idx := range ns {
		curr := ns[idx]
		currBytes, err := curr.ToBytes()
		if err != nil {
			pp.Println("can't parse", err)
			return err
		}
		if err := db.Put(curr.ID.ToBytes(), currBytes); err != nil {
			pp.Println("can't add", err)
			return err
		}
	}

	return nil
}
