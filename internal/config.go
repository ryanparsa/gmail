package internal

import (
	"fmt"
	"google.golang.org/api/gmail/v1"
	"gopkg.in/yaml.v3"
	"os"
)

type Filters []*gmail.Filter
type Labels []*gmail.Label

type Messages []*gmail.Message

type Config struct {
	Labels  Labels  `yaml:"labels" json:"filters" jsonschema_description:"Filters to be applied to emails"`
	Filters Filters `yaml:"filters" json:"labels" jsonschema_description:"Labels to be created"`
}

func NewConfig(f Filters, l Labels) *Config {
	return &Config{
		Labels:  l,
		Filters: f,
	}

}

func NewConfigFromYAML(configFile string) (*Config, error) {
	fileExists(configFile)

	// Read the contents of the file
	cf, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(cf, &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %v", err)
	}
	return &c, nil
}

func (c *Config) SaveToFile(outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Marshal the data to YAML and write it to the file
	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("failed to encode data to YAML: %v", err)
	}

	return nil
}
