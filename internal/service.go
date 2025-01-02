package internal

import (
	"context"
	"fmt"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Service struct {
	*gmail.Service
}

func NewService(credentialsPath, tokenPath string, scopes []string) (*Service, error) {
	// Authenticate and get an HTTP client
	client, err := getGmailClient(credentialsPath, tokenPath, scopes)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Gmail client: %v", err)
	}

	// Create Gmail service
	svc, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gmail service: %v", err)
	}

	return &Service{svc}, nil

}
