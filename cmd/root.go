/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gcscli",
	Short: "CLI tool for managing files in Google Cloud Storage (GCS)",

	Long: `This CLI tool allows developers to upload files to Google Cloud Storage (GCS) 
and get the content of GCS bucket, essentially to confirm that the upload was succesful.

Example Usage:

- Upload a file to a specific GCS path:
gcscli upload -p <object path in the bucket> -o <file name in the bucket> -f <path to file> bucket-name

- get files attributes in a GCS bucket or path if it exists:
gcscli get -bucketname <bucket-name> -o <target-path>

Powered by Cobra for efficient command-line interfaces.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gcs-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
