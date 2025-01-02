package internal

import (
	"google.golang.org/api/gmail/v1"
	"log"
	"strings"
)

func (s *Service) Filters() (Filters, error) {
	res, err := s.Service.Users.Settings.Filters.List("me").Do()
	if err != nil {
		return nil, err
	}

	return res.Filter, nil
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
