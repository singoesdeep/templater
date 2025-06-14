# Installation Guide

## Prerequisites
- Go 1.21 or higher
- Git

## Installation Methods

### 1. Using Go Install
```bash
go install github.com/singoesdeep/templater@latest
```

### 2. Building from Source
```bash
# Clone the repository
git clone https://github.com/singoesdeep/templater.git
cd templater

# Build the project
go build -o templater ./cmd/templater

# Move the binary to your PATH (optional)
mv templater /usr/local/bin/
```

### 3. Using Precompiled Binaries
Download the latest release from the [Releases page](https://github.com/singoesdeep/templater/releases) and extract the binary for your platform.

## Configuration

### Basic Configuration
Create a `.templater.yaml` file in your project root:

```yaml
defaults:
  output_dir: "generated"
  watch_interval: "1s"
  backup: true
  language: "en"
```

### Environment Variables
- `TEMPLATER_CONFIG`: Path to custom config file
- `TEMPLATER_DEBUG`: Enable debug mode (set to "true")

## Verification

Verify the installation:
```bash
templater --version
```

## Troubleshooting

### Common Issues

1. **Command not found**
   - Ensure the binary is in your PATH
   - Try using the full path to the binary

2. **Permission denied**
   - Check file permissions
   - Run with appropriate privileges

3. **Go version mismatch**
   - Update Go to version 1.21 or higher
   - Check version with `go version`

## Next Steps

- Read the [Quick Start Guide](QUICKSTART.md)
- Check the [Template Syntax Reference](TEMPLATE_SYNTAX.md)
- Review the [CLI Command Reference](CLI_REFERENCE.md) 