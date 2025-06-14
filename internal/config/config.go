package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the templater configuration
type Config struct {
	Defaults struct {
		OutputDir     string `yaml:"output_dir"`
		WatchInterval string `yaml:"watch_interval"`
		Backup        bool   `yaml:"backup"`
		Language      string `yaml:"language"`
	} `yaml:"defaults"`
}

// LoadConfig loads the configuration from .templater.yaml
func LoadConfig() (*Config, error) {
	// Look for config in current directory and parent directories
	configPath := findConfigFile()
	if configPath == "" {
		return &Config{}, nil // Return empty config if no file found
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return &config, nil
}

// findConfigFile looks for .templater.yaml in current and parent directories
func findConfigFile() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		configPath := filepath.Join(dir, ".templater.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // Reached root directory
		}
		dir = parent
	}

	return ""
}

// GetOutputDir returns the configured output directory or a default
func (c *Config) GetOutputDir() string {
	if c.Defaults.OutputDir != "" {
		return c.Defaults.OutputDir
	}
	return "generated"
}

// GetWatchInterval returns the configured watch interval or a default
func (c *Config) GetWatchInterval() string {
	if c.Defaults.WatchInterval != "" {
		return c.Defaults.WatchInterval
	}
	return "1s"
}

// ShouldBackup returns whether files should be backed up before overwriting
func (c *Config) ShouldBackup() bool {
	return c.Defaults.Backup
}

// GetLanguage returns the configured language or a default
func (c *Config) GetLanguage() string {
	if c.Defaults.Language != "" {
		return c.Defaults.Language
	}
	return "en"
}

// GetBackup returns whether backups should be created
func (c *Config) GetBackup() bool {
	if c == nil {
		return false
	}
	return c.Defaults.Backup
}
