FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o loans-service ./cmd/api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/loans-service .
EXPOSE 8081
CMD ["./loans-service"]