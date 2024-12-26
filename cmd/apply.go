package cmd

import (
	"github.com/ryanparsa/gmail/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(applyCmd)

}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply filters to Gmail messages",
	Long: `The "apply" command retrieves Gmail filters and applies their actions
to messages matching the filter criteria.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve configuration using the helper function
		config := internal.GetConfigFromFlags(cmd)

		logrus.Infof("Starting Gmail service with credentials: %s, token: %s", config.Credentials, config.Token)

		// Initialize Gmail service
		srv := internal.NewService(config.Credentials, config.Token, config.Scopes)

		// Fetch filters
		filtersResp, err := srv.Filters()
		if err != nil {
			logrus.Fatalf("Failed to fetch filters: %v", err)
		}
		logrus.Infof("Successfully fetched %d filters", len(filtersResp))

		// Process each filter
		for _, filter := range filtersResp {
			logrus.Infof("Processing filter: %s (ID: %s)", filter.Criteria.Query, filter.Id)

			// Build the query from filter criteria
			query := srv.BuildQueryFromFilter(filter.Criteria)
			if query == "" {
				logrus.Warnf("Skipping filter with empty query: %s (ID: %s)", filter.Criteria.Query, filter.Id)
				continue
			}
			logrus.Infof("Built query: %s", query)

			// Fetch emails matching the query
			messages, err := srv.FetchEmails(query)
			if err != nil {
				logrus.Errorf("Failed to fetch emails for query: %s, error: %v", query, err)
				continue
			}
			logrus.Infof("Fetched %d messages for query: %s", len(messages), query)

			// Apply filter actions to the messages
			if err := srv.ApplyFilterActions(filter.Action, messages); err != nil {
				logrus.Errorf("Failed to apply filter actions for filter ID: %s, error: %v", filter.Id, err)
				continue
			}
			logrus.Infof("Successfully applied actions for filter ID: %s", filter.Id)
		}

		logrus.Info("Successfully applied all filters")
	},
}
