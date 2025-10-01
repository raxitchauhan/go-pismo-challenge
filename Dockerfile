FROM golang:1.25.1-alpine3.22 AS build

WORKDIR /go/src/app

COPY . .