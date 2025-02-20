FROM golang:1.24-alpine AS builder

WORKDIR /build/

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /build/main /build/main.go


FROM alpine:latest

WORKDIR /app/

COPY --from=builder /build/main /app/main

USER nobody:nobody
CMD ["/app/main"]
