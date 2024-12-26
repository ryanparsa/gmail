package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/api/gmail/v1"
)

// Default Gmail API scopes
var (
	defaultScopes = []string{
		gmail.GmailModifyScope,
		gmail.GmailLabelsScope,
		gmail.GmailSettingsBasicScope,
	}
)

// Root command for the CLI application
var rootCmd = &cobra.Command{
	Use:   "gmail",
	Short: "CLI tool to manage Gmail labels and filters",
	Long: `This application provides a command-line interface (CLI) to interact 
with the Gmail API, allowing you to manage labels, filters, and other Gmail settings.`,
}

// Initialization function to define persistent flags
func init() {
	// Define persistent flags for the root command
	rootCmd.PersistentFlags().String("config", "config.yaml", "Path to the YAML configuration file")
	rootCmd.PersistentFlags().String("credentials", "credentials.json", "Path to the credentials JSON file")
	rootCmd.PersistentFlags().String("token", "token.json", "Path to the token JSON file")
	rootCmd.PersistentFlags().Int("wait", 5, "Wait time in seconds between operations")
	rootCmd.PersistentFlags().StringArray("scopes", defaultScopes, "Gmail API scopes to request")
}

// Execute starts the CLI application
func Execute() {

	// Execute the root command
	err := rootCmd.Execute()
	if err != nil {
		logrus.Fatalf("Error executing command: %v", err)
	}

}
