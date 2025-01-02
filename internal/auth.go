package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"os"
	"regexp"
)

// Authenticate handles Gmail API authentication
func Authenticate(credentialsPath, tokenPath string, scopes []string) error {
	// Check if credentials.json exists
	if !fileExists(credentialsPath) {
		logrus.Error("credentials.json not found. Please download it from Google Cloud Console")
		return errors.New("credentials.json not found. Please download it from Google Cloud Console")
	}

	// Load OAuth2 config
	config, err := loadCredentials(credentialsPath, scopes)
	if err != nil {
		logrus.Errorf("Failed to load credentials: %v", err)
		return fmt.Errorf("failed to load credentials: %v", err)
	}

	// Get or refresh token
	token, err := getToken(config, tokenPath)
	if err != nil {
		logrus.Warnf("Failed to authenticate with saved token: %v", err)
		logrus.Info("Attempting to authenticate via web...")
		token, err = getTokenFromWeb(config)
		if err != nil {
			logrus.Errorf("Failed to authenticate via web: %v", err)
			return fmt.Errorf("failed to authenticate: %v", err)
		}
	}

	// Validate token with a simple API call
	err = validateToken(config, token)
	if err != nil {
		logrus.Errorf("Failed to validate token: %v", err)
		return fmt.Errorf("failed to validate token: %v", err)
	}

	// Save the token
	err = saveToken(tokenPath, token)
	if err != nil {
		logrus.Errorf("Failed to save token: %v", err)
		return fmt.Errorf("failed to save token: %v", err)
	}

	logrus.Info("Authentication successful.")
	return nil
}

// loadCredentials loads OAuth2 config from the provided credentials file and custom scopes
func loadCredentials(credentialsPath string, scopes []string) (*oauth2.Config, error) {
	// Read the contents of the file
	credentialsData, err := os.ReadFile(credentialsPath)
	if err != nil {
		logrus.Errorf("Failed to read credentials file: %v", err)
		return nil, err
	}

	// Parse the JSON and create an OAuth2 config with custom scopes
	config, err := google.ConfigFromJSON(credentialsData, scopes...)
	if err != nil {
		logrus.Errorf("Failed to parse credentials JSON: %v", err)
		return nil, err
	}

	logrus.Info("OAuth2 configuration loaded successfully.")
	return config, nil
}

// getRedirectUrl validates and returns the RedirectURL from the OAuth2 config
func getRedirectUrl(config *oauth2.Config) (string, error) {
	// Ensure the RedirectURL is not empty
	if config.RedirectURL == "" {
		logrus.Error("No redirect URI found in credentials.")
		return "", fmt.Errorf("no redirect URI found in credentials")
	}

	// Regex to validate localhost URI with a port (e.g., http://localhost:8888)
	regex := `^http://localhost:\d+$`
	matched, err := regexp.MatchString(regex, config.RedirectURL)
	if err != nil {
		logrus.Errorf("Error validating redirect URI: %v", err)
		return "", fmt.Errorf("error validating redirect URI: %v", err)
	}

	// If the redirect URL doesn't match the regex, return an error
	if !matched {
		logrus.Errorf("Invalid redirect URI found: %s", config.RedirectURL)
		return "", fmt.Errorf("no valid localhost redirect URI found in credentials. Found: %s", config.RedirectURL)
	}

	logrus.Infof("Redirect URI validated successfully: %s", config.RedirectURL)
	return config.RedirectURL, nil
}

// getGmailClient initializes a Gmail API client using OAuth2 credentials and token
func getGmailClient(credentialsPath, tokenPath string, scopes []string) (*http.Client, error) {
	// Load OAuth2 config
	config, err := loadCredentials(credentialsPath, scopes)
	if err != nil {
		logrus.Errorf("Failed to load credentials: %v", err)
		return nil, fmt.Errorf("failed to load credentials: %v", err)
	}

	// Get or refresh token
	token, err := getToken(config, tokenPath)
	if err != nil {
		logrus.Errorf("Failed to get token: %v", err)
		return nil, fmt.Errorf("failed to get token: %v", err)
	}

	// Validate token with a simple API call
	err = validateToken(config, token)
	if err != nil {
		logrus.Errorf("Failed to validate token: %v", err)
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}

	logrus.Info("Gmail API client initialized successfully.")
	return config.Client(context.Background(), token), nil
}
