package cmd

import (
	"fmt"
	"github.com/ryanparsa/gmail/internal"
	"log"

	"github.com/spf13/cobra"
)

var outputPath string

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := internal.NewService(credentialsPath, tokenPath, scopes)

		// Fetch Gmail settings
		filters, err := svc.Filters()
		if err != nil {
			log.Fatal(err)
		}
		labels, err := svc.Labels()
		if err != nil {
			log.Fatal(err)
		}

		b := internal.NewConfig(filters, labels)

		// Save the backup to a file
		err = b.SaveToFile(outputPath)
		if err != nil {
			fmt.Printf("Failed to save backup: %v\n", err)
			return
		}

		fmt.Printf("Backup saved successfully to %s\n", outputPath)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().StringVar(&outputPath, "output", "backup.yaml", "Path to save the backup YAML file")
}
