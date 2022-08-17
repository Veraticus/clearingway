FROM golang:1.18.1 AS builder
WORKDIR /src
COPY . /src
RUN make build

FROM alpine:3.16.0
WORKDIR /clearingway
COPY --from=builder /src/clearingway .
ENTRYPOINT /clearingway/clearingway
