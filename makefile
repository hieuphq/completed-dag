.PHONY: run

build:
	go build -o dag cmd/dag/*.go

run: build
	./dag; rm dag