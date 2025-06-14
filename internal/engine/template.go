package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"text/template"
	"time"
	"unicode"

	"github.com/singoesdeep/templater/internal/reliability"
	"github.com/singoesdeep/templater/internal/security"
	"gopkg.in/yaml.v3"
)

// titleCase converts a string to title case
func titleCase(s string) string {
	if s == "" {
		return s
	}
	prev := ' '
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(prev) {
			prev = r
			return unicode.ToTitle(r)
		}
		prev = r
		return unicode.ToLower(r)
	}, s)
}

// TemplateFuncs provides additional functions for templates
var TemplateFuncs = template.FuncMap{
	"now": func() string {
		return time.Now().Format(time.RFC3339)
	},
	"upper": strings.ToUpper,
	"lower": strings.ToLower,
	"title": titleCase,
	"join":  strings.Join,
}

// TemplateMetadata represents metadata about a template
type TemplateMetadata struct {
	Name         string   `yaml:"name"`
	Description  string   `yaml:"description"`
	Author       string   `yaml:"author"`
	Version      string   `yaml:"version"`
	DependsOn    []string `yaml:"depends_on"`
	RequiredKeys []string `yaml:"required_keys"`
}

// ParseTemplateMetadata extracts metadata from template comments
func ParseTemplateMetadata(content string) (*TemplateMetadata, error) {
	// Look for metadata in comments at the start of the file
	metadataStart := "{{/*"
	metadataEnd := "*/}}"

	startIdx := strings.Index(content, metadataStart)
	if startIdx == -1 {
		return nil, nil // No metadata found
	}

	endIdx := strings.Index(content[startIdx:], metadataEnd)
	if endIdx == -1 {
		return nil, fmt.Errorf("unclosed metadata block")
	}

	metadataContent := content[startIdx+len(metadataStart) : startIdx+endIdx]
	metadata := &TemplateMetadata{}

	if err := yaml.Unmarshal([]byte(metadataContent), metadata); err != nil {
		return nil, fmt.Errorf("error parsing template metadata: %w", err)
	}

	return metadata, nil
}

// GetTemplateDependencies returns a list of template dependencies
func GetTemplateDependencies(templatePath string) ([]string, error) {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("error reading template file: %w", err)
	}

	metadata, err := ParseTemplateMetadata(string(content))
	if err != nil {
		return nil, err
	}

	if metadata == nil {
		return nil, nil
	}

	return metadata.DependsOn, nil
}

// ValidateTemplateDependencies checks if all required template dependencies exist
func ValidateTemplateDependencies(templatePath string) error {
	deps, err := GetTemplateDependencies(templatePath)
	if err != nil {
		return err
	}

	for _, dep := range deps {
		if _, err := os.Stat(dep); os.IsNotExist(err) {
			return fmt.Errorf("missing template dependency: %s", dep)
		}
	}

	return nil
}

var (
	// templateCache stores parsed templates
	templateCache = struct {
		sync.RWMutex
		templates map[string]*templateInfo
	}{
		templates: make(map[string]*templateInfo),
	}

	// cacheExpiry is the duration after which cached templates expire
	cacheExpiry = 5 * time.Minute

	// maxCacheSize limits the number of cached templates
	maxCacheSize = 100

	// fileContentCache stores file contents to reduce I/O
	fileContentCache = struct {
		sync.RWMutex
		contents map[string]*fileContent
	}{
		contents: make(map[string]*fileContent),
	}
)

// templateInfo stores template and its metadata
type templateInfo struct {
	template     *template.Template
	lastUsed     time.Time
	lastModified time.Time
}

// fileContent stores file content and metadata
type fileContent struct {
	content     []byte
	lastUsed    time.Time
	lastMod     time.Time
	accessCount int
}

// getFileContent retrieves file content from cache or reads from disk
func getFileContent(path string) ([]byte, error) {
	// Check cache first
	fileContentCache.RLock()
	if content, exists := fileContentCache.contents[path]; exists {
		// Check if file is still valid
		fileInfo, err := os.Stat(path)
		if err == nil && fileInfo.ModTime().Equal(content.lastMod) {
			content.lastUsed = time.Now()
			content.accessCount++
			fileContentCache.RUnlock()
			return content.content, nil
		}
	}
	fileContentCache.RUnlock()

	// Read file
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Get file info
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("error getting file info: %w", err)
	}

	// Update cache
	fileContentCache.Lock()
	// Remove least accessed content if cache is full
	if len(fileContentCache.contents) >= maxCacheSize {
		var leastAccessed *fileContent
		var leastPath string
		for p, c := range fileContentCache.contents {
			if leastAccessed == nil || c.accessCount < leastAccessed.accessCount {
				leastAccessed = c
				leastPath = p
			}
		}
		if leastAccessed != nil {
			delete(fileContentCache.contents, leastPath)
		}
	}
	fileContentCache.contents[path] = &fileContent{
		content:     content,
		lastUsed:    time.Now(),
		lastMod:     fileInfo.ModTime(),
		accessCount: 1,
	}
	fileContentCache.Unlock()

	return content, nil
}

