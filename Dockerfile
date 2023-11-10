FROM golang:1.21.4 AS builder
WORKDIR /src
COPY . /src
RUN make build

FROM alpine:3.18.4
WORKDIR /clearingway
COPY --from=builder /src/clearingway .
COPY --from=builder /src/config.yaml .
ENTRYPOINT /clearingway/clearingway
