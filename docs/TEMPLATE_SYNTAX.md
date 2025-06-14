# Template Syntax Reference

## Basic Syntax

### Variables
```go
{{.VariableName}}  // Access a variable
{{.Nested.Variable}}  // Access nested variables
{{$local := .Value}}  // Define local variable
```

### Actions

#### If Statements
```go
{{if .Condition}}
    // Content when condition is true
{{else if .OtherCondition}}
    // Content when other condition is true
{{else}}
    // Content when all conditions are false
{{end}}
```

#### Range Loops
```go
{{range .Items}}
    {{.}}  // Current item
{{end}}

{{range $index, $item := .Items}}
    {{$index}}: {{$item}}  // Index and item
{{end}}
```

#### With Blocks
```go
{{with .User}}
    {{.Name}}  // Access User.Name
{{end}}
```

### Functions

#### Built-in Functions
```go
{{len .Items}}  // Length of array/map
{{index .Array 0}}  // Array indexing
{{printf "%s" .Value}}  // Formatted string
{{html .Content}}  // HTML escaping
```

#### Custom Functions
```go
{{title .Name}}  // Capitalize first letter
{{default "empty" .Value}}  // Default value
{{join .Items ","}}  // Join array with separator
```

## Advanced Features

### Template Composition

#### Define Blocks
```go
{{define "header"}}
    <header>{{.Title}}</header>
{{end}}
```

#### Include Templates
```go
{{template "header" .}}
```

#### Nested Templates
```go
{{block "content" .}}
    Default content
{{end}}
```

### Pipeline Operations
```go
{{.Value | upper | trim}}  // Chain operations
{{.Name | default "Anonymous" | title}}  // Multiple operations
```

### Comments
```go
{{/* This is a comment */}}
```

## Data Types

### Supported Types
- Strings
- Numbers
- Booleans
- Arrays/Slices
- Maps
- Structs

### Type Handling
```go
{{if eq .Type "string"}}
    String value: {{.Value}}
{{else if eq .Type "number"}}
    Number value: {{.Value}}
{{end}}
```

## Best Practices

### Error Handling
```go
{{if .Error}}
    Error: {{.Error}}
{{else}}
    Success: {{.Result}}
{{end}}
```

### Default Values
```go
{{.Value | default "empty"}}
{{.Count | default 0}}
```

### Conditional Rendering
```go
{{if and .Enabled .Visible}}
    Content is enabled and visible
{{end}}
```

## Security Considerations

### HTML Escaping
```go
{{.Content}}  // Auto-escaped
{{.Content | safeHTML}}  // Explicitly marked safe
```

### JavaScript Escaping
```go
{{.Value | js}}  // JavaScript escaping
```

## Examples

### Complex Template
```go
{{define "user-card"}}
    {{with .User}}
        <div class="user-card">
            {{if .Avatar}}
                <img src="{{.Avatar}}" alt="{{.Name}}">
            {{end}}
            <h2>{{.Name | title}}</h2>
            {{if .Bio}}
                <p>{{.Bio}}</p>
            {{end}}
            {{range .Tags}}
                <span class="tag">{{.}}</span>
            {{end}}
        </div>
    {{end}}
{{end}}
```

### Data Processing
```go
{{define "stats"}}
    {{$total := len .Items}}
    {{$active := 0}}
    {{range .Items}}
        {{if .Active}}
            {{$active = add $active 1}}
        {{end}}
    {{end}}
    <div>Total: {{$total}}</div>
    <div>Active: {{$active}}</div>
    <div>Inactive: {{subtract $total $active}}</div>
{{end}}
``` 