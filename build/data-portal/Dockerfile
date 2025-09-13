FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY . .
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" -o data-portal ./cmd/data-portal

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/data-portal /app/

ENTRYPOINT ["/app/data-portal"]
