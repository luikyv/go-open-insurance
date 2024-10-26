# Step 1: Build Stage.
FROM golang:1.23-alpine AS builder

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main ./cmd/server

# Step 2: Run Stage.
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main ./cmd/server/
COPY ./keys/ ./keys/
COPY ./templates/ ./templates/

EXPOSE 80

CMD ["./cmd/server/main"]
