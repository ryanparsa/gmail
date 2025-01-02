package internal

import (
	"strings"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/gmail/v1"
)

// Filters fetches all Gmail filters for the user.
func (s *Service) Filters() (Filters, error) {
	logrus.Info("Fetching Gmail filters...")
	res, err := s.Service.Users.Settings.Filters.List(userId).Do()
	if err != nil {
		logrus.Errorf("Failed to fetch Gmail filters: %v", err)
		return nil, err
	}
	logrus.Infof("Fetched %d filters successfully.", len(res.Filter))
	return res.Filter, nil
}

// ApplyFilterActions applies the actions defined in a Gmail filter to a set of messages.
func (s *Service) ApplyFilterActions(action *gmail.FilterAction, messages []*gmail.Message) error {
	logrus.Infof("Applying filter actions to %d messages...", len(messages))

	for _, msg := range messages {
		modifyReq := &gmail.ModifyMessageRequest{
			AddLabelIds:    action.AddLabelIds,
			RemoveLabelIds: action.RemoveLabelIds,
		}

		_, err := s.Users.Messages.Modify(userId, msg.Id, modifyReq).Do()
		if err != nil {
			logrus.Errorf("Failed to apply filter actions to message %s: %v", msg.Id, err)
			return err
		}
		logrus.Debugf("Filter actions applied successfully to message %s", msg.Id)
	}

	logrus.Info("Filter actions applied successfully to all messages.")
	return nil
}

// BuildQueryFromFilter constructs a Gmail search query from a filter's criteria.
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
		logrus.Warn("No criteria provided for filter. Returning an empty query.")
		return ""
	}

	query := "(" + strings.Join(queryParts, " ") + ")"
	logrus.Debugf("Constructed query from filter criteria: %s", query)
	return query
}

// DeleteFilters deletes Gmail filters for the user.
func (s *Service) DeleteFilters(filters Filters) error {
	logrus.Infof("Deleting %d Gmail filters...", len(filters))
	for _, filter := range filters {
		err := s.Users.Settings.Filters.Delete(userId, filter.Id).Do()
		if err != nil {
			logrus.Errorf("Failed to delete filter %s: %v", filter.Id, err)
		} else {
			logrus.Infof("Filter %s deleted successfully.", filter.Id)
		}
	}
	logrus.Info("All specified filters deleted successfully.")
	return nil
}

// CreateFilters creates new Gmail filters for the user.
func (s *Service) CreateFilters(f Filters) error {
	logrus.Infof("Creating %d Gmail filters...", len(f))

	for _, filter := range f {
		newFilter, err := s.Service.Users.Settings.Filters.Create(userId, filter).Do()
		if err != nil {
			logrus.Errorf("Failed to create filter: %v", err)
		} else {
			logrus.Infof("Filter %s created successfully.", newFilter.Id)
		}
	}

	logrus.Info("All specified filters created successfully.")
	return nil
}
