package cmd

import (
	"github.com/ryanparsa/gmail/internal"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clear the saved token file to reset authentication",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Starting the 'clean' command...")

		// Attempt to clear the token file
		logrus.Infof("Attempting to clear the token file at path: %s", tokenPath)
		err := internal.ClearToken(tokenPath)
		if err != nil {
			// Log an error message if clearing the token fails
			logrus.Errorf("Failed to clear token: %v", err)
		} else {
			// Log a success message if the token is cleared successfully
			logrus.Info("Token cleared successfully.")
		}

		logrus.Info("Clean command completed.")
	},
}
