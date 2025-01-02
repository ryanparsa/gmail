/*
Copyright © 2025 Ryan Parsa <imryanparsa@gmail.com>
*/
package cmd

import (
	"github.com/ryanparsa/gmail/internal"
	"log"

	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "A brief description of your command",

	Run: func(cmd *cobra.Command, args []string) {
		err := internal.ClearToken(tokenPath)
		if err != nil {
			log.Printf("Failed to clear token: %v\n", err)
		} else {
			log.Printf("Token cleared successfully\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
