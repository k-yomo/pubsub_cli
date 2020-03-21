package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/colorstring"
	"github.com/spf13/cobra"
)

var version string
var projectID string
var emulatorHost string
var gcpCredentialFilePath string

// Exec executes command
func Exec() {
	rootCmd := newRootCmd()
	if err := rootCmd.Execute(); err != nil {
		_, _ = colorstring.Println(fmt.Sprintf("[red][error]%v", err))
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:     "pubsub_cli",
		Short:   "pubsub_cli is a handy cloud Pub/Sub CLI",
		Long:    "Very simple cloud Pub/Sub CLI used as publisher / subscriber",
		Version: version,
	}

	projectID = os.Getenv("GCP_PROJECT_ID")
	emulatorHost = os.Getenv("PUBSUB_EMULATOR_HOST")
	gcpCredentialFilePath = os.Getenv("GCP_CREDENTIAL_FILE_PATH")
	rootCmd.PersistentFlags().Bool("help", false, fmt.Sprintf("help for %s", rootCmd.Name()))
	rootCmd.PersistentFlags().StringVarP(&projectID, "project", "p", projectID, "gcp project id (You can also set 'GCP_PROJECT_ID' to env variable)")
	rootCmd.PersistentFlags().StringVarP(&emulatorHost, "host", "h", emulatorHost, "emulator host (You can also set 'PUBSUB_EMULATOR_HOST' to env variable)")
	rootCmd.PersistentFlags().StringVarP(&gcpCredentialFilePath, "cred-file", "c", gcpCredentialFilePath, "gcp credential file path (You can also set 'GCP_CREDENTIAL_FILE_PATH' to env variable)")

	rootCmd.AddCommand(newPublishCmd(os.Stdin), newSubscribeCmd(os.Stdin), newRegisterPushCmd(os.Stdin))
	return rootCmd
}
