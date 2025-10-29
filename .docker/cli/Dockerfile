FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .
RUN go build -o openkms ./cmd/cli

FROM alpine:3.22 AS production

WORKDIR /app

COPY --from=builder /app/openkms ./openkms

ENTRYPOINT ["./openkms"]
