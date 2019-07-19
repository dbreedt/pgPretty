GOFLAGS=-i
LDFLAGS=-ldflags="-s -w"

all: build

build:
	go build $(GOFLAGS) $(LDFLAGS) -o pgPretty main.go

.PHONY:	build