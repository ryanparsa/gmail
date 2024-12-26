package internal

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Config struct {
	Credentials string
	Token       string
	Scopes      []string
	ConfigFile  string
	Wait        int
}

// GetConfigFromFlags retrieves and validates required flags for Gmail commands
func GetConfigFromFlags(cmd *cobra.Command) *Config {
	credentials, err := cmd.Flags().GetString("credentials")
	if err != nil {
		logrus.Fatalf("Failed to get 'credentials' flag: %v", err)
	}
	if credentials == "" {
		logrus.Fatal("The 'credentials' flag is required")
	}

	token, err := cmd.Flags().GetString("token")
	if err != nil {
		logrus.Fatalf("Failed to get 'token' flag: %v", err)
	}
	if token == "" {
		logrus.Fatal("The 'token' flag is required")
	}

	scopes, err := cmd.Flags().GetStringArray("scopes")
	if err != nil {
		logrus.Fatalf("Failed to get 'scopes' flag: %v", err)
	}

	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		logrus.Fatalf("Failed to get 'config' flag: %v", err)
	}

	wait, err := cmd.Flags().GetInt("wait")
	if err != nil {
		logrus.Fatalf("Failed to get 'wait' flag: %v", err)
	}

	return &Config{
		Credentials: credentials,
		Token:       token,
		Scopes:      scopes,
		ConfigFile:  configFile,
		Wait:        wait,
	}
}
