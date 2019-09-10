GOFLAGS=-i
LDFLAGS=-ldflags="-s -w"

all: build

build:
	go build $(GOFLAGS) $(LDFLAGS) -o pgPretty *.go

run:
	go run ./...

.PHONY:	build