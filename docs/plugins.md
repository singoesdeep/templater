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