

all: run

PROTO_FILE=order.proto
SERVER_FILE=cmd/server/main.go
# Установка зависимостей
.PHONY: install-deps
install-deps:
	@echo "Installing protobuf dependencies..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Компиляция
.PHONY: protoc
protoc: install-deps
	@echo "Compiling .proto files"
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./${PROTO_FILE}

# Создание сервера
.PHONY: build
build: protoc
	@echo "Building binary..."
	go build -o main ${SERVER_FILE}

# Запуск сервера
.PHONY: run
run: build
	@echo "Running server..."
	./main