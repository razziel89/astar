SHELL := /bin/bash

default: lint

build: astar

astar: *.go
	go build -o astar ./...

.PHONY: lint
lint:
	golangci-lint run .

test: .test.log

.test.log: go.* *.go
	set -o pipefail && \
		go test ./... | tee .test.log || \
		rm .test.log

coverage.out: go.* *.go
	go test -coverprofile=coverage.out 
	go tool cover -html=coverage.out -o coverage.html

.PHONY: coverage
coverage: coverage.out
	xdg-open coverage.html

.PHONY: performance
performance:
	cd ./tests && echo MAPPED GRAPH && QUIET=$${QUIET-1} GRAPH_TYPE=MAPPED go run .
	cd ./tests && echo HEAPED GRAPH && QUIET=$${QUIET-1} GRAPH_TYPE=HEAPED go run .

bench: .bench.log

.bench.log: go.* *.go
	set -o pipefail && \
		go test -bench=. ./... | tee .bench.log || \
		rm .bench.log
