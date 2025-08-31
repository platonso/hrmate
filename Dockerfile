FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/api

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/server .

CMD ["./server"]
