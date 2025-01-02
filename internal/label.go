package internal

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/api/gmail/v1"
)

// DeleteLabels deletes user-defined labels (ignoring system labels).
func (s *Service) DeleteLabels(labels Labels) error {
	for _, label := range labels {
		if label.Type == "system" {
			logrus.Infof("Skipping system label: %s", label.Name)
			continue
		}
		err := s.Users.Labels.Delete("me", label.Id).Do()
		if err != nil {
			logrus.Errorf("Failed to delete label %s: %v", label.Name, err)
		} else {
			logrus.Infof("Label %s deleted successfully", label.Name)
		}
	}
	return nil
}

// CreateLabels creates new labels in the user's Gmail account.
func (s *Service) CreateLabels(l Labels) error {
	for _, label := range l {
		newLabel, err := s.Service.Users.Labels.Create("me", label).Do()
		if err != nil {
			logrus.Errorf("Failed to create label %s: %v", label.Name, err)
		} else {
			logrus.Infof("Label %s created successfully (ID: %s)", label.Name, newLabel.Id)
		}
	}
	return nil
}

// Labels retrieves all labels in the user's Gmail account.
func (s *Service) Labels() (Labels, error) {
	logrus.Info("Fetching all labels from Gmail...")
	res, err := s.Service.Users.Labels.List("me").Do()
	if err != nil {
		logrus.Errorf("Failed to fetch labels: %v", err)
		return nil, err
	}
	logrus.Infof("Fetched %d labels successfully", len(res.Labels))
	return res.Labels, nil
}

// LabelsMap creates a map of label names to Gmail label objects.
func (s *Service) LabelsMap() (map[string]*gmail.Label, error) {
	logrus.Info("Creating labels map...")
	labels, err := s.Labels()
	if err != nil {
		logrus.Errorf("Failed to fetch labels for mapping: %v", err)
		return nil, err
	}

	labelMap := make(map[string]*gmail.Label)
	for _, label := range labels {
		labelMap[label.Name] = label
	}
	logrus.Infof("Labels map created successfully with %d entries", len(labelMap))
	return labelMap, nil
}
