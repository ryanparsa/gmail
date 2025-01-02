package cmd

import (
	"fmt"
	"github.com/openai/openai-go"
	"github.com/ryanparsa/gmail/internal"
	"github.com/spf13/cobra"
	"log"
)

var (
	openAIHost  string
	openAIKey   string
	openAIModel string
	numEmails   int64
)

func init() {
	rootCmd.AddCommand(smartCmd)
	// Flags for OpenAI host, model, and API key
	smartCmd.Flags().StringVar(&openAIHost, "openai-host", "https://api.openai.com", "The OpenAI API host")
	smartCmd.Flags().StringVar(&openAIKey, "openai-key", "", "The OpenAI API key (required)")
	smartCmd.Flags().StringVarP(&openAIModel, "openai-model", "m", openai.ChatModelGPT4o, "The OpenAI model to use")

	// Flag for the number of emails to fetch
	smartCmd.Flags().Int64VarP(&numEmails, "size", "s", 1000, "Number of emails to fetch")

}

// smartCmd represents the smart command
var smartCmd = &cobra.Command{
	Use: "smart",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if openAIKey == "" {
			return fmt.Errorf("openai-key is required")
		}

		if numEmails < 1 {
			return fmt.Errorf("size must be greater than 0")
		}

		if openAIModel == "" {
			return fmt.Errorf("openai-model is required")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := internal.NewService(credentialsPath, tokenPath, scopes)
		if err != nil {
			log.Fatalln(err)
		}

		m, err := svc.Messages(numEmails)
		if err != nil {
			log.Fatalln(err)
		}

		// Call OpenAI to generate filters and labels
		response, err := internal.GetFiltersAndLabelsFromAI(openAIKey, openAIHost, openAIModel, m)
		if err != nil {
			log.Fatalf("Failed to generate filters and labels: %v", err)
		}

		// Save the generated filters and labels to a file
		outputFile := "filters_and_labels.json"

		err = response.SaveToFile(outputFile)
		if err != nil {
			log.Fatalln(err)
		} else {
			fmt.Printf("Filters and labels saved to %s\n", outputFile)
		}

	},
}
