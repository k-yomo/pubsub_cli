package cmd

import (
	"fmt"
	"github.com/mitchellh/colorstring"
	"github.com/spf13/cobra"
	"os"
)

var projectID string
var emulatorHost string
var gcpCredentialFilePath string

func init() {
	projectID = os.Getenv("GCP_PROJECT_ID")
	emulatorHost = os.Getenv("PUBSUB_EMULATOR_HOST")
	gcpCredentialFilePath = os.Getenv("GCP_CREDENTIAL_FILE_PATH")
}

var rootCmd = &cobra.Command{
	Use:   "pubsub_cli",
	Short: "pubsub_cli is a handy cloud Pub/Sub CLI",
	Long:  "Very simple cloud Pub/Sub CLI used as publisher / subscriber",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Exec() {
	rootCmd.PersistentFlags().StringVar(&projectID, "project", projectID, "gcp project id (You can also set 'GCP_PROJECT_ID' to env variable)")
	rootCmd.PersistentFlags().StringVar(&emulatorHost, "host", emulatorHost, "emulator host (You can also set 'PUBSUB_EMULATOR_HOST' to env variable)")
	rootCmd.PersistentFlags().StringVar(&gcpCredentialFilePath, "cred-file", gcpCredentialFilePath, "gcp credential file path (You can also set 'GCP_CREDENTIAL_FILE_PATH' to env variable)")
	if err := rootCmd.Execute(); err != nil {
		_, _ = colorstring.Println(fmt.Sprintf("[red][error]%v", err))
		os.Exit(1)
	}
}
