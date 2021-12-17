default: lint

build: astar

astar: *.go
	go build -o astar ./...

.PHONY: lint
lint:
	golangci-lint run .
