FROM golang:1.20.5-alpine AS builder
WORKDIR /app
COPY go_api/. .
RUN go mod download
RUN go build -o main .

FROM cgr.dev/chainguard/wolfi-base:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
