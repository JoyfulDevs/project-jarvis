FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY . .
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" -o weekly-report ./cmd/weekly-report

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/weekly-report /app/

ENTRYPOINT ["/app/weekly-report"]
