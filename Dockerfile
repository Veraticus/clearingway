FROM golang:1.19.0 AS builder
WORKDIR /src
COPY . /src
RUN make build

FROM alpine:3.16.0
WORKDIR /clearingway
COPY --from=builder /src/clearingway .
COPY --from=builder /src/config.yaml .
ENTRYPOINT /clearingway/clearingway
