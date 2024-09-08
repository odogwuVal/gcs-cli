/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"embed"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// func init() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatalf("Error loading .env file")
// 	}
// }

var authorizedUsers []string

//go:embed assets/.env
var envFile embed.FS

//go:embed assets/key.json
var keyFile embed.FS

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload [file_path] [bucket_name]",
	Short: "Upload document to zimvest GCS bucket",
	Long: `Uploads documents to zimvest/zedcrestwealth bucket
	For example:

[gcscli upload -p <object path in the bucket> -o <file name in the bucket> -f <path to file> bucket-name]

- -p: Path in the bucket where the object will be uploaded.
- -o: Name of the object in the bucket.
- -f: File path on the local file system to upload.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Parse the embedded .env file
		loadEmbeddedEnv()
		// Get user token from environment variable
		userToken := os.Getenv("GCSCLI_USER_TOKEN")
		if userToken == "" {
			log.Fatal("Missing user token. Please set the environment variable 'GCSCLI_USER_TOKEN'.")
		}

		// Validate the user token
		if !validateToken(userToken) {
			log.Fatal("Unauthorized user. Your token is not in the list of authorized users.")
		}
	},

	Args: cobra.ExactArgs(1), //bucketname should be provided as an argument
	Run: func(cmd *cobra.Command, args []string) {
		bucketName := args[0]
		filePath, _ := cmd.Flags().GetString("filepath")
		objectPath, _ := cmd.Flags().GetString("objectpath")
		objectName, _ := cmd.Flags().GetString("objectname")

		if objectPath == "" || objectName == "" || filePath == "" {
			log.Fatal("You must provide the path, object name, and file path.")
		}

		// Combine object path and object name
		fullObjectPath := objectPath + "/" + objectName

		if err := uploadToGCS(os.Stdout, bucketName, fullObjectPath, filePath); err != nil {
			log.Fatalf("Failed to upload file: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	uploadCmd.Flags().StringP("filepath", "f", "", "Path to the file on the local filesystem")
	uploadCmd.Flags().StringP("objectname", "o", "", "Name of the object in the bucket")
	uploadCmd.Flags().StringP("objectpath", "p", "", "path to upload document in gcs bucket")
}

// Function to load the embedded .env file and parse its contents
func loadEmbeddedEnv() {
	data, err := envFile.ReadFile("assets/.env")
	if err != nil {
		log.Fatalf("Failed to read embedded .env file: %v", err)
	}

	// Parse the .env file line by line
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Fatalf("Invalid line in .env file: %s", line)
		}

		// Set the environment variable
		os.Setenv(parts[0], parts[1])

		// If the line contains a token, add it to the authorized tokens list
		if strings.HasPrefix(parts[0], "TOKEN") {
			authorizedUsers = append(authorizedUsers, parts[1])
		}
	}
}

func uploadToGCS(w io.Writer, bucketName, objectPath, filePath string) error {
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

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	// Get file size for progress tracking
	fileInfo, err := f.Stat()
	if err != nil {
		return fmt.Errorf("f.Stat: %v", err)
	}
	fileSize := fileInfo.Size()

	bar := progressbar.NewOptions64(fileSize, progressbar.OptionSetDescription("Uploading..."))

	// Set a timeout for the upload operation.
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Create an object handle with the full path (including the folder structure).
	o := client.Bucket(bucketName).Object(objectPath)

	// Optional: Add a condition to prevent overwriting if the object already exists.
	o = o.If(storage.Conditions{DoesNotExist: true})

	// Create a writer to upload the file.
	wc := o.NewWriter(ctx)

	// Create a multi-writer to update the progress bar while uploading
	writer := io.MultiWriter(wc, bar)

	if _, err = io.Copy(writer, f); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		// Check if the error is a googleapi error with a 412 status code
		if gErr, ok := err.(*googleapi.Error); ok && gErr.Code == 412 {
			fmt.Fprintln(w, "\nUpload failed: Precondition not met. The object might already exist.")
			return nil
		}

		return fmt.Errorf("Writer.Close: %w", err)
	}

	fmt.Fprintf(w, "\nFile %v uploaded to bucket %v at path %v.\n", filePath, bucketName, objectPath)
	return nil
}

// validateToken checks if the provided token or username is in the list of authorized users
func validateToken(userToken string) bool {
	// Fetch the list of authorized users from the environment variable
	// authorizedUsers := os.Getenv("GCSCLI_AUTHORIZED_USERS")
	// if authorizedUsers == "" {
	// 	log.Fatal("Missing authorized users list. Please set the environment variable 'GCSCLI_AUTHORIZED_USERS'.")
	// }

	// Check if the userToken is in the list of authorized users
	for _, authorizedUser := range authorizedUsers {
		if userToken == authorizedUser {
			return true
		}
	}
	return false
}
