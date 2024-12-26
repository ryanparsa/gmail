package internal

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/api/gmail/v1"
	"strings"
)

// DeleteFilter deletes a specific Gmail filter by its ID.
func (s *Service) DeleteFilter(filter *gmail.Filter) error {
	logrus.Infof("Deleting filter with ID: %s", filter.Id)

	err := s.Users.Settings.Filters.Delete("me", filter.Id).Do()
	if err != nil {
		logrus.Errorf("Failed to delete filter with ID %s: %v", filter.Id, err)
		return err
	}

	logrus.Infof("Successfully deleted filter with ID: %s", filter.Id)
	return nil
}

// ApplyFilterActions applies the actions defined in a Gmail filter to a set of messages.
func (s *Service) ApplyFilterActions(action *gmail.FilterAction, messages []*gmail.Message) error {
	logrus.Infof("Applying filter actions: AddLabelIds=%v, RemoveLabelIds=%v", action.AddLabelIds, action.RemoveLabelIds)

	for _, msg := range messages {
		logrus.Infof("Applying actions to message ID: %s", msg.Id)

		modifyReq := &gmail.ModifyMessageRequest{
			AddLabelIds:    action.AddLabelIds,
			RemoveLabelIds: action.RemoveLabelIds,
		}

		_, err := s.Users.Messages.Modify("me", msg.Id, modifyReq).Do()
		if err != nil {
			logrus.Errorf("Failed to apply actions to message ID %s: %v", msg.Id, err)
			return err
		}
	}

	logrus.Info("Successfully applied filter actions to all messages")
	return nil
}

// DeleteFilters deletes a list of Gmail filters.
func (s *Service) DeleteFilters(filters []*gmail.Filter) error {
	logrus.Infof("Deleting %d filters", len(filters))

	for _, filter := range filters {
		err := s.DeleteFilter(filter)
		if err != nil {
			logrus.Errorf("Failed to delete filter with ID %s: %v", filter.Id, err)
			return err
		}
	}

	logrus.Info("Successfully deleted all filters")
	return nil
}

// Filters retrieves the list of Gmail filters for the current user.
func (s *Service) Filters() ([]*gmail.Filter, error) {
	logrus.Info("Fetching all Gmail filters")

	filters, err := s.Users.Settings.Filters.List("me").Do()
	if err != nil {
		logrus.Errorf("Failed to fetch Gmail filters: %v", err)
		return nil, err
	}

	logrus.Infof("Successfully fetched %d filters", len(filters.Filter))
	return filters.Filter, nil
}

// CreateFilter creates a new Gmail filter.
func (s *Service) CreateFilter(filter *gmail.Filter) error {
	logrus.Infof("Creating new filter with ID: %s", filter.Id)

	f, err := s.Users.Settings.Filters.Create("me", filter).Do()
	if err != nil {
		logrus.Errorf("Failed to create filter with ID %s: %v", filter.Id, err)
		return err
	}

	logrus.Infof("Successfully created filter with ID: %s", f.Id)
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
	logrus.Info("Building Gmail query from filter criteria")

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
		logrus.Warn("No filter criteria provided; returning an empty query")
		return ""
	}

	query := "(" + strings.Join(queryParts, " ") + ")"
	logrus.Infof("Built Gmail query: %s", query)
	return query
}
