all: install

.PHONY: build
build: OUT?=build/faucet
build:
	mkdir -p build
	go build -o=$(OUT) ./cmd/faucet

install:
	go install -mod=readonly ./cmd/faucet