// GetCachedTemplate retrieves a template from cache or parses a new one
func GetCachedTemplate(templatePath string) (*template.Template, error) {
	// Check cache first
	templateCache.RLock()
	if info, exists := templateCache.templates[templatePath]; exists {
		// Check if template is still valid
		if time.Since(info.lastUsed) < cacheExpiry {
			// Update last used time
			info.lastUsed = time.Now()
			templateCache.RUnlock()
			return info.template, nil
		}
	}
	templateCache.RUnlock()

	// Get template content from cache or file
	tmplContent, err := getFileContent(templatePath)
	if err != nil {
		return nil, err
	}

	// Get file modification time
	fileInfo, err := os.Stat(templatePath)
	if err != nil {
		return nil, fmt.Errorf("error getting file info: %w", err)
	}

	// Parse template with functions
	tmpl, err := template.New(filepath.Base(templatePath)).
		Funcs(TemplateFuncs).
		Parse(string(tmplContent))
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %w", err)
	}

	// Update cache
	templateCache.Lock()
	// Remove oldest template if cache is full
	if len(templateCache.templates) >= maxCacheSize {
		var oldest *templateInfo
		var oldestPath string
		for p, i := range templateCache.templates {
			if oldest == nil || i.lastUsed.Before(oldest.lastUsed) {
				oldest = i
				oldestPath = p
			}
		}
		if oldest != nil {
			delete(templateCache.templates, oldestPath)
		}
	}
	templateCache.templates[templatePath] = &templateInfo{
		template:     tmpl,
		lastUsed:     time.Now(),
		lastModified: fileInfo.ModTime(),
	}
	templateCache.Unlock()

	return tmpl, nil
}

// CleanupCache removes expired templates and file contents
func CleanupCache() {
	now := time.Now()

	// Cleanup template cache
	templateCache.Lock()
	for path, info := range templateCache.templates {
		if now.Sub(info.lastUsed) > cacheExpiry {
			delete(templateCache.templates, path)
		}
	}
	templateCache.Unlock()

	// Cleanup file content cache
	fileContentCache.Lock()
	for path, content := range fileContentCache.contents {
		if now.Sub(content.lastUsed) > cacheExpiry {
			delete(fileContentCache.contents, path)
		}
	}
	fileContentCache.Unlock()

	// Trigger garbage collection
	runtime.GC()
}

// RenderTemplate processes a template file with the given data
func RenderTemplate(templatePath string, data map[string]string) (string, error) {
	// Get template content from cache or file
	tmplContent, err := getFileContent(templatePath)
	if err != nil {
		return "", err
	}

	// Validate template content
	if err := security.ValidateTemplateContent(string(tmplContent)); err != nil {
		return "", fmt.Errorf("template validation error: %w", err)
	}

	// Sanitize data
	data = security.SanitizeData(data)

	// Get template from cache or parse new
	tmpl, err := GetCachedTemplate(templatePath)
	if err != nil {
		return "", err
	}

	// Execute template
	var result bytes.Buffer
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return result.String(), nil
}

// LoadData loads data from a JSON or YAML file
func LoadData(dataPath string) (map[string]string, error) {
	if dataPath == "" {
		return make(map[string]string), nil
	}

	data, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, fmt.Errorf("error reading data file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(dataPath))
	var result map[string]string
	if ext == ".yaml" || ext == ".yml" {
		if err := yaml.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("error parsing YAML data: %w", err)
		}
	} else {
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("error parsing JSON data: %w", err)
		}
	}

	return result, nil
}

// WriteToFile writes the content to a file with security and reliability checks
func WriteToFile(filePath string, content string) error {
	// Validate output path
	allowedDirs := []string{
		filepath.Dir(filePath),
	}
	if err := security.ValidateOutputPath(filePath, allowedDirs); err != nil {
		return fmt.Errorf("security error: %w", err)
	}

	// Create backup if file exists
	if err := reliability.BackupFile(filePath); err != nil {
		return fmt.Errorf("backup error: %w", err)
	}

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		// Attempt to restore from backup on error
		if restoreErr := reliability.RestoreFromBackup(filePath); restoreErr != nil {
			return fmt.Errorf("write error: %w (restore failed: %v)", err, restoreErr)
		}
		return fmt.Errorf("write error: %w (restored from backup)", err)
	}

	// Verify file integrity
	if err := reliability.VerifyFileIntegrity(filePath, []byte(content)); err != nil {
		// Attempt to restore from backup on integrity check failure
		if restoreErr := reliability.RestoreFromBackup(filePath); restoreErr != nil {
			return fmt.Errorf("integrity check failed: %w (restore failed: %v)", err, restoreErr)
		}
		return fmt.Errorf("integrity check failed: %w (restored from backup)", err)
	}

	// Cleanup old backups (keep last 7 days)
	if err := reliability.CleanupOldBackups(filePath, 7*24*time.Hour); err != nil {
		// Log but don't fail on cleanup error
		fmt.Printf("Warning: Failed to cleanup old backups: %v\n", err)
	}

	return nil
}

// ExtractTemplateKeys extracts all keys used in the template
func ExtractTemplateKeys(templatePath string) ([]string, error) {
	// Read template content directly for key extraction
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("error reading template file: %w", err)
	}

	// Simple key extraction (can be enhanced for more complex cases)
	keys := make(map[string]bool)
	parts := strings.Split(string(content), "{{.")
	for _, part := range parts[1:] {
		if idx := strings.Index(part, "}}"); idx != -1 {
			key := strings.TrimSpace(part[:idx])
			keys[key] = true
		}
	}

	result := make([]string, 0, len(keys))
	for key := range keys {
		result = append(result, key)
	}
	return result, nil
}

// ValidateTemplateData validates that all required template keys are present in the data
func ValidateTemplateData(templatePath string, data map[string]string) error {
	keys, err := ExtractTemplateKeys(templatePath)
	if err != nil {
		return fmt.Errorf("error extracting template keys: %w", err)
	}

	missing := []string{}
	for _, k := range keys {
		if _, ok := data[k]; !ok {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required data keys: %v", missing)
	}
	return nil
}

// ClearTemplateCache clears the template cache
func ClearTemplateCache() {
	templateCache = struct {
		sync.RWMutex
		templates map[string]*templateInfo
	}{
		templates: make(map[string]*templateInfo),
	}
}
