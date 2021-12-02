CURRENT_DIR = $(shell pwd)

build:
	go build ./...

test: build 
	golangci-lint run
	golangci-lint run --build-tags release

	GOBIN=$(CURRENT_DIR)/bin GO111MODULE=off go get gotest.tools/gotestsum
	$(CURRENT_DIR)/bin/gotestsum --format dots -- -count=1 -parallel 8 -race -v ./...
	$(CURRENT_DIR)/bin/gotestsum --format dots -- -count=1 -parallel 8 -race -v -tags release ./...

bench: build
	# Run benchmarks with -race for testing purposes (since -race adds overhead to real benchmarks).
	go test -bench=. -benchmem -count=1 -parallel 8 -race
	go test -bench=. -benchmem -count=1 -parallel 8 -race -tags release 
	#
	# *** Run for real ***
	#
	go test -bench=. -benchmem -count=1 -parallel 8 
	go test -bench=. -benchmem -count=1 -parallel 8 -tags release

	go test -bench=. -benchmem -count=1 -cpu 8 -parallel 8
	go test -bench=. -benchmem -count=1 -cpu 8 -parallel 8 -tags release
