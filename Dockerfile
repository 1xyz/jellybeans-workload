FROM golang:1.13.8-alpine3.11 AS builder

RUN apk update && apk add make git build-base curl && \
     rm -rf /var/cache/apk/*

ADD . /go/src/github.com/1xyz/jellybeans-workload
WORKDIR /go/src/github.com/1xyz/jellybeans-workload
RUN make release/linux

###

FROM alpine:latest AS jellybeans-workload

RUN apk update && apk add ca-certificates bash
WORKDIR /root/
COPY --from=builder /go/src/github.com/1xyz/jellybeans-workload/bin/linux/jellybeans-workload .