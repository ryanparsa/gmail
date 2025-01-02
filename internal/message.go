package internal

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/gmail/v1"
)

// Messages fetches Gmail messages with optional query and limit.
//
// Parameters:
// - max: The maximum number of messages to fetch.
// - query: A Gmail search query to filter messages.
//
// Returns:
// - A slice of Gmail messages.
// - An error if the API request fails.
func (s *Service) Messages(max int64, query string) (Messages, error) {
	logrus.Infof("Fetching Gmail messages with query: '%s' and max results: %d", query, max)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var messages Messages

	// Initialize the request
	req := s.Users.Messages.List("me").Q(query).MaxResults(max)

	// Iterate through the pages of messages
	err := req.Pages(ctx, func(page *gmail.ListMessagesResponse) error {
		for _, m := range page.Messages {
			// Fetch full message details
			msg, err := s.Users.Messages.Get("me", m.Id).Format("full").Do()
			if err != nil {
				logrus.Errorf("Error fetching message %s: %v", m.Id, err)
				continue
			}
			messages = append(messages, msg)
			logrus.Debugf("Fetched message ID: %s", m.Id)
		}
		return nil
	})

	if err != nil {
		logrus.Errorf("Failed to fetch messages: %v", err)
		return nil, err
	}

	logrus.Infof("Successfully fetched %d messages.", len(messages))
	return messages, nil
}
