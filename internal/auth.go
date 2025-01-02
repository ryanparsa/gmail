package internal

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"os"
	"regexp"
)

// Authenticate handles Gmail API authentication
func Authenticate(credentialsPath, tokenPath string, scopes []string) error {
	// Check if credentials.json exists
	if !fileExists(credentialsPath) {
		return errors.New("credentials.json not found. Please download it from Google Cloud Console")
	}

	// Load OAuth2 config
	config, err := loadCredentials(credentialsPath, scopes)
	if err != nil {
		return fmt.Errorf("failed to load credentials: %v", err)
	}

	// Get or refresh token
	token, err := getToken(config, tokenPath)
	if err != nil {
		log.Printf("failed to authenticate: %v\n", err)
		token, err = getTokenFromWeb(config)
		if err != nil {
			return fmt.Errorf("failed to authenticate: %v", err)
		}
	}

	// Validate token with a simple API call
	err = validateToken(config, token)
	if err != nil {
		return fmt.Errorf("failed to validate token: %v", err)
	}

	// Save the token
	err = saveToken(tokenPath, token)
	if err != nil {
		return fmt.Errorf("failed to save token: %v", err)
	}

	return nil
}

// loadCredentials loads OAuth2 config from the provided credentials file and custom scopes
func loadCredentials(credentialsPath string, scopes []string) (*oauth2.Config, error) {
	// Read the contents of the file
	credentialsData, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, err
	}

	// Parse the JSON and create an OAuth2 config with custom scopes
	config, err := google.ConfigFromJSON(credentialsData, scopes...)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// getRedirectUrl validates and returns the RedirectURL from the OAuth2 config
func getRedirectUrl(config *oauth2.Config) (string, error) {
	// Ensure the RedirectURL is not empty
	if config.RedirectURL == "" {
		return "", fmt.Errorf("no redirect URI found in credentials")
	}

	// Regex to validate localhost URI with a port (e.g., http://localhost:8888)
	regex := `^http://localhost:\d+$`
	matched, err := regexp.MatchString(regex, config.RedirectURL)
	if err != nil {
		return "", fmt.Errorf("error validating redirect URI: %v", err)
	}

	// If the redirect URL doesn't match the regex, return an error
	if !matched {
		return "", fmt.Errorf("no valid localhost redirect URI found in credentials. Found: %s", config.RedirectURL)
	}

	// Return the valid RedirectURL
	return config.RedirectURL, nil
}

func getGmailClient(credentialsPath, tokenPath string, scopes []string) (*http.Client, error) {

	// Load OAuth2 config
	config, err := loadCredentials(credentialsPath, scopes)
	if err != nil {
		return nil, fmt.Errorf("failed to load credentials: %v", err)
	}

	// Get or refresh token
	token, err := getToken(config, tokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %v", err)
	}

	// Validate token with a simple API call
	err = validateToken(config, token)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}

	return config.Client(context.Background(), token), nil

}
