.PHONY: all

build:
	go build cmd/phi/main.go

exec:
	./main

run:
	go run cmd/phi/main.go

all: build exec
