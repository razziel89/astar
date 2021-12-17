default: lint

build: astar

astar: *.go
	go build -o astar ./...

.PHONY: lint
lint:
	golangci-lint run .

test: .test.log

.test.log: *.go go.*
	go test ./... | tee .test.log
