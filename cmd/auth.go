package cmd

import (
	"fmt"
	"github.com/ryanparsa/gmail/internal"
	"github.com/spf13/cobra"
	"google.golang.org/api/gmail/v1"
	"log"
)

var scopes []string // Local flag for specifying API scopes

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.Flags().StringSliceVar(&scopes, "scopes", []string{
		gmail.GmailLabelsScope,
		gmail.GmailReadonlyScope,
		gmail.GmailSettingsBasicScope,
	}, "")

}

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.Authenticate(credentialsPath, tokenPath, scopes)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Authenticated successfully")
		}
	},
}
