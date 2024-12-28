package cmd

import (
	"github.com/ryanparsa/gmail/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(wipeCmd)
}

var wipeCmd = &cobra.Command{
	Use:   "wipe",
	Short: "Delete all user-created Gmail labels and filters",
	Long: `The "wipe" command deletes all user-created Gmail labels and filters. 
It does not remove system labels such as "INBOX" or "SPAM".`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve configuration using the helper function
		config := internal.GetConfigFromFlags(cmd)

		logrus.Infof("Initializing Gmail service with credentials: %s and token: %s", config.Credentials, config.Token)

		// Initialize Gmail service
		srv := internal.NewService(config.Credentials, config.Token, config.Scopes)

		// Fetch and delete labels
		logrus.Info("Fetching user-created labels...")
		labelsResp, err := srv.Labels()
		if err != nil {
			logrus.Fatalf("Failed to fetch labels: %v", err)
		}
		logrus.Infof("Fetched %d labels.", len(labelsResp))

		logrus.Info("Deleting user-created labels...")
		err = srv.DeleteLabel(labelsResp...)
		if err != nil {
			logrus.Fatalf("Failed to delete labels: %v", err)
		}
		logrus.Info("All user-created labels deleted successfully.")

		// Fetch and delete filters
		logrus.Info("Fetching Gmail filters...")
		filtersResp, err := srv.Filters()
		if err != nil {
			logrus.Fatalf("Failed to fetch filters: %v", err)
		}
		logrus.Infof("Fetched %d filters.", len(filtersResp))

		logrus.Info("Deleting Gmail filters...")
		err = srv.DeleteFilter(filtersResp...)
		if err != nil {
			logrus.Fatalf("Failed to delete filters: %v", err)
		}
		logrus.Info("All Gmail filters deleted successfully.")

		logrus.Info("Gmail wipe process completed successfully.")
	},
}
