/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// listCmd represents the list command
var getCmd = &cobra.Command{
	Use:   "get [bucket_name] and attributes",
	Short: "get attributes of zimvest GCS bucket object if it exists example: [gcscli --objectname 'path/to/object.pdf' 'bucketname']",
	Long: `get attributes of zimvest GCS bucket object if it exists 
	for example: 
[gcscli -o '<path/to/object.pdf>' -b '<bucketname>']
- -o: object that is being serched for including path prefix.
- -b: bucketname.`,
	// Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// bucketName := args[1]
		bucketName, _ := cmd.Flags().GetString("bucketname")
		objectName, _ := cmd.Flags().GetString("objectname")

		if bucketName != "" {
			if err := listBucketContents(bucketName, objectName); err != nil {
				log.Fatalf("Failed to get object: %v", err)
			}
		} else {
			if err := listBucketContents("prod-eu-zimvest", objectName); err != nil {
				log.Fatalf("Failed to get object: %v", err)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	getCmd.Flags().StringP("bucketname", "b", "prod-eu-zimvest", "gcs bucket")
	getCmd.Flags().StringP("objectname", "o", "", "object to return inclusing any path prefix")
}

func listBucketContents(bucketName, objectName string) error {
	// Read the embedded key.json file
	keyData, err := keyFile.ReadFile("assets/key.json")
	if err != nil {
		return fmt.Errorf("failed to read embedded key.json: %v", err)
	}

	// Create a credentials object using the key data
	creds, err := google.CredentialsFromJSON(context.Background(), keyData, storage.ScopeFullControl)
	if err != nil {
		return fmt.Errorf("failed to create credentials from JSON: %v", err)
	}

	// Set a timeout for the context to ensure it doesn't hang indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	client, err := storage.NewClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	fmt.Printf("Checking if object exists in bucket: %s with name: %s\n", bucketName, objectName)
	object := client.Bucket(bucketName).Object(objectName)

	// Get the object's attributes (metadata)
	attrs, err := object.Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			fmt.Printf("Object %s does not exist in bucket %s\n", objectName, bucketName)
			return nil
		}
		return fmt.Errorf("Object(%q).Attrs: %v", objectName, err)
	}

	// If the object exists, print some of its attributes
	fmt.Printf("Object found:\n")
	fmt.Printf("Name: %s\n", attrs.Name)
	fmt.Printf("Content-Type: %s\n", attrs.ContentType)
	fmt.Printf("Last Modified: %s\n", attrs.Updated)
	fmt.Printf("Storage Class: %s\n", attrs.StorageClass)

	return nil
}
