FROM golang:1.25.1-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o server ./cmd/app/

RUN go build -o migrator ./cmd/migrator/

FROM alpine:latest AS app

WORKDIR /app
COPY --from=builder /build/server .
CMD ["./server"]

FROM alpine:latest AS migrator

WORKDIR /app
COPY --from=builder /build/migrator .
COPY --from=builder /build/migrations/ ./migrations/
CMD ["./migrator"]
