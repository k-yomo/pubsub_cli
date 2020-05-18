package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/colorstring"
	"github.com/spf13/cobra"
)

const projectFlagName = "project"
const hostFlagName = "host"
const credFileFlagName = "cred-file"

var version string

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

	projectID := os.Getenv("GCP_PROJECT_ID")
	emulatorHost := os.Getenv("PUBSUB_EMULATOR_HOST")
	gcpCredentialFilePath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	rootCmd.PersistentFlags().Bool("help", false, fmt.Sprintf("help for %s", rootCmd.Name()))
	rootCmd.PersistentFlags().StringVarP(&projectID, projectFlagName, "p", projectID, "gcp project id (You can also set 'GCP_PROJECT_ID' to env variable)")
	rootCmd.PersistentFlags().StringVarP(&emulatorHost, hostFlagName, "h", emulatorHost, "emulator host (You can also set 'PUBSUB_EMULATOR_HOST' to env variable)")
	rootCmd.PersistentFlags().StringVarP(&gcpCredentialFilePath, credFileFlagName, "c", gcpCredentialFilePath, "gcp credential file path (You can also set 'GOOGLE_APPLICATION_CREDENTIALS' to env variable)")

	rootCmd.AddCommand(
		newPublishCmd(os.Stdin),
		newSubscribeCmd(os.Stdin),
		newRegisterPushCmd(os.Stdin),
	)
	return rootCmd
}
