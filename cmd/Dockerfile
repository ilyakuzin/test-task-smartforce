FROM golang:1.21-alpine AS build 

WORKDIR /src
COPY go.mod go.sum ./
COPY . .
RUN go build -o ./cmd/main.go