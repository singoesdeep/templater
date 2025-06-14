# 🎨 Templater

A powerful, secure, and flexible template processing tool for generating code and text files from templates using Go's text/template engine.

[![Go Report Card](https://goreportcard.com/badge/github.com/singoesdeep/templater)](https://goreportcard.com/report/github.com/singoesdeep/templater)
[![GoDoc](https://godoc.org/github.com/singoesdeep/templater?status.svg)](https://godoc.org/github.com/singoesdeep/templater)
[![License](https://img.shields.io/github/license/singoesdeep/templater)](LICENSE)
[![Release](https://img.shields.io/github/v/release/singoesdeep/templater)](https://github.com/singoesdeep/templater/releases)

## ✨ Features

- 🚀 **Fast & Efficient**: Process templates in under 100ms for files under 1MB
- 🔒 **Secure**: Sandboxed environment with dangerous command prevention
- 📦 **Multiple Formats**: Support for JSON and YAML data sources
- 🔄 **Watch Mode**: Automatic regeneration on file changes
- 🎯 **CLI & Library**: Use as a command-line tool or import as a Go library
- 🐳 **Docker Support**: Ready-to-use Docker image
- 🔌 **Plugin System**: Extend functionality with custom plugins

## 🚀 Quick Start

### Installation

```bash
# Using Go
go install github.com/singoesdeep/templater@latest

# Using Docker
docker pull singoesdeep/templater:latest
```

### Basic Usage

```bash
# Generate output from a template
templater generate -t template.tmpl -d data.json -o output.go

# Watch for changes
templater watch -t template.tmpl -d data.json -o output.go

# Process multiple templates
templater generate-all -t templates/ -d data/ -o output/
```

### As a Library

```go
package main

import (
    "fmt"
    "github.com/singoesdeep/templater"
)

func main() {
    // Load data
    data, err := templater.LoadData("data.json")
    if err != nil {
        log.Fatal(err)
    }

    // Render template
    result, err := templater.Render("template.tmpl", data)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(result)
}
```

## 📚 Documentation

- [Installation Guide](docs/INSTALL.md)
- [Quick Start Tutorial](docs/QUICKSTART.md)
- [Template Syntax](docs/TEMPLATE_SYNTAX.md)
- [CLI Reference](docs/CLI.md)
- [API Documentation](docs/API.md)
- [Security Guide](docs/SECURITY.md)

## 🔧 Configuration

Create a `.templater.yaml` file in your project root:

```yaml
defaults:
  output_dir: "generated"
  watch_interval: "1s"
  backup: true
  language: "en"
```

## 🐳 Docker Usage

```bash
# Run with Docker
docker run -v $(pwd)/templates:/app/templates \
          -v $(pwd)/output:/app/output \
          singoesdeep/templater generate \
          -t /app/templates/example.tmpl \
          -o /app/output/result.go
```

## 🔌 Plugin System

Create custom plugins to extend functionality:

```go
package main

import "github.com/singoesdeep/templater/pkg/templater/plugin"

type MyPlugin struct{}

func (p *MyPlugin) Name() string {
    return "myplugin"
}

func (p *MyPlugin) Process(input string) (string, error) {
    // Process input
    return input, nil
}

func main() {
    plugin := &MyPlugin{}
    // Register plugin
}
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Go's text/template](https://pkg.go.dev/text/template)
- [Cobra](https://github.com/spf13/cobra)
- [YAML.v3](https://github.com/go-yaml/yaml)

## 📞 Support

- [GitHub Issues](https://github.com/singoesdeep/templater/issues)
- [Documentation](https://github.com/singoesdeep/templater/tree/main/docs)
- [Discussions](https://github.com/singoesdeep/templater/discussions)

---

Made with ❤️ by Singoesdeep 