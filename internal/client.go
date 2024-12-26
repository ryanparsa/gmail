package internal

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// getClient initializes an HTTP client authenticated using OAuth2 tokens.
// It attempts to read the token from the specified token file. If the token is not found or invalid,
// it fetches a new token using the provided OAuth2 configuration, saves it, and returns the HTTP client.
//
// Parameters:
// - authConfig: The OAuth2 configuration used to authenticate the client.
// - tokenFile: The file path where the OAuth2 token is stored.
//
// Returns:
// - A configured HTTP client that uses the OAuth2 token for requests.
func getClient(authConfig *oauth2.Config, tokenFile string) *http.Client {
	logrus.Info("Initializing OAuth2 client...")

	// Attempt to read the token from the file
	logrus.Infof("Reading token from file: %s", tokenFile)
	token, err := ReadToken(tokenFile)
	if err != nil {
		logrus.Warnf("Failed to read token from file: %s, error: %v", tokenFile, err)

		// Fetch a new token if reading fails
		logrus.Info("Fetching a new OAuth2 token...")
		token = GetToken(authConfig)

		// Save the new token to the file
		if err = SaveToken(tokenFile, token); err != nil {
			logrus.Errorf("Failed to save token to file: %s, error: %v", tokenFile, err)
		} else {
			logrus.Infof("Token successfully saved to file: %s", tokenFile)
		}
	} else {
		logrus.Infof("Token successfully read from file: %s", tokenFile)
	}

	// Return the authenticated HTTP client
	logrus.Info("OAuth2 client successfully initialized")
	return authConfig.Client(context.Background(), token)
}
