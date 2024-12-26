package cmd

import (
	"github.com/ryanparsa/gmail/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration files from templates",
	Long: `The "init" command checks for the existence of required configuration files 
and creates them using predefined templates if they do not exist.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Initializing configuration files...")

		// Loop through all templates
		for file, template := range internal.Templates {
			logrus.Infof("Checking file: %s", file)

			// Check if the file exists
			if _, err := os.Stat(file); os.IsNotExist(err) {
				logrus.Infof("File %s does not exist. Creating it...", file)

				// Save the template file
				err = internal.SaveTemplate(file, template)
				if err != nil {
					logrus.Fatalf("Failed to create file %s: %v", file, err)
				}
				logrus.Infof("File %s created successfully.", file)
			} else if err != nil {
				logrus.Fatalf("Failed to check file %s: %v", file, err)
			} else {
				logrus.Warnf("File %s already exists. Skipping...", file)
			}
		}

		logrus.Info("Initialization complete.")
	},
}
