package internal

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/api/gmail/v1"
)

// Labels retrieves all Gmail labels for the authenticated user.
//
// Returns:
// - A slice of Gmail labels.
// - An error if the label retrieval fails.
func (s *Service) Labels() ([]*gmail.Label, error) {
	logrus.Info("Fetching all Gmail labels")

	labelsResp, err := s.Users.Labels.List("me").Do()
	if err != nil {
		logrus.Errorf("Failed to fetch Gmail labels: %v", err)
		return nil, err
	}

	logrus.Infof("Successfully fetched %d labels", len(labelsResp.Labels))
	return labelsResp.Labels, nil
}

// DeleteLabel deletes a single Gmail label.
//
// Parameters:
// - label: The Gmail label to be deleted.
//
// Returns:
// - An error if the label deletion fails or if the label is a system label (which cannot be deleted).
func (s *Service) DeleteLabel(labels ...*gmail.Label) error {
	for _, label := range labels {
		if label.Type == "system" {
			logrus.Warnf("Skipping deletion of system label: %s (ID: %s)", label.Name, label.Id)
			return nil
		}

		logrus.Infof("Deleting label: %s (ID: %s)", label.Name, label.Id)
		err := s.Users.Labels.Delete("me", label.Id).Do()
		if err != nil {
			logrus.Errorf("Failed to delete label: %s (ID: %s), error: %v", label.Name, label.Id, err)
			return err
		}

		logrus.Infof("Successfully deleted label: %s (ID: %s)", label.Name, label.Id)
	}
	return nil
}

// LabelsMap creates a map of Gmail labels where the key is the label name and the value is the Gmail label.
//
// Returns:
// - A map of Gmail labels keyed by their names.
// - An error if the label retrieval fails.
func (s *Service) LabelsMap() (map[string]*gmail.Label, error) {
	logrus.Info("Building map of Gmail labels")

	labels, err := s.Labels()
	if err != nil {
		logrus.Errorf("Failed to fetch Gmail labels for map creation: %v", err)
		return nil, err
	}

	labelMap := make(map[string]*gmail.Label)
	for _, label := range labels {
		logrus.Infof("Adding label to map: %s (ID: %s)", label.Name, label.Id)
		labelMap[label.Name] = label
	}

	logrus.Info("Successfully created map of Gmail labels")
	return labelMap, nil
}

// CreateLabel creates a new Gmail label.
//
// Parameters:
// - label: The Gmail label to be created.
//
// Returns:
// - An error if the label creation fails.
func (s *Service) CreateLabel(labels ...*gmail.Label) error {
	for _, label := range labels {

		logrus.Infof("Creating label: %s", label.Name)

		_, err := s.Users.Labels.Create("me", label).Do()
		if err != nil {
			logrus.Errorf("Failed to create label: %s, error: %v", label.Name, err)
			return err
		}

		logrus.Infof("Successfully created label: %s", label.Name)
	}
	return nil
}
