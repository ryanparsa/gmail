package cmd

import (
	"fmt"
	"github.com/ryanparsa/gmail/internal"
	"log"
	"time"

	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "A brief description of your command",

	Run: func(cmd *cobra.Command, args []string) {
		svc, err := internal.NewService(credentialsPath, tokenPath, scopes)
		if err != nil {
			log.Fatalln(err)
		}
		config, err := internal.NewConfigFromYAML(cfgFile)
		if err != nil {
			fmt.Printf("Failed to load configuration: %v\n", err)
			return
		}

		err = svc.CreateLabels(config.Labels)
		if err != nil {
			return
		}

		log.Println("Waiting for 2s to labels get be created...")
		time.Sleep(2 * time.Second)
		lm, err := svc.LabelsMap()
		if err != nil {
			log.Fatalln(err)
		}

		var updatedFilters internal.Filters

		for _, filter := range config.Filters {
			// Update AddLabelIds with existing label IDs
			for id, label := range filter.Action.AddLabelIds {
				if labelID, exists := lm[label]; exists {
					filter.Action.AddLabelIds[id] = labelID.Id
				} else {
					log.Printf("Label %s does not exist. Skipping AddLabelId mapping.", label)
				}
			}

			// Update RemoveLabelIds with existing label IDs
			for id, label := range filter.Action.RemoveLabelIds {
				if labelID, exists := lm[label]; exists {
					filter.Action.RemoveLabelIds[id] = labelID.Id
				} else {
					log.Printf("Label %s does not exist. Skipping RemoveLabelId mapping.", label)

				}
			}

			updatedFilters = append(updatedFilters, filter)
		}

		err = svc.CreateFilters(updatedFilters)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
