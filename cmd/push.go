package cmd

import (
	"github.com/ryanparsa/gmail/internal"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pushCmd)
}

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push Gmail labels and filters configuration",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Starting the 'push' command...")

		// Step 1: Initialize Gmail Service
		logrus.Info("Initializing Gmail service...")
		svc, err := internal.NewService(credentialsPath, tokenPath, scopes)
		if err != nil {
			logrus.Fatalf("Failed to initialize Gmail service: %v", err)
		}

		// Step 2: Load Configuration
		logrus.Infof("Loading configuration from file: %s", cfgFile)
		config, err := internal.NewConfigFromYAML(cfgFile)
		if err != nil {
			logrus.Errorf("Failed to load configuration: %v", err)
			return
		}
		logrus.Info("Configuration loaded successfully.")

		// Step 3: Create Labels
		logrus.Info("Creating labels...")
		err = svc.CreateLabels(config.Labels)
		if err != nil {
			logrus.Errorf("Failed to create labels: %v", err)
			return
		}
		logrus.Info("Labels created successfully. Waiting 2 seconds for propagation...")
		time.Sleep(2 * time.Second)

		// Step 4: Fetch Existing Labels Map
		logrus.Info("Fetching existing labels map...")
		lm, err := svc.LabelsMap()
		if err != nil {
			logrus.Fatalf("Failed to fetch labels map: %v", err)
		}
		logrus.Infof("Labels map retrieved: %d labels found.", len(lm))

		// Step 5: Update Filters with Label IDs
		logrus.Info("Updating filters with existing label IDs...")
		var updatedFilters internal.Filters
		for _, filter := range config.Filters {
			logrus.Infof("Processing filter: %v", filter)

			// Update AddLabelIds
			for id, label := range filter.Action.AddLabelIds {
				if labelID, exists := lm[label]; exists {
					filter.Action.AddLabelIds[id] = labelID.Id
				} else {
					logrus.Warnf("Label '%s' does not exist. Skipping AddLabelId mapping.", label)
				}
			}

			// Update RemoveLabelIds
			for id, label := range filter.Action.RemoveLabelIds {
				if labelID, exists := lm[label]; exists {
					filter.Action.RemoveLabelIds[id] = labelID.Id
					logrus.Infof("Mapped RemoveLabelId '%s' to ID '%s'.", label, labelID.Id)
				} else {
					logrus.Warnf("Label '%s' does not exist. Skipping RemoveLabelId mapping.", label)
				}
			}

			updatedFilters = append(updatedFilters, filter)
		}
		logrus.Info("Filters updated successfully.")

		// Step 6: Create Filters
		logrus.Info("Creating filters...")
		err = svc.CreateFilters(updatedFilters)
		if err != nil {
			logrus.Errorf("Failed to create filters: %v", err)
			return
		}
		logrus.Info("Filters created successfully.")
		logrus.Info("Push command completed.")
	},
}
