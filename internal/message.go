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
	var fetched int64

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var messages Messages

	// Initialize the request
	req := s.Users.Messages.List(userId).Q(query)

	// Iterate through the pages of messages
	err := req.Pages(ctx, func(page *gmail.ListMessagesResponse) error {

		if fetched >= max {
			logrus.Infof("Fetched %d messages, stopping...", fetched)
			return nil
		}

		for _, m := range page.Messages {
			// Fetch full message details
			msg, err := s.Users.Messages.Get(userId, m.Id).Format("full").Do()
			if err != nil {
				logrus.Errorf("Error fetching message %s: %v", m.Id, err)
				continue
			}
			messages = append(messages, msg)
			fetched++
			logrus.Infof("Fetched message ID: %s | %d", m.Id, fetched)

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
