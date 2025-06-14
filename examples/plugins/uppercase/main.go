package uppercase

import (
	"fmt"
)

func main() {
	// Example template
	template := "Hello, {{.name}}!"

	// Process the template
	result, err := PluginInstance.(*UppercasePlugin).ProcessTemplate(template)
	if err != nil {
		fmt.Printf("Error processing template: %v\n", err)
		return
	}

	fmt.Printf("Original template: %s\n", template)
	fmt.Printf("Processed result: %s\n", result)
}
