package cmd

import (
	"github.com/ryanparsa/gmail/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(retroCmd)

}

var retroCmd = &cobra.Command{
	Use:   "retro",
	Short: "Retroactively apply Gmail filters to existing messages",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Starting the 'retro' command...")

		// Step 1: Initialize Gmail Service
		logrus.Info("Initializing Gmail service...")
		svc, err := internal.NewService(credentialsPath, tokenPath, scopes)
		if err != nil {
			logrus.Fatalf("Failed to initialize Gmail service: %v", err)
		}
		logrus.Info("Gmail service initialized successfully.")

		// Step 2: Fetch Filters
		logrus.Info("Fetching Gmail filters...")
		filters, err := svc.Filters()
		if err != nil {
			logrus.Fatalf("Failed to fetch Gmail filters: %v", err)
		}
		logrus.Infof("Fetched %d filters successfully.", len(filters))

		// Step 3: Process Each Filter
		for _, filter := range filters {
			logrus.Infof("Processing filter: %v", filter)

			// Build the query from filter criteria
			query := svc.BuildQueryFromFilter(filter.Criteria)
			if query == "" {
				logrus.Warn("Empty query generated from filter criteria. Skipping this filter.")
				continue
			}
			logrus.Infof("Built query: %s", query)

			// Fetch emails matching the query
			logrus.Infof("Fetching emails matching query: %s", query)
			messages, err := svc.Messages(100000, query)
			if err != nil {
				logrus.Errorf("Failed to fetch emails for query '%s': %v", query, err)
				continue
			}
			logrus.Infof("Fetched %d messages matching the query.", len(messages))

			// Apply filter actions to the messages
			logrus.Info("Applying filter actions to the fetched messages...")
			if err := svc.ApplyFilterActions(filter.Action, messages); err != nil {
				logrus.Errorf("Failed to apply filter actions: %v", err)
				continue
			}
			logrus.Info("Filter actions applied successfully.")
		}

		logrus.Info("Retro command completed.")
	},
}
