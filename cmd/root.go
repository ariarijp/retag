package cmd

import (
	"github.com/spf13/cobra"
)

var baseUrl string
var apiKey string

var rootCmd = &cobra.Command{
	Use:   "retag",
	Short: "Exporting and Importing tool for Redash Query name and tags",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVar(&baseUrl, "url", "http://demo.redash.io", "Redash URL")
	rootCmd.MarkPersistentFlagRequired("url")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "Redash Admin user's API Key")
	rootCmd.MarkPersistentFlagRequired("api-key")
}
