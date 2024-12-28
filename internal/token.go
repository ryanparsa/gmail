package internal

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// ReadToken reads an OAuth2 token from a file.
//
// Parameters:
// - file: The path to the file containing the token.
//
// Returns:
// - The OAuth2 token.
// - An error if the file cannot be read or the token cannot be decoded.
func ReadToken(file string) (*oauth2.Token, error) {
	logrus.Infof("Reading token from file: %s", file)

	f, err := os.Open(file)
	if err != nil {
		logrus.Errorf("Failed to open token file: %s, error: %v", file, err)
		return nil, err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			logrus.Errorf("Failed to close token file: %s, error: %v", file, cerr)
		}
	}()

	var token oauth2.Token
	if err = json.NewDecoder(f).Decode(&token); err != nil {
		logrus.Errorf("Failed to decode token from file: %s, error: %v", file, err)
		return nil, err
	}

	logrus.Infof("Successfully read token from file: %s", file)
	return &token, nil
}

// SaveToken saves an OAuth2 token to a file.
//
// Parameters:
// - path: The path to save the token.
// - token: The OAuth2 token to be saved.
//
// Returns:
// - An error if the token cannot be written to the file.
func SaveToken(path string, token *oauth2.Token) error {
	logrus.Infof("Saving token to file: %s", path)

	file, err := os.Create(path)
	if err != nil {
		logrus.Errorf("Failed to create token file: %s, error: %v", path, err)
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			logrus.Errorf("Failed to close token file: %s, error: %v", path, cerr)
		}
	}()

	if err = json.NewEncoder(file).Encode(token); err != nil {
		logrus.Errorf("Failed to encode token to file: %s, error: %v", path, err)
		return err
	}

	logrus.Infof("Successfully saved token to file: %s", path)
	return nil
}

// GetToken initiates the OAuth2 authentication flow to retrieve an OAuth2 token.
//
// Parameters:
// - authConfig: The OAuth2 configuration object.
//
// Returns:
// - The OAuth2 token.
func GetToken(authConfig *oauth2.Config) *oauth2.Token {
	logrus.Info("Starting OAuth2 authentication flow")

	server := &http.Server{Addr: strings.TrimPrefix(authConfig.RedirectURL, "http://")}
	codeChan := make(chan string)

	authURL := authConfig.AuthCodeURL(generateState(), oauth2.AccessTypeOffline)
	logrus.Infof("Visit the following URL for authentication: %s", authURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		fmt.Fprintf(w, "Authentication successful. You can now close this tab.")
		codeChan <- code
	})

	// Start the HTTP server in a goroutine
	go func() {
		logrus.Infof("Starting local HTTP server on: %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Errorf("HTTP server error: %v", err)
		}
	}()

	// Wait for the auth code
	code := <-codeChan
	logrus.Info("Received authorization code")

	// Shut down the server
	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Failed to shut down HTTP server: %v", err)
	}

	// Exchange the authorization code for a token
	token, err := authConfig.Exchange(context.Background(), code)
	if err != nil {
		logrus.Fatalf("Failed to exchange authorization code for token: %v", err)
	}

	logrus.Info("Successfully obtained OAuth2 token")
	return token
}

// generateState generates a random state string for OAuth2 authentication flow.
// This helps prevent CSRF attacks by verifying the state parameter in the callback.
func generateState() string {
	// Define the length of the random state string (32 bytes in this example).
	const stateLength = 32

	// Generate a random byte array.
	stateBytes := make([]byte, stateLength)
	if _, err := rand.Read(stateBytes); err != nil {
		logrus.Fatalf("Failed to generate random state: %v", err)
	}

	// Encode the random bytes to a base64 string for safe usage in URLs.
	return base64.URLEncoding.EncodeToString(stateBytes)
}
