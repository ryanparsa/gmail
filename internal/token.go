package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"net/http"
	"os"
	"strings"
)

// getToken retrieves or generates a new token
func getToken(config *oauth2.Config, tokenPath string) (*oauth2.Token, error) {
	token, err := loadToken(tokenPath)
	if err != nil {
		logrus.Errorf("Failed to load token: %v", err)
		return nil, err
	}

	if token.Valid() {
		logrus.Info("Token is valid.")
		return token, nil
	}

	logrus.Warn("Token is expired, attempting to refresh.")
	token, err = refreshToken(config, token)
	if err != nil {
		logrus.Errorf("Failed to refresh token: %v", err)
		return nil, err
	}

	err = saveToken(tokenPath, token)
	if err != nil {
		logrus.Errorf("Failed to save refreshed token: %v", err)
		return nil, err
	}

	logrus.Info("Token refreshed and saved successfully.")
	return token, nil
}

func refreshToken(config *oauth2.Config, token *oauth2.Token) (*oauth2.Token, error) {
	tokenSource := config.TokenSource(context.Background(), token)
	newToken, err := tokenSource.Token()
	if err != nil {
		logrus.Errorf("Failed to refresh token: %v", err)
		return nil, fmt.Errorf("failed to refresh token: %v", err)
	}

	logrus.Info("Token refreshed successfully.")
	return newToken, nil
}

// loadToken reads the token from token.json
func loadToken(tokenPath string) (*oauth2.Token, error) {
	if !fileExists(tokenPath) {
		logrus.Warn("Token file not found.")
		return nil, fmt.Errorf("no token file found")
	}

	file, err := os.Open(tokenPath)
	if err != nil {
		logrus.Errorf("Failed to open token file: %v", err)
		return nil, err
	}
	defer file.Close()

	var token oauth2.Token
	err = json.NewDecoder(file).Decode(&token)
	if err != nil {
		logrus.Errorf("Failed to decode token: %v", err)
		return nil, err
	}

	logrus.Info("Token loaded successfully.")
	return &token, nil
}

// saveToken saves the token to token.json
func saveToken(tokenPath string, token *oauth2.Token) error {
	file, err := os.Create(tokenPath)
	if err != nil {
		logrus.Errorf("Failed to create token file: %v", err)
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(token)
	if err != nil {
		logrus.Errorf("Failed to save token to file: %v", err)
		return err
	}

	logrus.Info("Token saved successfully.")
	return nil
}

// validateToken makes a simple API call to validate the token
func validateToken(config *oauth2.Config, token *oauth2.Token) error {
	ctx := context.Background()
	client := config.Client(ctx, token)
	svc, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		logrus.Errorf("Failed to create Gmail service: %v", err)
		return err
	}

	_, err = svc.Users.GetProfile("me").Do()
	if err != nil {
		logrus.Errorf("Failed to validate token: %v", err)
		return err
	}

	logrus.Info("Token validated successfully.")
	return nil
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	redirectURL, err := getRedirectUrl(config)
	if err != nil {
		logrus.Errorf("Failed to get redirect URL: %v", err)
		return nil, fmt.Errorf("failed to get redirect URL: %v", err)
	}

	port := strings.Split(redirectURL, ":")[2]
	authCodeChan := make(chan string)

	server := &http.Server{Addr: ":" + port}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		authCode := r.URL.Query().Get("code")
		if authCode == "" {
			http.Error(w, "Authorization code not found", http.StatusBadRequest)
			logrus.Warn("Authorization code not found in request.")
			return
		}

		fmt.Fprintln(w, "Authentication successful! You can close this browser window.")
		logrus.Info("Authorization code received.")
		authCodeChan <- authCode
	})

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("Error starting HTTP server: %v", err)
		}
	}()

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	logrus.Infof("Open the following URL in your browser to authenticate: %s", authURL)

	authCode := <-authCodeChan

	go func() {
		_ = server.Close()
	}()

	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		logrus.Errorf("Failed to exchange authorization code for token: %v", err)
		return nil, fmt.Errorf("failed to exchange authorization code for token: %v", err)
	}

	logrus.Info("Token obtained successfully from web authorization.")
	return token, nil
}

// ClearToken deletes the token.json file
func ClearToken(tokenPath string) error {
	err := os.Remove(tokenPath)
	if err != nil {
		logrus.Errorf("Failed to delete token file: %v", err)
		return err
	}

	logrus.Info("Token file deleted successfully.")
	return nil
}
