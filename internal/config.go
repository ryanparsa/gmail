package internal

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/gmail/v1"
	"gopkg.in/yaml.v3"
)

// Filters and Labels represent Gmail Filters and Labels
type Filters []*gmail.Filter
type Labels []*gmail.Label

// Messages represent Gmail Messages
type Messages []*gmail.Message

// Config represents the configuration containing filters and labels
type Config struct {
	Labels  Labels  `yaml:"labels" json:"labels" jsonschema_description:"Labels to be created"`
	Filters Filters `yaml:"filters" json:"filters" jsonschema_description:"Filters to be applied to emails"`
}

// NewConfig creates a new Config instance from filters and labels
func NewConfig(f Filters, l Labels) *Config {
	return &Config{
		Labels:  l,
		Filters: f,
	}
}

// NewConfigFromYAML loads a Config from a YAML file
func NewConfigFromYAML(configFile string) (*Config, error) {
	// Check if the file exists
	if !fileExists(configFile) {
		logrus.Errorf("Configuration file does not exist: %s", configFile)
		return nil, fmt.Errorf("configuration file does not exist: %s", configFile)
	}

	// Read the contents of the file
	data, err := os.ReadFile(configFile)
	if err != nil {
		logrus.Errorf("Failed to read configuration file: %v", err)
		return nil, fmt.Errorf("failed to read configuration file: %v", err)
	}

	// Unmarshal the YAML data into the Config struct
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		logrus.Errorf("Failed to unmarshal YAML data: %v", err)
		return nil, fmt.Errorf("failed to unmarshal YAML data: %v", err)
	}

	logrus.Infof("Configuration loaded successfully from file: %s", configFile)
	return &config, nil
}

// SaveToFile saves the Config to a specified file in YAML format
func (c *Config) SaveToFile(outputPath string) error {
	// Create or overwrite the file
	file, err := os.Create(outputPath)
	if err != nil {
		logrus.Errorf("Failed to create file: %v", err)
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Write the YAML-encoded data to the file
	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	if err := encoder.Encode(c); err != nil {
		logrus.Errorf("Failed to encode data to YAML: %v", err)
		return fmt.Errorf("failed to encode data to YAML: %v", err)
	}

	logrus.Infof("Configuration saved successfully to file: %s", outputPath)
	return nil
}

// Helper function to check if a file exists
func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		logrus.Warnf("File does not exist: %s", filePath)
		return false
	}
	return !info.IsDir()
}
