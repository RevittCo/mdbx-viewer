# syntax = docker/dockerfile:1.2
FROM docker.io/library/golang:1.23-alpine3.20 AS base

RUN apk --no-cache add build-base linux-headers git bash ca-certificates libstdc++

WORKDIR /api

COPY go.mod .
COPY go.sum .

RUN go mod download

FROM base AS builder

COPY . .

RUN go build -o apiexecutable ./cmd/api

FROM docker.io/library/alpine:3.20

COPY --from=builder /api/apiexecutable /apiexecutable

EXPOSE 8080

CMD ["/apiexecutable"]