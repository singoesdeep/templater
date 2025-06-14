package security

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// DangerousCommands is a list of potentially dangerous commands that should be prevented
var DangerousCommands = []string{
	"rm", "del", "delete", "format", "mkfs", "dd",
	"shutdown", "reboot", "halt", "poweroff",
	"chmod", "chown", "chattr",
	"wget", "curl", "nc", "netcat",
	"bash", "sh", "zsh", "powershell",
	"sudo", "su", "doas",
}

// DangerousPatterns are regex patterns that match potentially dangerous operations
var DangerousPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(exec|system|eval|spawn|fork)`),
	regexp.MustCompile(`(?i)(file|directory|path)\.(delete|remove|unlink)`),
	regexp.MustCompile(`(?i)(os|process)\.(exec|command|run)`),
	regexp.MustCompile(`(?i)(shell|bash|sh|zsh|powershell)`),
	regexp.MustCompile(`(?i)(sudo|su|doas)`),
	regexp.MustCompile(`(?i)(wget|curl|nc|netcat)`),
}

// DangerousPathPatterns are regex patterns that match potentially dangerous paths
var DangerousPathPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(/etc|/var|/usr|/bin|/sbin|/dev|/proc|/sys)`),
	regexp.MustCompile(`(?i)(\.\.|/\.\.)`),
	regexp.MustCompile(`(?i)(/root|/home/[^/]+/\.)`),
	regexp.MustCompile(`(?i)(/tmp/[^/]+/\.\.)`),
}

// ValidateTemplateContent checks if the template contains any dangerous operations
func ValidateTemplateContent(content string) error {
	// Check for dangerous commands
	words := strings.Fields(content)
	for _, word := range words {
		for _, cmd := range DangerousCommands {
			if strings.EqualFold(word, cmd) {
				return fmt.Errorf("dangerous command detected: %s", cmd)
			}
		}
	}

	// Check for dangerous patterns
	for _, pattern := range DangerousPatterns {
		if pattern.MatchString(content) {
			return fmt.Errorf("dangerous pattern detected: %s", pattern.String())
		}
	}

	return nil
}

// SanitizeData removes potentially dangerous content from data values
func SanitizeData(data map[string]string) map[string]string {
	sanitized := make(map[string]string)
	for k, v := range data {
		// Remove any shell command injection attempts
		v = strings.ReplaceAll(v, "`", "")
		v = strings.ReplaceAll(v, "$(", "")
		v = strings.ReplaceAll(v, "&&", "")
		v = strings.ReplaceAll(v, "||", "")
		v = strings.ReplaceAll(v, ";", "")
		sanitized[k] = v
	}
	return sanitized
}

// ValidateOutputPath ensures the output path is within allowed directories
func ValidateOutputPath(path string, allowedDirs []string) error {
	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Check for dangerous path patterns
	for _, pattern := range DangerousPathPatterns {
		if pattern.MatchString(absPath) {
			return fmt.Errorf("dangerous path pattern detected: %s", pattern.String())
		}
	}

	// Check if path is within allowed directories
	for _, dir := range allowedDirs {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			continue
		}
		if strings.HasPrefix(absPath, absDir) {
			// Verify directory exists and is accessible
			if _, err := os.Stat(absDir); err != nil {
				return fmt.Errorf("allowed directory not accessible: %w", err)
			}
			return nil
		}
	}

	return fmt.Errorf("output path %s is not within allowed directories", path)
}

// SanitizePath removes potentially dangerous elements from a path
func SanitizePath(path string) string {
	// Remove any path traversal attempts
	path = filepath.Clean(path)
	path = strings.ReplaceAll(path, "..", "")
	path = strings.ReplaceAll(path, "./", "")
	path = strings.ReplaceAll(path, "/.", "")

	// Remove any dangerous characters
	path = strings.Map(func(r rune) rune {
		if r < 32 || r == 127 {
			return -1
		}
		return r
	}, path)

	return path
}
