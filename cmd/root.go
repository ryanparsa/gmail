/*
Copyright © 2025 Ryan Parsa <imryanparsa@gmail.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var credentialsPath string
var tokenPath string
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gmail",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "config.yaml", "config file")
	rootCmd.PersistentFlags().StringVarP(&credentialsPath, "credentials", "c", "credentials.json", "Path to the credentials JSON file")
	rootCmd.PersistentFlags().StringVarP(&tokenPath, "token", "t", "token.json", "Path to the token JSON file")

}
