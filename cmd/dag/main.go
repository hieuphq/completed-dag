package main

import (
	"github.com/k0kubun/pp"
)

func main() {
	ns := GenerateNodes(100)
	nwes := Nodes(ns).GenerateEdges()

	db, err := NewDB("./temp.leveldb")
	if err != nil {
		pp.Println(err)
		return
	}

	defer db.Close()

	println(nwes)
}
