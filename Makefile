

all: run

PROTO_FILE=order.proto
SERVER_FILE=cmd/server/main.go
# Установка зависимостей
.PHONY: install-deps
install-deps:
	@echo "Installing protobuf dependencies..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	export PATH="$PATH:$(go env GOPATH)/bin"

# Компиляция
.PHONY: generate
generate: install-deps
	@echo "Compiling .proto files"
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative ./${PROTO_FILE}

.PHONY: build
build: generate
	@echo "Building binary..."
	go build -o main ${SERVER_FILE}

.PHONY: run
run: build
	@echo "Running server..."
	./main