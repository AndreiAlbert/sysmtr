PROTO_SRC_DIR=../proto
PB_DEST_DIR=pb
BIN_DIR=bin
SERVER_BIN=$(BIN_DIR)/server
AGENT_BIN=$(BIN_DIR)/agent

.PHONY: all gen build build-server build-agent run-server run-agent clean deps

all: gen build

gen:
	@echo -e "Generating Protocol Buffers\n"
	@mkdir -p $(PB_DEST_DIR)
	protoc --proto_path=$(PROTO_SRC_DIR) \
		--go_out=$(PB_DEST_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PB_DEST_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_SRC_DIR)/monitor.proto

build: build-server build-agent

test:
	@echo "Running tests..."
	go test -v ./...

build-server:
	@echo "Building server..."
	@mkdir -p $(BIN_DIR)
	go build -o $(SERVER_BIN) ./cmd/server


build-agent:
	@echo "Building agennt..."
	@mkdir -p $(BIN_DIR)
	go build -o $(AGENT_BIN) ./cmd/agent


run-server:
	go run ./cmd/server

run-agent:
	go run ./cmd/agent

clean:
	@echo "Cleaning up..."
	rm -rf $(BIN_DIR)

deps:
	@echo "Tidying and downloading Go modules..."
	go mod tidy
	go mod download
	@echo "Installing Go protocol buffer plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Checking for protoc..."
	@if ! command -v protoc >/dev/null 2>&1; then \
		echo "Error: 'protoc' compiler not found."; \
		echo "Please install it via your system package manager:"; \
		exit 1; \
	else \
		echo "protoc is installed."; \
	fi
