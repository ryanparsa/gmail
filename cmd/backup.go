package cmd

import (
	"github.com/ryanparsa/gmail/internal"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var outputPath string

func init() {
	rootCmd.AddCommand(backupCmd)

	// Define and attach the `--output` flag
	backupCmd.Flags().StringVar(&outputPath, "output", "backup.yaml", "Path to save the backup YAML file")
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup Gmail settings (filters and labels) to a YAML file",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Starting the 'backup' command...")

		// Step 1: Initialize Gmail Service
		logrus.Info("Initializing Gmail service...")
		svc, err := internal.NewService(credentialsPath, tokenPath, scopes)
		if err != nil {
			logrus.Fatalf("Failed to initialize Gmail service: %v", err)
		}

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

		// Step 4: Create Backup Configuration
		logrus.Info("Creating backup configuration...")
		backupConfig := internal.NewConfig(filters, labels)

		// Step 5: Save Backup to File
		logrus.Infof("Saving backup to file: %s", outputPath)
		err = backupConfig.SaveToFile(outputPath)
		if err != nil {
			logrus.Errorf("Failed to save backup to file: %v", err)
			return
		}

		logrus.Infof("Backup saved successfully to %s.", outputPath)
		logrus.Info("Backup command completed.")
	},
}
