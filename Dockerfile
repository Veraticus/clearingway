FROM golang:1.22.1 AS builder
WORKDIR /src
COPY . /src
RUN make build

FROM alpine:3.19.1
WORKDIR /clearingway
COPY --from=builder /src/clearingway .
COPY --from=builder /src/config.yaml .
ENTRYPOINT /clearingway/clearingway
