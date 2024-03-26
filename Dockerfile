FROM golang:1.20.7 as builder

WORKDIR /go/mx-chain-ws-connector-firehose-go
COPY . .
RUN go mod tidy
# Multiversx node
WORKDIR /go/mx-chain-ws-connector-firehose-go/cmd/connector
RUN go build -v -ldflags="-X main.appVersion=$(git describe --tags --long --dirty)"
RUN mv connector /usr/bin/connector

# ===== SECOND STAGE ======
FROM ubuntu:22.04
RUN apt-get update && apt-get upgrade -y
COPY --from=builder /usr/bin/connector /app/connector
WORKDIR /app
ENTRYPOINT ["./connector"]