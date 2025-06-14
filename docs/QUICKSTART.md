# Quick Start Guide

## Basic Usage

### 1. Create a Template
Create a file named `hello.tmpl`:
```go
package main

func {{.FuncName}}() string {
    return "{{.Message}}"
}
```

### 2. Create Data File
Create a file named `data.json`:
```json
{
    "FuncName": "Hello",
    "Message": "Hello, World!"
}
```

### 3. Generate Output
```bash
templater generate -t hello.tmpl -d data.json -o hello.go
```

The generated `hello.go` will contain:
```go
package main

func Hello() string {
    return "Hello, World!"
}
```

## Advanced Features

### Multiple Templates
```bash
# Process all templates in a directory
templater generate-all -t ./templates -d ./data -o ./generated
```

### Watch Mode
```bash
# Automatically regenerate on file changes
templater watch -t hello.tmpl -d data.json -o hello.go
```

### Using YAML Data
Create `data.yaml`:
```yaml
FuncName: Hello
Message: Hello from YAML!
```

```bash
templater generate -t hello.tmpl -d data.yaml -o hello.go
```

## Template Features

### Conditionals
```go
{{if .Enabled}}
    // This code will be included if .Enabled is true
{{else}}
    // This code will be included if .Enabled is false
{{end}}
```

### Loops
```go
{{range .Items}}
    // Process each item
    {{.Name}}
{{end}}
```

### Functions
```go
{{.FunctionName | title}}  // Capitalize first letter
{{.Value | default "empty"}}  // Use default value if empty
```

## Best Practices

1. **Template Organization**
   - Keep templates in a dedicated directory
   - Use meaningful file names
   - Group related templates together

2. **Data Management**
   - Validate data before processing
   - Use consistent data structures
   - Keep sensitive data separate

3. **Output Handling**
   - Use appropriate file extensions
   - Implement backup mechanisms
   - Monitor generated files

## Next Steps

- Explore [Template Syntax Reference](TEMPLATE_SYNTAX.md)
- Learn about [CLI Commands](CLI_REFERENCE.md)
- Check out [Example Projects](../examples/) 