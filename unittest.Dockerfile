FROM golang:1.25.1-alpine3.22  AS golang

RUN apk add --no-cache \
	bash build-base

RUN go env -w CGO_ENABLED=1

WORKDIR $GOPATH/src/go-pismo-challenge

COPY --from=golangci/golangci-lint:v2.5-alpine /usr/bin/golangci-lint /usr/local/bin/golangci-lint

COPY .golangci.yml .

COPY . .