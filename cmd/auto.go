package cmd

import (
	"fmt"
	"github.com/openai/openai-go"
	"github.com/ryanparsa/gmail/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	openAIHost  string
	openAIKey   string
	openAIModel string
	numEmails   int64
)

func init() {
	rootCmd.AddCommand(autoCmd)

	// Flags for OpenAI host, model, and API key
	autoCmd.Flags().StringVar(&openAIHost, "openai-host", "https://api.openai.com", "The OpenAI API host")
	autoCmd.Flags().StringVar(&openAIKey, "openai-key", "", "The OpenAI API key (required)")
	autoCmd.Flags().StringVarP(&openAIModel, "openai-model", "m", openai.ChatModelGPT4o, "The OpenAI model to use")

	// Flag for the number of emails to fetch
	autoCmd.Flags().Int64VarP(&numEmails, "size", "s", 1000, "Number of emails to fetch")
}

// autoCmd represents the auto command
var autoCmd = &cobra.Command{
	Use:   "auto",
	Short: "Generate filters and labels using OpenAI and Gmail messages",
	Long:  `Fetch Gmail messages, analyze them with OpenAI, and generate filters and labels.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Validate required flags
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
		logrus.Info("Starting the 'auto' command...")

		// Step 1: Initialize Gmail Service
		logrus.Info("Initializing Gmail service...")
		svc, err := internal.NewService(credentialsPath, tokenPath, scopes)
		if err != nil {
			logrus.Fatalf("Failed to initialize Gmail service: %v", err)
		}

		// Step 2: Fetch Gmail Messages
		logrus.Infof("Fetching the latest %d emails...", numEmails)
		messages, err := svc.Messages(numEmails, "")
		if err != nil {
			logrus.Fatalf("Failed to fetch emails: %v", err)
		}
		logrus.Infof("Fetched %d emails successfully.", len(messages))

		// Step 3: Call OpenAI API
		logrus.Infof("Generating filters and labels using OpenAI model '%s'...", openAIModel)
		response, err := internal.GetFiltersAndLabelsFromAI(openAIKey, openAIHost, openAIModel, messages)
		if err != nil {
			logrus.Fatalf("Failed to generate filters and labels: %v", err)
		}
		logrus.Info("Filters and labels generated successfully.")

		// Step 4: Save Output to File
		outputFile := "filters_and_labels.json"
		logrus.Infof("Saving filters and labels to file: %s", outputFile)
		err = response.SaveToFile(outputFile)
		if err != nil {
			logrus.Fatalf("Failed to save filters and labels to file: %v", err)
		}
		logrus.Infof("Filters and labels saved successfully to %s.", outputFile)

		logrus.Info("auto command completed successfully.")
	},
}
