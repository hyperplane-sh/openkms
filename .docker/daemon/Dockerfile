FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .
RUN go build -o openkms-daemon ./cmd/daemon

FROM alpine:3.22 AS production

# Create openkms user and group.
#
RUN addgroup -S openkms && \
    adduser -S -G openkms openkms

# Setup openkms directory
#
RUN mkdir -p /etc/hyperplane/openkms
RUN mkdir -p /etc/hyperplane/openkms/configs
RUN mkdir -p /etc/hyperplane/openkms/data
RUN mkdir -p /etc/hyperplane/openkms/logs
RUN mkdir -p /etc/hyperplane/openkms/binaries

# Set permissions
#
RUN chmod -R 755 /etc/hyperplane/openkms

# Change ownership to openkms user
#
RUN chown -R openkms:openkms /etc/hyperplane/openkms

WORKDIR /app

COPY --from=builder /app/openkms-daemon /etc/hyperplane/openkms/binaries/openkms-daemon
COPY ./openkms.yaml /etc/hyperplane/openkms/configs/openkms.yaml

RUN chown openkms:openkms /etc/hyperplane/openkms/binaries/openkms-daemon
RUN chown openkms:openkms /etc/hyperplane/openkms/configs/openkms.yaml
RUN chmod 755 /etc/hyperplane/openkms/binaries/openkms-daemon
RUN chmod 644 /etc/hyperplane/openkms/configs/openkms.yaml

USER openkms

ENTRYPOINT ["/etc/hyperplane/openkms/binaries/openkms-daemon"]