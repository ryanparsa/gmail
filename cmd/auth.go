package cmd

import (
	"github.com/ryanparsa/gmail/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/api/gmail/v1"
)

var scopes []string

func init() {
	rootCmd.AddCommand(authCmd)

	authCmd.Flags().StringSliceVar(&scopes, "scopes", []string{
		gmail.GmailLabelsScope,
		gmail.GmailReadonlyScope,
		gmail.GmailSettingsBasicScope,
	}, "OAuth 2.0 scopes to request")
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with the Gmail API",
	Long:  `Authenticate with the Gmail API using the specified credentials and token files.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Starting the 'auth' command...")

		// Log the scopes being requested
		logrus.Infof("Requested OAuth 2.0 scopes: %v", scopes)

		// Attempt authentication
		logrus.Info("Attempting to authenticate with the Gmail API...")
		err := internal.Authenticate(credentialsPath, tokenPath, scopes)
		if err != nil {
			logrus.Fatalf("Authentication failed: %v", err)
		} else {
			logrus.Info("Authentication successful.")
		}

		logrus.Info("Auth command completed.")
	},
}
