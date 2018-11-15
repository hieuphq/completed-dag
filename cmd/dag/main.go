package main

import (
	"github.com/k0kubun/pp"
)

func main() {
	ns := GenerateNodes(100)
	nwes := Nodes(ns).GenerateEdges()

	pp.Println(nwes)
}
