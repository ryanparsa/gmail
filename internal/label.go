package internal

import (
	"google.golang.org/api/gmail/v1"
	"log"
)

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

func (s *Service) Labels() (Labels, error) {
	res, err := s.Service.Users.Labels.List("me").Do()

	if err != nil {
		return nil, err
	}

	return res.Labels, nil
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
