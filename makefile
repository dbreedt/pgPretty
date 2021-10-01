GOFLAGS=-i
LDFLAGS=-ldflags="-s -w"

all: build

build:
	go build $(GOFLAGS) $(LDFLAGS) -o pgPretty *.go

run:
	go run *.go

test:
	cd test && go test

cover:
	go test -coverprofile /tmp/pgPretty.out -covermode=atomic -coverpkg github.com/dbreedt/pgPretty/... ./test/...
	go tool cover -html=/tmp/pgPretty.out -o /tmp/cover.html && google-chrome /tmp/cover.html

.PHONY:	build test cover