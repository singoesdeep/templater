# Templater API Documentation

## Overview

The Templater API provides a set of functions for template processing, code generation, and execution. This document describes the available functions and their usage.

## Core Functions

### Render

```go
func Render(templatePath string, data map[string]string) (string, error)
```

Processes a template file with the given data and returns the result.

**Parameters:**
- `templatePath`: Path to the template file
- `data`: Map of string key-value pairs for template variables

**Returns:**
- `string`: The processed template output
- `error`: Any error that occurred during processing

**Example:**
```go
data := map[string]string{
    "Name": "World",
    "Message": "Hello",
}
result, err := templater.Render("template.tmpl", data)
if err != nil {
    log.Fatal(err)
}
fmt.Println(result)
```

### Generate

```go
func Generate(templatePath, dataPath, outputPath string) error
```

Processes a template file with data from a file and saves the result.

**Parameters:**
- `templatePath`: Path to the template file
- `dataPath`: Path to the data file (JSON or YAML)
- `outputPath`: Path where the result should be saved

**Returns:**
- `error`: Any error that occurred during processing

**Example:**
```go
err := templater.Generate("template.tmpl", "data.json", "output.go")
if err != nil {
    log.Fatal(err)
}
```

### RunGoCode

```go
func RunGoCode(code string) error
```

Executes the generated Go code in a temporary file.

**Parameters:**
- `code`: The Go code to execute

**Returns:**
- `error`: Any error that occurred during execution

**Example:**
```go
code := `package main
import "fmt"
func main() {
    fmt.Println("Hello, World!")
}`
err := templater.RunGoCode(code)
if err != nil {
    log.Fatal(err)
}
```

### LoadData

```go
func LoadData(path string) (map[string]string, error)
```

Loads data from a JSON or YAML file.

**Parameters:**
- `path`: Path to the data file

**Returns:**
- `map[string]string`: The loaded data
- `error`: Any error that occurred during loading

**Example:**
```go
data, err := templater.LoadData("data.json")
if err != nil {
    log.Fatal(err)
}
```

### ValidateTemplate

```go
func ValidateTemplate(path string) error
```

Validates a template file.

**Parameters:**
- `path`: Path to the template file

**Returns:**
- `error`: Any validation errors

**Example:**
```go
err := templater.ValidateTemplate("template.tmpl")
if err != nil {
    log.Fatal(err)
}
```

### ValidateData

```go
func ValidateData(templatePath string, data map[string]string) error
```

Validates data against a template.

**Parameters:**
- `templatePath`: Path to the template file
- `data`: The data to validate

**Returns:**
- `error`: Any validation errors

**Example:**
```go
err := templater.ValidateData("template.tmpl", data)
if err != nil {
    log.Fatal(err)
}
```

## Error Handling

All functions return errors that should be checked and handled appropriately. Common error types include:

- Template validation errors
- Data loading errors
- File system errors
- Execution errors

## Security Considerations

The API includes several security features:

1. Template content validation
2. Data sanitization
3. Output path validation
4. Secure temporary file handling

## Best Practices

1. Always validate templates before processing
2. Sanitize input data
3. Use proper error handling
4. Clean up temporary files
5. Validate output paths

## Examples

### Complete Example

```go
package main

import (
    "fmt"
    "log"

    "github.com/singoesdeep/templater"
)

func main() {
    // Load and validate data
    data, err := templater.LoadData("data.json")
    if err != nil {
        log.Fatal(err)
    }

    // Validate template
    if err := templater.ValidateTemplate("template.tmpl"); err != nil {
        log.Fatal(err)
    }

    // Validate data against template
    if err := templater.ValidateData("template.tmpl", data); err != nil {
        log.Fatal(err)
    }

    // Render template
    result, err := templater.Render("template.tmpl", data)
    if err != nil {
        log.Fatal(err)
    }

    // Save result
    if err := templater.Generate("template.tmpl", "data.json", "output.go"); err != nil {
        log.Fatal(err)
    }

    // Run the generated code
    if err := templater.RunGoCode(result); err != nil {
        log.Fatal(err)
    }
}
```

## Troubleshooting

### Common Issues

1. **Template Not Found**
   - Ensure the template path is correct
   - Check file permissions
   - Verify the template file exists

2. **Invalid Template Syntax**
   - Check for missing closing tags
   - Verify template variable names
   - Ensure proper Go template syntax

3. **Data Loading Errors**
   - Verify JSON/YAML syntax
   - Check file permissions
   - Ensure data structure matches template

4. **Output Path Issues**
   - Check directory permissions
   - Verify path is within allowed directories
   - Ensure sufficient disk space

### Debugging Tips

1. Use `ValidateTemplate` to check template syntax
2. Use `ValidateData` to verify data structure
3. Check error messages for specific issues
4. Enable debug logging for detailed information

## Performance Considerations

1. Use concurrent processing for multiple templates
2. Enable caching for frequently used templates
3. Monitor memory usage during processing
4. Clean up temporary files regularly

## Version Compatibility

The API is designed to maintain backward compatibility. When breaking changes are necessary, they will be introduced in new major versions with appropriate migration guides. 