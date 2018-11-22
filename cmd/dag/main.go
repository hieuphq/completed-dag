package main

import (
	"fmt"
	"time"

	"github.com/hieuphq/completed-dag/domain"
	"github.com/hieuphq/completed-dag/generator"
	"github.com/hieuphq/completed-dag/graph"
	"github.com/hieuphq/completed-dag/helper"
	"github.com/hieuphq/completed-dag/store"
	"github.com/hieuphq/completed-dag/util"
	"github.com/k0kubun/pp"
)

func generateData(db store.DB) error {
	size := 100000
	gen := generator.NewSimpleGenerator()
	ns := gen.Generate(size, size-1)
	err := saveToDB(ns, db)
	if err != nil {
		return err
	}
	fmt.Println(ns[0].ToString(0))

	return nil
}
func main() {

	// defer profile.Start(profile.CPUProfile).Stop()

	db, err := store.NewDB("./temp.leveldb")
	if err != nil {
		return
	}

	defer db.Close()

	// generateData(db)

	g := graph.NewSimpleConnectedGraph(
		db,
		helper.NewParallelReachHelper(db),
		helper.NewSimpleListHelper(db),
		helper.NewSimpleInsertHelper(db),
	)

	// g = graph.NewSimpleConnectedGraph(
	// 	db,
	// 	helper.NewSimpleReachHelper(db),
	// 	helper.NewSimpleListHelper(db),
	// 	helper.NewSimpleInsertHelper(db),
	// )

	ID, _ := domain.NewUUIDFromString(util.RootID)

	start := time.Now()
	rs := g.Reach(*ID)
	pp.Println(time.Since(start).Seconds())
	pp.Println(rs)
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
