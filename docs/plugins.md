# Plugin System Documentation

## Overview
The `templater` tool supports a plugin system that allows users to extend its functionality. Plugins can be dynamically loaded or statically registered within the application.

## Available Plugins

### UppercasePlugin
- **Version**: 1.0.0
- **Description**: Converts template output to uppercase.
- **Usage**: This plugin is automatically registered and can be used to process templates, converting the output to uppercase.

## How to Use Plugins
1. **Static Registration**: Plugins can be statically registered in the `main.go` file. This is useful for platforms where dynamic loading is not supported, such as Windows.
2. **Dynamic Loading**: Plugins can be loaded from a specified directory. Ensure the plugin is compiled as a shared library (e.g., `.so` for Linux/macOS).

## Example
To use the `UppercasePlugin`, ensure it is registered in the `main.go` file. The plugin will automatically process templates and convert the output to uppercase.

## Adding New Plugins
To add a new plugin, implement the `Plugin` interface and register it using the `plugin.Registry`. Ensure the plugin is compiled as a shared library if dynamic loading is desired.

## Troubleshooting
- Ensure the plugin is correctly registered and compiled.
- Check the console output for any errors related to plugin loading or registration.

# Templater Plugins

Templater supports plugins to extend its functionality. Plugins are Go packages that implement the `templater.Plugin` interface.

## Plugin Interface

A plugin must implement the following interface:

```go
type Plugin interface {
    Name() string
    Description() string
    Execute(input string) (string, error)
}
```

## Example: Uppercase Plugin

Here's a simple example of a plugin that converts text to uppercase:

```go
package uppercase

import (
    "strings"
)

type UppercasePlugin struct{}

func (p *UppercasePlugin) Name() string {
    return "uppercase"
}

func (p *UppercasePlugin) Description() string {
    return "Converts text to uppercase"
}

func (p *UppercasePlugin) Execute(input string) (string, error) {
    return strings.ToUpper(input), nil
}

func New() *UppercasePlugin {
    return &UppercasePlugin{}
}
```

### Usage

1. **Create a new directory for your plugin:**

   ```sh
   mkdir -p plugins/uppercase
   cd plugins/uppercase
   ```

2. **Create a `go.mod` file:**

   ```sh
   go mod init github.com/singoesdeep/templater/plugins/uppercase
   ```

3. **Create the plugin file (`uppercase.go`):**

   ```go
   package uppercase

   import (
       "strings"
   )

   type UppercasePlugin struct{}

   func (p *UppercasePlugin) Name() string {
       return "uppercase"
   }

   func (p *UppercasePlugin) Description() string {
       return "Converts text to uppercase"
   }

   func (p *UppercasePlugin) Execute(input string) (string, error) {
       return strings.ToUpper(input), nil
   }

   func New() *UppercasePlugin {
       return &UppercasePlugin{}
   }
   ```

4. **Create a `main.go` file to test the plugin:**

   ```go
   package main

   import (
       "fmt"
       "log"

       "github.com/singoesdeep/templater/plugins/uppercase"
   )

   func main() {
       plugin := uppercase.New()
       result, err := plugin.Execute("hello, world!")
       if err != nil {
           log.Fatal(err)
       }
       fmt.Println(result) // Output: HELLO, WORLD!
   }
   ```

5. **Run the example:**

   ```sh
   go run main.go
   ```

## Integrating Plugins with Templater

To use your plugin with Templater, you need to register it in your application. Here's how:

1. **Import your plugin:**

   ```go
   import (
       "github.com/singoesdeep/templater/plugins/uppercase"
   )
   ```

2. **Register the plugin:**

   ```go
   func main() {
       plugin := uppercase.New()
       templater.RegisterPlugin(plugin)
   }
   ```

3. **Use the plugin in your templates:**

   ```go
   // In your template
   {{ uppercase "hello, world!" }}
   ```

## Best Practices

- **Keep plugins simple and focused:** Each plugin should do one thing well.
- **Handle errors gracefully:** Always return meaningful error messages.
- **Test your plugins:** Write unit tests to ensure your plugin works as expected.
- **Document your plugins:** Provide clear documentation and examples.

## Conclusion

Plugins are a powerful way to extend Templater's functionality. By following the `templater.Plugin` interface, you can create custom plugins that integrate seamlessly with Templater.

For more examples and advanced usage, check out the [examples/plugins](examples/plugins) directory in the Templater repository. 