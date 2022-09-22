# syntax = docker/dockerfile:1

FROM golang:1.19-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.* .

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=linux GOARCH=amd64 go build -o /build/server hello/server/main.go

FROM scratch
COPY --from=build /build/server /
ENTRYPOINT ["/server"]
