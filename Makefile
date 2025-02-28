SERVER_BIN=bin/server
CLIENT_BIN=bin/client
PROTO_DIR=api/proto
PROTO_OUT=api/pb
PROTO_FILES=$(wildcard $(PROTO_DIR)/*.proto)

GO=go
GO_FLAGS=-mod=vendor

.PHONY: all build run test proto clean

all: build

run-server:
	$(GO) run ./cmd/server/main.go

run-client:
	$(GO) run ./cmd/client/main.go

build-server:
	$(GO) build -o $(SERVER_BIN) ./cmd/server

build-client:
	$(GO) build -o $(CLIENT_BIN) ./cmd/client

proto:
	@mkdir -p $(PROTO_OUT)
	protoc --go_out=$(PROTO_OUT) --go_opt=paths=source_relative --go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative --proto_path=$(PROTO_DIR) $(PROTO_FILES)

test:
	$(GO) vet ./...
	$(GO) test ./...

clean:
	rm -rf bin/