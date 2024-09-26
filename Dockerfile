FROM golang:1.22.1-alpine3.18 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main cmd/app/main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/main .

COPY internal/pkg/config/auth.conf internal/pkg/config/auth.conf
COPY internal/pkg/config/auth.csv internal/pkg/config/auth.csv


CMD ["/app/main"]