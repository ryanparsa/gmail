package internal

import (
	"context"
	"google.golang.org/api/gmail/v1"
	"log"
	"time"
)

func (s *Service) Messages(max int64, query string) (Messages, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var messages Messages

	req := s.Users.Messages.List("me").Q(query).MaxResults(max)

	err := req.Pages(ctx, func(page *gmail.ListMessagesResponse) error {
		for _, m := range page.Messages {
			msg, err := s.Users.Messages.Get("me", m.Id).Format("full").Do()
			if err != nil {
				log.Printf("Error fetching message %s: %v", m.Id, err)
				continue
			}
			messages = append(messages, msg)
		}
		return nil

	})
	if err != nil {
		return nil, err
	}

	return messages, nil
}
