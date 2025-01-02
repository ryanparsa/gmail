/*
Copyright © 2025 Ryan Parsa <imryanparsa@gmail.com>
*/
package cmd

import (
	"github.com/ryanparsa/gmail/internal"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(wipeCmd)

}

// wipeCmd represents the wipe command
var wipeCmd = &cobra.Command{
	Use:   "wipe",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := internal.NewService(credentialsPath, tokenPath, scopes)
		if err != nil {
			log.Fatal(err)
		}

		// Fetch Gmail settings
		filters, err := svc.Filters()
		if err != nil {
			log.Fatal(err)
		}
		labels, err := svc.Labels()
		if err != nil {
			log.Fatal(err)
		}

		err = svc.DeleteFilters(filters)
		if err != nil {
			log.Printf("Failed to delete filters: %v\n", err)
		}

		err = svc.DeleteLabels(labels)
		if err != nil {
			log.Printf("Failed to delete labels: %v\n", err)
		}
	},
}
