FROM golang:1.22-alpine3.19 as build
WORKDIR /go/src/github.com/isdzulqor/donation-hub

COPY go.mod go.sum ./
RUN go mod download -x

COPY ./.air.toml ./.air.toml
COPY ./cmd ./cmd
COPY ./internal ./internal

RUN go install github.com/cosmtrek/air@latest

ENTRYPOINT [ "air", "-c", "./.air.toml" ]