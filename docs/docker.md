# Docker Support

This document describes how to use `templater` with Docker.

## Quick Start

1. Build the Docker image:
   ```bash
   docker build -t templater .
   ```

2. Run the container:
   ```bash
   docker run -v $(pwd)/templates:/app/templates -v $(pwd)/output:/app/output templater generate -t /app/templates/example.tmpl -o /app/output/result.go
   ```

## Using Docker Compose

1. Start the service:
   ```bash
   docker-compose up -d
   ```

2. Execute commands:
   ```bash
   docker-compose exec templater generate -t /app/templates/example.tmpl -o /app/output/result.go
   ```

## Volume Mounting

The Docker setup supports two main volume mounts:

- `/app/templates`: Mount your template files here
- `/app/output`: Mount your output directory here

Example:
```bash
docker run -v /path/to/templates:/app/templates -v /path/to/output:/app/output templater
```

## Health Checks

The container includes a health check that runs every 30 seconds. It verifies that the `templater` binary is working correctly by running the `--version` command.

## Security

The container runs as a non-root user `templater` for enhanced security. All files in the container are owned by this user.

## Environment Variables

- `TZ`: Set the timezone (default: UTC)

## Building from Source

To build the Docker image from source:

1. Clone the repository:
   ```bash
   git clone https://github.com/singoesdeep/templater.git
   cd templater
   ```

2. Build the image:
   ```bash
   docker build -t templater .
   ```

## Development

For development, use the provided `docker-compose.yml`:

1. Start the development environment:
   ```bash
   docker-compose up -d
   ```

2. Make changes to your templates in the `templates` directory
3. The changes will be reflected in the container
4. Use `docker-compose logs -f` to monitor the application

## Troubleshooting

1. Permission Issues:
   - Ensure the mounted volumes have the correct permissions
   - The container user `templater` needs read/write access

2. Volume Mounting:
   - Use absolute paths for volume mounting
   - Check that the directories exist on the host

3. Health Check Failures:
   - Check the container logs: `docker logs <container_id>`
   - Verify that the binary is executable
   - Check for any permission issues 