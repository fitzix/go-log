FROM golang:latest AS builder

ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.io

WORKDIR /github.com/fitzix/go-log

COPY . .
RUN go mod download
RUN go build -o 'go-log' server.go

FROM scratch

WORKDIR /go/bin
COPY --from=builder /github.com/fitzix/go-log/go-log .
# COPY config.toml .
# ENTRYPOINT [ "/go/bin/go-log" ]