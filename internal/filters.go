package internal

import (
	"context"
	"google.golang.org/api/gmail/v1"
	"log"
	"strings"
	"time"
)

func (s *Service) Filters() (Filters, error) {
	res, err := s.Service.Users.Settings.Filters.List("me").Do()
	if err != nil {
		return nil, err
	}

	return res.Filter, nil
}

func (s *Service) LabelsMap() (map[string]*gmail.Label, error) {

	labels, err := s.Labels()
	if err != nil {
		return nil, err
	}

	labelMap := make(map[string]*gmail.Label)
	for _, label := range labels {
		labelMap[label.Name] = label
	}

	return labelMap, nil
}

func (s *Service) CreateFilters(f Filters) error {

	for _, filter := range f {
		newFilter, err := s.Service.Users.Settings.Filters.Create("me", filter).Do()
		if err != nil {
			log.Printf("Failed to create filter %s: %v", filter.Id, err)
		} else {
			log.Printf("Filter %s created successfully", newFilter.Id)
		}
	}

	return nil
}

func (s *Service) Labels() (Labels, error) {
	res, err := s.Service.Users.Labels.List("me").Do()

	if err != nil {
		return nil, err
	}

	return res.Labels, nil
}

func (s *Service) CreateLabels(l Labels) error {
	for _, label := range l {
		newLabel, err := s.Service.Users.Labels.Create("me", label).Do()
		if err != nil {
			log.Printf("Failed to create label %s: %v", label.Name, err)
		} else {
			log.Printf("Label %s created successfully", newLabel.Id)
		}
	}

	return nil
}

func (s *Service) DeleteFilters(filters Filters) error {
	for _, filter := range filters {
		err := s.Users.Settings.Filters.Delete("me", filter.Id).Do()
		if err != nil {
			log.Printf("Failed to delete filter %s: %v", filter.Id, err)
		} else {
			log.Printf("Filter %s deleted successfully", filter.Id)
		}
	}
	return nil
}

func (s *Service) DeleteLabels(labels Labels) error {
	for _, label := range labels {
		if label.Type == "system" {
			log.Printf("Skipping system label %s", label.Name)
			continue
		}
		err := s.Users.Labels.Delete("me", label.Id).Do()
		if err != nil {
			log.Printf("Failed to delete label %s: %v", label.Name, err)
		} else {
			log.Printf("Label %s deleted successfully", label.Name)
		}
	}

	return nil
}

// BuildQueryFromFilter constructs a Gmail search query from a filter's criteria.
//
// Example:
// - From: "example@example.com"
// - To: "recipient@example.com"
// - Subject: "Meeting"
// - Query: "is:unread"
func (s *Service) BuildQueryFromFilter(criteria *gmail.FilterCriteria) string {

	var queryParts []string
	if criteria.From != "" {
		queryParts = append(queryParts, "from:"+criteria.From)
	}
	if criteria.To != "" {
		queryParts = append(queryParts, "to:"+criteria.To)
	}
	if criteria.Subject != "" {
		queryParts = append(queryParts, "subject:"+criteria.Subject)
	}
	if criteria.Query != "" {
		queryParts = append(queryParts, criteria.Query)
	}

	// If no criteria is provided, return an empty query
	if len(queryParts) == 0 {
		return ""
	}

	query := "(" + strings.Join(queryParts, " ") + ")"
	return query
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

	var messages []*gmail.Message
	pageToken := ""

	for {
		// List messages with the specified query and page token
		req := s.Users.Messages.List("me").Q(query).PageToken(pageToken)
		resp, err := req.Do()
		if err != nil {
			return nil, err
		}

		// Append fetched messages
		messages = append(messages, resp.Messages...)

		// Break if there are no more pages
		if resp.NextPageToken == "" {
			break
		}

		// Update page token for the next iteration
		pageToken = resp.NextPageToken
	}

	return messages, nil
}

// ApplyFilterActions applies the actions defined in a Gmail filter to a set of messages.
func (s *Service) ApplyFilterActions(action *gmail.FilterAction, messages []*gmail.Message) error {

	for _, msg := range messages {

		modifyReq := &gmail.ModifyMessageRequest{
			AddLabelIds:    action.AddLabelIds,
			RemoveLabelIds: action.RemoveLabelIds,
		}

		_, err := s.Users.Messages.Modify("me", msg.Id, modifyReq).Do()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Messages(max int64) (Messages, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var messages Messages

	req := s.Users.Messages.List("me").MaxResults(max)

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
