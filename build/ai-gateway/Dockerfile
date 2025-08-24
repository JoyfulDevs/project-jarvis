FROM golang:1.25 AS builder
WORKDIR /app

COPY . .
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" -o ai-gateway ./cmd/ai-gateway

FROM scratch

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /app/ai-gateway /app/

ENTRYPOINT ["/app/ai-gateway"]
