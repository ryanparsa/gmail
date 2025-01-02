package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var credentialsPath string
var tokenPath string
var cfgFile string

func init() {

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "config.yaml", "config file")
	rootCmd.PersistentFlags().StringVarP(&credentialsPath, "credentials", "c", "credentials.json", "Path to the credentials JSON file")
	rootCmd.PersistentFlags().StringVarP(&tokenPath, "token", "t", "token.json", "Path to the token JSON file")

}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gmail",
	Short: "A CLI tool to interact with Gmail",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logrus.Fatalln(err)
	}
}
