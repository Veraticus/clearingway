FROM golang:1.23.1 AS builder
WORKDIR /src
COPY go.mod go.sum /src/
RUN go mod download
COPY . /src
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make clearingway

FROM gcr.io/distroless/static-debian12
WORKDIR /clearingway
COPY --from=builder /src/clearingway .
COPY --from=builder /src/config.yaml .
USER nonroot:nonroot
ENTRYPOINT ["/clearingway/clearingway"]
