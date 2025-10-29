# Hyperplane - Open Key Management System

[Documentation](https://docs.open-mks.com/) | [Hyperplane](https://www.hyperplane.sh) | [Developer Documentation](https://github.com/hyperplane-sh/openkms/wiki)

OpenKMS is an open-source Key Management System (KMS) designed to provide secure key storage and management for
applications and services. It offers a robust and scalable solution for handling cryptographic keys, ensuring data
security.

## Getting Started with OpenKMS

**Work in progress**

Getting started with OpenKMS is straightforward. You can deploy OpenKMS using Docker with the following
`docker-compose.yml`.

```yaml
volumes:
  openkms-daemon-data:

services:
  daemon:
    image: ghcr.io/hyperplane-sh/openkms:daemon-0.0.1
    container_name: openkms-daemon
    volumes:
      - openkms-daemon-data:/etc/hyperplane/openkms
  cli:
    image: ghcr.io/hyperplane-sh/openkms:cli-0.0.1
    container_name: openkms-cli
    depends_on:
      - daemon
```

