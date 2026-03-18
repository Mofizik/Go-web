FROM golang:1.25.6-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server ./cmd/server/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /out/server /app/server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 50051
EXPOSE 8080

CMD ["/app/server"]