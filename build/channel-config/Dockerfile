FROM golang:1.25 AS builder
WORKDIR /app

COPY . .
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" -o channel-config ./cmd/channel-config

FROM scratch

COPY --from=builder /app/channel-config /app/

ENTRYPOINT ["/app/channel-config"]
