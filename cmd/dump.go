package cmd

import (
	"fmt"
	"time"

	"github.com/ryanparsa/gmail/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(dumpCmd)
}

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump Gmail labels and filters to a YAML file",
	Long: `The "dump" command retrieves Gmail labels and filters and writes them to a YAML file.
This is useful for backing up Gmail configurations.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve configuration using the helper function
		config := internal.GetConfigFromFlags(cmd)

		logrus.Infof("Initializing Gmail service with credentials: %s and token: %s", config.Credentials, config.Token)

		// Initialize Gmail service
		srv := internal.NewService(config.Credentials, config.Token, config.Scopes)
		// Fetch Gmail labels
		logrus.Info("Fetching Gmail labels...")
		labelsResp, err := srv.Labels()
		if err != nil {
			logrus.Fatalf("Failed to fetch Gmail labels: %v", err)
		}
		logrus.Infof("Fetched %d labels", len(labelsResp))

		// Fetch Gmail filters
		logrus.Info("Fetching Gmail filters...")
		filtersResp, err := srv.Filters()
		if err != nil {
			logrus.Fatalf("Failed to fetch Gmail filters: %v", err)
		}
		logrus.Infof("Fetched %d filters", len(filtersResp))

		// Set the labels and filters in viper
		viper.Set("filters", filtersResp)
		viper.Set("labels", labelsResp)

		// Generate a filename for the dump file
		filename := fmt.Sprintf("dump_%s.yaml", time.Now().Format("2006_01_02__15_04_05"))

		// Write to YAML file
		logrus.Infof("Writing labels and filters to YAML file: %s", filename)
		if err := viper.WriteConfigAs(filename); err != nil {
			logrus.Fatalf("Failed to write YAML file: %v", err)
		}

		logrus.Infof("Successfully wrote labels and filters to file: %s", filename)
	},
}
