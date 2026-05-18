# Not used by goreleaser, for local development
FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /usr/bin/seapack cmd/cli/main.go

FROM ghcr.io/gitlayzer/seapack-alpine:latest

COPY --from=builder /usr/bin/seapack /usr/bin/seapack
ENTRYPOINT ["/usr/bin/seapack", "frontend"]
