package cmd

import (
	"github.com/ryanparsa/gmail/internal"
	"github.com/spf13/cobra"
	"log"
)

// retroCmd represents the retro command
var retroCmd = &cobra.Command{
	Use:   "retro",
	Short: "A brief description of your command",

	Run: func(cmd *cobra.Command, args []string) {
		svc, err := internal.NewService(credentialsPath, tokenPath, scopes)
		if err != nil {
			log.Fatalln(err)
		}

		filters, err := svc.Filters()
		if err != nil {
			log.Fatalln(err)
		}

		for _, filter := range filters {

			// Build the query from filter criteria
			query := svc.BuildQueryFromFilter(filter.Criteria)
			if query == "" {
				continue
			}

			// Fetch emails matching the query
			messages, err := svc.FetchEmails(query)
			if err != nil {
				continue
			}

			// Apply filter actions to the messages
			if err := svc.ApplyFilterActions(filter.Action, messages); err != nil {
				continue
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(retroCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// retroCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// retroCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
