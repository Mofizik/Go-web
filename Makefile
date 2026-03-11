all: run

PROTO_FILE=order.proto
SERVER_FILE=cmd/server/main.go

# Установка зависимостей
.PHONY: install-deps
install-deps:
	@echo "Installing protobuf dependencies..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	export PATH="$$PATH:$(shell go env GOPATH)/bin"

# Компиляция
.PHONY: generate
generate:
	@echo "Compiling .proto files"
	protoc --go_out=. \
       --go-grpc_out=. \
       --grpc-gateway_out=. \
       -I . -I third_party \
       ./order.proto
	       ./${PROTO_FILE}

.PHONY: build
build:
	@echo "Building binary..."
	go build -o main ${SERVER_FILE}

.PHONY: run
run: build
	@echo "Running server..."
	./main