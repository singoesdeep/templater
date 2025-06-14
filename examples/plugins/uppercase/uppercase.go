package uppercase

import (
	"strings"

	"github.com/singoesdeep/templater/pkg/templater/plugin"
)

// UppercasePlugin is a simple plugin that converts template output to uppercase
type UppercasePlugin struct {
	*plugin.BasePlugin
}

// NewUppercasePlugin creates a new UppercasePlugin
func NewUppercasePlugin() *UppercasePlugin {
	return &UppercasePlugin{
		BasePlugin: &plugin.BasePlugin{
			PluginName:        "uppercase",
			PluginVersion:     "1.0.0",
			PluginDescription: "Converts template output to uppercase",
		},
	}
}

// ProcessTemplate processes the template content
func (p *UppercasePlugin) ProcessTemplate(content string) (string, error) {
	return strings.ToUpper(content), nil
}

// GetTemplateFuncs returns additional template functions
func (p *UppercasePlugin) GetTemplateFuncs() map[string]interface{} {
	return map[string]interface{}{
		"uppercase": strings.ToUpper,
	}
}

// Initialize is called when the plugin is loaded
func (p *UppercasePlugin) Initialize() error {
	return nil
}

// Shutdown is called when the plugin is unloaded
func (p *UppercasePlugin) Shutdown() error {
	return nil
}

// Export PluginInstance for dynamic loading
var PluginInstance plugin.Plugin = NewUppercasePlugin()
