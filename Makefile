SHELL := /bin/bash

default: lint

.PHONY: setup
setup:
	go mod download
	go mod tidy

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

coverage.html: go.* *.go
	go test -covermode=count -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

coverage_badge_report.out: go.* *.go
	go test -covermode=count -coverprofile=coverage.out
	go tool cover -func=coverage.out -o=coverage_badge_report.out

.PHONY: coverage
coverage: coverage.html
	xdg-open coverage.html

.PHONY: performance
performance:
	cd ./tests && echo MAPPED GRAPH && QUIET=$${QUIET-1} GRAPH_TYPE=MAPPED go run .
	cd ./tests && echo HEAPED GRAPH && QUIET=$${QUIET-1} GRAPH_TYPE=HEAPED go run .

.PHONY: readme_test
readme_test:
	# Ensure the example in the readme actually works. Do this by extracting the
	# quoted code from the readme using awk into a temporary directory, then
	# initialise a temporary go module there, then run the example.
	repodir=$$(pwd) && \
	tmpdir=$$(mktemp -d) && \
	awk 'BEGIN{printme=false}{if($$1 == "```"){printme=0};if(printme){print};if($$1 == "```go"){printme=1};}' README.md > $$tmpdir/main.go && \
	cd $$tmpdir && go mod init main  && \
	echo "replace github.com/razziel89/astar => $${repodir}" >> go.mod && \
	go mod tidy && \
	go run main.go && \
	rm -r $$tmpdir

bench: .bench.log

.bench.log: go.* *.go
	set -o pipefail && \
		go test -bench=. ./... | tee .bench.log || \
		rm .bench.log
