package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
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

// getToken retrieves or generates a new token
func getToken(config *oauth2.Config, tokenPath string) (*oauth2.Token, error) {

	token, err := loadToken(tokenPath)
	if err != nil {
		return nil, err
	}

	if token.Valid() {
		return token, nil
	}

	token, err = refreshToken(config, token)
	if err != nil {
		return nil, err
	}
	err = saveToken(tokenPath, token)
	if err != nil {
		return nil, err
	}
	return token, nil

}

func refreshToken(config *oauth2.Config, token *oauth2.Token) (*oauth2.Token, error) {

	// If token is expired, refresh it
	fmt.Println("Token is expired, attempting to refresh.")
	tokenSource := config.TokenSource(context.Background(), token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %v", err)
	}

	return newToken, nil
}

// loadToken reads the token from token.json
func loadToken(tokenPath string) (*oauth2.Token, error) {
	if !fileExists(tokenPath) {
		return nil, fmt.Errorf("no token file found")
	}

	file, err := os.Open(tokenPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var token oauth2.Token
	err = json.NewDecoder(file).Decode(&token)
	return &token, err
}

// saveToken saves the token to token.json
func saveToken(tokenPath string, token *oauth2.Token) error {
	file, err := os.Create(tokenPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(token)
}

// validateToken makes a simple API call to validate the token
func validateToken(config *oauth2.Config, token *oauth2.Token) error {
	ctx := context.Background()
	client := config.Client(ctx, token)
	svc, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	// Validate by fetching the user's Gmail profile
	_, err = svc.Users.GetProfile("me").Do()
	return err
}
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	// Get the RedirectURL from the OAuth2 config
	redirectURL, err := getRedirectUrl(config)
	if err != nil {
		return nil, fmt.Errorf("failed to get redirect URL: %v", err)
	}

	// Extract the port from the RedirectURL
	port := strings.Split(redirectURL, ":")[2]

	// Channel to receive the authorization code
	authCodeChan := make(chan string)

	// Start a temporary HTTP server
	server := &http.Server{Addr: ":" + port}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Extract the authorization code from the URL
		authCode := r.URL.Query().Get("code")
		if authCode == "" {
			http.Error(w, "Authorization code not found", http.StatusBadRequest)
			return
		}

		// Notify the CLI and user
		fmt.Fprintln(w, "Authentication successful! You can close this browser window.")
		authCodeChan <- authCode
	})

	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting HTTP server: %v\n", err)
		}
	}()

	// Generate the authorization URL
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser to authorize the application:\n%v\n", authURL)

	// Wait for the authorization code from the channel
	authCode := <-authCodeChan

	// Shutdown the server
	go func() {
		_ = server.Close()
	}()

	// Exchange the authorization code for a token
	token, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange authorization code for token: %v", err)
	}

	return token, nil
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

// ClearToken deletes the token.json file
func ClearToken(tokenPath string) error {
	return os.Remove(tokenPath)
}

// fileExists checks if a file exists at the given path
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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
