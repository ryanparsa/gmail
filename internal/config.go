package internal

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
)

// GetConfig loads the OAuth2 configuration from a credentials file and requested scopes.
//
// Parameters:
// - credentials: The file path to the JSON credentials file.
// - scopes: A list of OAuth2 scopes required for the application.
//
// Returns:
// - A pointer to an oauth2.Config object for authenticating requests.
// - If any error occurs (e.g., reading the credentials file or parsing JSON), the function logs the error and exits the application.
func GetConfig(credentials string, scopes ...string) *oauth2.Config {
	logrus.Infof("Reading OAuth2 credentials from file: %s", credentials)

	// Read the credentials file
	credentialsFile, err := os.ReadFile(credentials)
	if err != nil {
		logrus.Fatalf("Failed to read credentials file: %s, error: %v", credentials, err)
	}

	// Parse the credentials into an OAuth2 config
	authConfig, err := google.ConfigFromJSON(credentialsFile, scopes...)
	if err != nil {
		logrus.Fatalf("Failed to parse credentials JSON file: %s, error: %v", credentials, err)
	}

	logrus.Info("OAuth2 configuration successfully loaded")
	return authConfig
}
