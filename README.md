# Templater

A powerful template engine that supports multiple file formats including Markdown, HTML, and DOCX.

## Features

- **Multiple Format Support**: Process templates in Markdown, HTML, and DOCX formats
- **Variable Substitution**: Replace placeholders with actual values
- **Conditional Logic**: Support for if/else conditions in templates
- **Loops**: Iterate over arrays and maps
- **Image Support**: Handle image placeholders in templates
- **Format Preservation**: Maintain original formatting in DOCX files

## Installation

```bash
go install github.com/singoesdeep/templater@latest
```

## Usage

### Basic Usage

```bash
templater process -t template.md -d data.json -o output.md
```

### Command Line Options

- `-t, --template`: Template file path (required)
- `-d, --data`: Data file path (required)
- `-o, --output`: Output file path (required)
- `-f, --format`: Output format (markdown, html, or docx)

### Template Syntax

#### Variables

```markdown
Hello {{name}}!
```

#### Conditionals

```markdown
{{if condition}}
  This will be shown if condition is true
{{else}}
  This will be shown if condition is false
{{end}}
```

#### Loops

```markdown
{{range items}}
  - {{.name}}
{{end}}
```

#### Images

```markdown
![alt text]({{image.placeholder}})
```

### Data Format

The data file should be in JSON format:

```json
{
  "name": "John",
  "items": [
    {"name": "Item 1"},
    {"name": "Item 2"}
  ],
  "image": {
    "placeholder": "/path/to/image.jpg"
  }
}
```

### DOCX Templates

DOCX templates support all the same features as Markdown templates, with some additional considerations:

1. **Format Preservation**: The engine preserves the original formatting of the DOCX file, including:
   - Font styles (bold, italic, underline)
   - Font sizes
   - Colors
   - Paragraph formatting
   - Lists and tables

2. **Image Handling**: Images in DOCX templates can be replaced using the same syntax as in Markdown:
   ```markdown
   {{image.placeholder}}
   ```
   The image will be inserted with a default size of 2x2 inches.

## License

This project is licensed under the MIT License.

## Acknowledgments

- [gooxml](https://github.com/plutext/gooxml) - Go library for working with Office Open XML (OOXML) files
- [Cobra](https://github.com/spf13/cobra) - A Commander for modern Go CLI interactions
- [YAML.v3](https://github.com/go-yaml/yaml) - YAML support for the Go language
- [Go's text/template](https://pkg.go.dev/text/template) - Go's built-in template engine 