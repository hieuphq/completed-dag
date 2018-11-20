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

	size := 10
	ns := gen.Generate(size, size-1)
	err := saveToDB(ns, db)
	if err != nil {
		pp.Println(err)
	}

	g := graph.NewSimpleConnectedGraph(db, helper.NewSimpleReachHelper(db))

	domain.Nodes(ns).Print()

	for idx := range ns {
		curr := ns[idx]
		pp.Println(g.Reach(curr.ID))
	}

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
