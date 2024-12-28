package cmd

import (
	"github.com/ryanparsa/gmail/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/gmail/v1"
)

func init() {
	rootCmd.AddCommand(loadCmd)

	// Define flags
	loadCmd.Flags().Int("wait", 2, "Wait time (in seconds) between operations")
}

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load and sync Gmail labels and filters from a YAML file",
	Long: `The "load" command reads Gmail labels and filters from a YAML configuration file 
and syncs them with the user's Gmail account.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve configuration using the helper function
		config := internal.GetConfigFromFlags(cmd)

		logrus.Infof("Initializing Gmail service with credentials: %s and token: %s", config.Credentials, config.Token)

		// Initialize Gmail service
		srv := internal.NewService(config.Credentials, config.Token, config.Scopes)
		internal.Wait(config.Wait, "Starting sync process...")

		// Read YAML configuration
		logrus.Infof("Reading configuration from file: %s", config.ConfigFile)
		dataConfig := viper.New()
		dataConfig.SetConfigFile(config.ConfigFile)
		dataConfig.SetConfigType("yaml")

		var data struct {
			Labels  []gmail.Label  `yaml:"labels"`
			Filters []gmail.Filter `yaml:"filters"`
		}

		if err := dataConfig.ReadInConfig(); err != nil {
			logrus.Fatalf("Failed to read configuration file: %v", err)
		}

		if err := dataConfig.Unmarshal(&data); err != nil {
			logrus.Fatalf("Failed to parse configuration file: %v", err)
		}

		// Sync Labels
		logrus.Info("Syncing Gmail labels...")
		for _, label := range data.Labels {
			logrus.Infof("Creating label: %s", label.Name)
			err := srv.CreateLabel(&label)
			if err != nil {
				logrus.Errorf("Failed to create label %s: %v", label.Name, err)
			} else {
				logrus.Infof("Label %s created successfully", label.Name)
			}
		}

		internal.Wait(config.Wait, "Waiting before syncing filters...")

		// Fetch existing labels and build a label map
		logrus.Info("Fetching existing labels...")
		labelMap, err := srv.LabelsMap()
		if err != nil {
			logrus.Fatalf("Failed to fetch existing labels: %v", err)
		}

		// Sync Filters
		logrus.Info("Syncing Gmail filters...")
		for _, filter := range data.Filters {
			// Update AddLabelIds with existing label IDs
			for id, label := range filter.Action.AddLabelIds {
				if labelID, exists := labelMap[label]; exists {
					filter.Action.AddLabelIds[id] = labelID.Id
				} else {
					logrus.Warnf("Label %s does not exist. Skipping AddLabelId mapping.", label)
				}
			}

			// Update RemoveLabelIds with existing label IDs
			for id, label := range filter.Action.RemoveLabelIds {
				if labelID, exists := labelMap[label]; exists {
					filter.Action.RemoveLabelIds[id] = labelID.Id
				} else {
					logrus.Warnf("Label %s does not exist. Skipping RemoveLabelId mapping.", label)
				}
			}

			// Create the filter
			logrus.Infof("Creating filter with criteria: %+v", filter.Criteria)
			err := srv.CreateFilter(&filter)
			if err != nil {
				logrus.Errorf("Failed to create filter: %v", err)
			} else {
				logrus.Info("Filter created successfully")
			}
		}

		logrus.Info("Gmail label and filter synchronization complete.")
	},
}
