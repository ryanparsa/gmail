package internal

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Templates defines a map of template names to their file paths.
var Templates = map[string]string{
	"config.yaml":      "internal/templates/config.yaml",
	"credentials.json": "internal/templates/credentials.json",
	"token.json":       "internal/templates/token.json",
}

// LoadTemplate reads the content of a template file at the given path and returns it as a string.
//
// Parameters:
// - path: The file path of the template.
//
// Returns:
// - The content of the file as a string.
// - Logs an error if the file cannot be read and exits the program.
func LoadTemplate(path string) string {
	logrus.Infof("Loading template from path: %s", path)

	data, err := os.ReadFile(path)
	if err != nil {
		logrus.Fatalf("Failed to read template from path: %s, error: %v", path, err)
	}

	logrus.Infof("Successfully loaded template from path: %s", path)
	return string(data)
}

// SaveTemplate saves the provided template content to the specified file path.
//
// Parameters:
// - path: The file path where the template content should be saved.
//
// Returns:
// - An error if the file cannot be created or written to.
func SaveTemplate(fileName, templateFilePath string) error {
	logrus.Infof("Saving template to path: %s", fileName)

	// Create or overwrite the file
	file, err := os.Create(fileName)
	if err != nil {
		logrus.Errorf("Failed to create file at path: %s, error: %v", fileName, err)
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			logrus.Errorf("Failed to close file at path: %s, error: %v", fileName, cerr)
		}
	}()

	// Write the template content to the file

	_, err = file.WriteString(LoadTemplate(templateFilePath))
	if err != nil {
		logrus.Errorf("Failed to write to file at path: %s, error: %v", fileName, err)
		return err
	}

	logrus.Infof("Successfully saved template to path: %s", fileName)
	return nil
}
