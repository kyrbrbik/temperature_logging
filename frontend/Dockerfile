FROM golang:1.20.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .
COPY ./static/images ./static/images
WORKDIR /app/templates
COPY ./templates .

FROM cgr.dev/chainguard/wolfi-base:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static/ ./static/
EXPOSE 8090
CMD ["./main"]

