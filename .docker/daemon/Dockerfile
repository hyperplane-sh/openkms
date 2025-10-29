FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .
RUN go build -o openkms-daemon ./cmd/daemon

FROM alpine:3.22 AS production

WORKDIR /app

COPY --from=builder /app/openkms-daemon ./openkms-daemon

ENTRYPOINT ["./openkms-daemon"]