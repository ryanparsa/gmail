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
		for fileName, templateFilePath := range internal.Templates {
			logrus.Infof("Checking file: %s", fileName)

			// Check if the file exists
			if _, err := os.Stat(fileName); os.IsNotExist(err) {
				logrus.Infof("File %s does not exist. Creating it...", fileName)

				// Save the template file
				err = internal.SaveTemplate(fileName, templateFilePath)
				if err != nil {
					logrus.Fatalf("Failed to create file %s: %v", fileName, err)
				}
				logrus.Infof("File %s created successfully.", fileName)
			} else if err != nil {
				logrus.Fatalf("Failed to check file %s: %v", fileName, err)
			} else {
				logrus.Warnf("File %s already exists. Skipping...", fileName)
			}
		}

		logrus.Info("Initialization complete.")
	},
}
