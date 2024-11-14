APP_NAME=taskPet
BIN_DIR=bin
SRC_DIR=./cmd/main.go

.PHONY: build run test

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(SRC_DIR)

run: build
	./$(BIN_DIR)/$(APP_NAME)

test:
	go test -v ./...
