# CLI Command Reference

## Global Options

```bash
--config string     # Path to config file
--debug            # Enable debug mode
--help             # Show help for command
--version          # Show version information
```

## Commands

### generate
Generate output from a template file.

```bash
templater generate [flags]
```

#### Flags
```bash
-t, --template string    # Template file path
-d, --data string       # Data file path (JSON/YAML)
-o, --output string     # Output file path
--stdout               # Write output to stdout
--backup               # Create backup before overwriting
--no-backup            # Disable backup creation
```

#### Examples
```bash
# Basic usage
templater generate -t template.tmpl -d data.json -o output.go

# Output to stdout
templater generate -t template.tmpl -d data.json --stdout

# With backup
templater generate -t template.tmpl -d data.json -o output.go --backup
```

### generate-all
Generate output from multiple templates in a directory.

```bash
templater generate-all [flags]
```

#### Flags
```bash
-t, --template-dir string    # Template directory path
-d, --data-dir string       # Data directory path
-o, --output-dir string     # Output directory path
-r, --recursive            # Process subdirectories
-m, --monitor              # Monitor resource usage
```

#### Examples
```bash
# Process all templates
templater generate-all -t ./templates -d ./data -o ./generated

# With recursion
templater generate-all -t ./templates -d ./data -o ./generated -r

# With monitoring
templater generate-all -t ./templates -d ./data -o ./generated -m
```

### watch
Watch for changes and regenerate automatically.

```bash
templater watch [flags]
```

#### Flags
```bash
-t, --template string    # Template file path
-d, --data string       # Data file path
-o, --output string     # Output file path
-i, --interval string   # Watch interval (default "1s")
--stdout               # Write output to stdout
```

#### Examples
```bash
# Basic watch
templater watch -t template.tmpl -d data.json -o output.go

# Custom interval
templater watch -t template.tmpl -d data.json -o output.go -i 5s

# Watch with stdout
templater watch -t template.tmpl -d data.json --stdout
```

### run
Run a generated Go file.

```bash
templater run [flags] <file>
```

#### Flags
```bash
--args string    # Arguments to pass to the program
--env string     # Environment variables (key=value)
```

#### Examples
```bash
# Run generated file
templater run output.go

# With arguments
templater run output.go --args "arg1 arg2"

# With environment
templater run output.go --env "DEBUG=true"
```

## Configuration

### Environment Variables
```bash
TEMPLATER_CONFIG    # Path to config file
TEMPLATER_DEBUG     # Enable debug mode
```

### Config File (.templater.yaml)
```yaml
defaults:
  output_dir: "generated"
  watch_interval: "1s"
  backup: true
  language: "en"
```

## Exit Codes

- `0`: Success
- `1`: General error
- `2`: Invalid arguments
- `3`: File not found
- `4`: Template error
- `5`: Data error
- `6`: Output error

## Examples

### Complex Workflow
```bash
# Generate multiple files
templater generate-all -t ./templates -d ./data -o ./generated

# Watch for changes
templater watch -t main.tmpl -d config.yaml -o main.go

# Run generated code
templater run main.go --args "serve --port 8080"
```

### Using Configuration
```bash
# Load custom config
TEMPLATER_CONFIG=./config.yaml templater generate -t template.tmpl

# Enable debug mode
TEMPLATER_DEBUG=true templater watch -t template.tmpl
``` 