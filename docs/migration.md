# Migration Guide to v1.0.0

This guide helps users migrate from previous versions to `templater` v1.0.0.

## Breaking Changes

### Template Processing
- Template validation is now stricter by default
- Template caching behavior has changed
- Template dependencies must be explicitly declared

### CLI Commands
- `--output` flag is now required for file generation
- Watch mode interval format has changed
- Configuration file format has been updated

### API Changes
- `Render` function signature has changed
- Plugin interface has been updated
- Error handling has been improved

## Migration Steps

### 1. Update Configuration

Old format:
```yaml
output_dir: "generated"
watch_interval: 1
backup: true
```

New format:
```yaml
defaults:
  output_dir: "generated"
  watch_interval: "1s"
  backup: true
  language: "en"
```

### 2. Update Template Files

Add template metadata:
```go
{{/* Template metadata
DependsOn:
  - other.tmpl
  - shared.tmpl
*/}}
```

### 3. Update CLI Usage

Old:
```bash
templater generate -t template.tmpl -d data.json
```

New:
```bash
templater generate -t template.tmpl -d data.json -o output.go
```

### 4. Update API Usage

Old:
```go
result, err := templater.Render(templatePath, data)
```

New:
```go
result, err := templater.Render(templatePath, data)
if err != nil {
    // Handle error
}
```

### 5. Update Plugin Code

Old:
```go
type Plugin interface {
    Name() string
    Process(string) string
}
```

New:
```go
type Plugin interface {
    Name() string
    Version() string
    Description() string
    Initialize() error
    Shutdown() error
}
```

## New Features

### 1. Enhanced Security
- Template content validation
- Sandbox environment
- Data sanitization
- File path validation

### 2. Improved Performance
- Template caching
- Concurrent processing
- Resource monitoring
- Memory optimization

### 3. Better Error Handling
- Detailed error messages
- Recovery mechanisms
- Backup system
- Validation checks

### 4. New Commands
- `templater validate` - Validate templates
- `templater watch` - Watch mode with status
- `templater plugins` - Plugin management

## Troubleshooting

### Common Issues

1. **Template Validation Errors**
   - Check template syntax
   - Verify dependencies
   - Update metadata format

2. **Configuration Issues**
   - Use new YAML format
   - Set required fields
   - Check file permissions

3. **Plugin Loading Errors**
   - Update plugin interface
   - Check version compatibility
   - Verify plugin paths

### Getting Help

- Check the documentation
- Review the changelog
- Open an issue
- Join the community

## Support

For additional help:
- Read the documentation
- Check the FAQ
- Join the community
- Open an issue 