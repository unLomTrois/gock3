run:
	go run ./cmd/main.go

build:
	cd ./cmd && go build main.go

.DEFAULT_GOAL := build
