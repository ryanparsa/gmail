package cmd

import (
	"github.com/ryanparsa/gmail/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(wipeCmd)
}

// wipeCmd represents the wipe command
var wipeCmd = &cobra.Command{
	Use:   "wipe",
	Short: "Delete all Gmail filters and labels",
	Long: `The wipe command deletes all filters and labels from the connected Gmail account.
It ensures a clean slate by removing user-defined filters and labels.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Starting the 'wipe' command...")

		// Step 1: Initialize Gmail Service
		logrus.Info("Initializing Gmail service...")
		svc, err := internal.NewService(credentialsPath, tokenPath, scopes)
		if err != nil {
			logrus.Fatalf("Failed to initialize Gmail service: %v", err)
		}
		logrus.Info("Gmail service initialized successfully.")

		// Step 2: Fetch Gmail Filters
		logrus.Info("Fetching Gmail filters...")
		filters, err := svc.Filters()
		if err != nil {
			logrus.Fatalf("Failed to fetch Gmail filters: %v", err)
		}
		logrus.Infof("Fetched %d filters successfully.", len(filters))

		// Step 3: Fetch Gmail Labels
		logrus.Info("Fetching Gmail labels...")
		labels, err := svc.Labels()
		if err != nil {
			logrus.Fatalf("Failed to fetch Gmail labels: %v", err)
		}
		logrus.Infof("Fetched %d labels successfully.", len(labels))

		// Step 4: Delete Filters
		if len(filters) > 0 {
			logrus.Info("Deleting Gmail filters...")
			err = svc.DeleteFilters(filters)
			if err != nil {
				logrus.Errorf("Failed to delete filters: %v", err)
			} else {
				logrus.Info("Filters deleted successfully.")
			}
		} else {
			logrus.Info("No filters to delete.")
		}

		// Step 5: Delete Labels
		if len(labels) > 0 {
			logrus.Info("Deleting Gmail labels...")
			err = svc.DeleteLabels(labels)
			if err != nil {
				logrus.Errorf("Failed to delete labels: %v", err)
			} else {
				logrus.Info("Labels deleted successfully.")
			}
		} else {
			logrus.Info("No labels to delete.")
		}

		logrus.Info("Wipe command completed successfully.")
	},
}
