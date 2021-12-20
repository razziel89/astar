default: lint

build: astar

astar: *.go
	go build -o astar ./...

.PHONY: lint
lint:
	golangci-lint run .

test: .test.log

.test.log: go.* *.go
	go test ./... | tee .test.log

coverage.out: go.* *.go
	go test -coverprofile=coverage.out 
	go tool cover -html=coverage.out -o coverage.html

.PHONY: coverage
coverage: coverage.out
	xdg-open coverage.html
