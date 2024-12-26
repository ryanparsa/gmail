package internal

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Service struct {
	*gmail.Service
}

// NewService creates a new Gmail service using the provided credentials, token, and scopes.
//
// Parameters:
// - credentials: Path to the credentials file.
// - token: Path to the token file for OAuth2 authentication.
// - scopes: List of scopes for Gmail API access.
//
// Returns:
// - A pointer to the initialized Service object.
// - Logs a fatal error and exits the program if initialization fails.
func NewService(credentials string, token string, scopes []string) *Service {
	logrus.Info("Initializing Gmail service")

	// Get OAuth2 config
	authConfig := GetConfig(credentials, scopes...)
	logrus.Info("OAuth2 configuration successfully loaded")

	// Get authenticated HTTP client
	client := getClient(authConfig, token)
	logrus.Info("OAuth2 client successfully initialized")

	// Create Gmail service
	service, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		logrus.Fatalf("Failed to initialize Gmail service: %v", err)
	}

	logrus.Info("Gmail service successfully initialized")
	return &Service{service}
}

// FetchEmails fetches emails from the Gmail API based on the provided query.
//
// Parameters:
// - query: The Gmail query string to filter emails.
//
// Returns:
// - A slice of Gmail message objects.
// - An error if the API request fails.
func (s *Service) FetchEmails(query string) ([]*gmail.Message, error) {
	logrus.Infof("Fetching emails with query: %s", query)

	var messages []*gmail.Message
	pageToken := ""

	for {
		// List messages with the specified query and page token
		req := s.Users.Messages.List("me").Q(query).PageToken(pageToken)
		resp, err := req.Do()
		if err != nil {
			logrus.Errorf("Failed to fetch emails: %v", err)
			return nil, err
		}

		// Append fetched messages
		messages = append(messages, resp.Messages...)
		logrus.Infof("Fetched %d messages so far", len(messages))

		// Break if there are no more pages
		if resp.NextPageToken == "" {
			break
		}

		// Update page token for the next iteration
		pageToken = resp.NextPageToken
	}

	logrus.Infof("Successfully fetched %d messages", len(messages))
	return messages, nil
}

// Wait pauses execution for the specified duration.
//
// Parameters:
// - duration: Duration to wait in seconds.
// - message: Message to log while waiting.
func Wait(duration int, message string) {
	if message == "" {
		logrus.Infof("Waiting for %d seconds", duration)
	} else {
		logrus.Infof("Waiting for %d seconds: %s", duration, message)
	}
	time.Sleep(time.Duration(duration) * time.Second)
	logrus.Info("Wait complete")
}
