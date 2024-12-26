package cmd

import (
	"github.com/ryanparsa/gmail/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(authCmd)
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with the Gmail API",
	Long: `The "auth" command initializes the Gmail API service
by reading the provided credentials and token files, and performs authentication.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve configuration using the helper function
		config := internal.GetConfigFromFlags(cmd)

		logrus.Infof("Initializing Gmail service with credentials: %s and token: %s", config.Credentials, config.Token)

		// Initialize Gmail service
		srv := internal.NewService(config.Credentials, config.Token, config.Scopes)
		if srv != nil {
			logrus.Info("Successfully authenticated with the Gmail API")
		} else {
			logrus.Error("Failed to authenticate with the Gmail API")
		}
	},
}
