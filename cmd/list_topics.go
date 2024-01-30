package cmd

import (
	"context"
	"io"

	"github.com/k-yomo/pubsub_cli/pkg"
	"github.com/mitchellh/colorstring"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// newCreateTopicCmd returns the command to create topics
func newListTopicsCmd(out io.Writer) *cobra.Command {
	command := &cobra.Command{
		Use:     "list",
		Short:   "lists all Pub/Sub topics in the given project",
		Long:    "lists all Pub/Sub topics in the given project",
		Example: "pubsub_cli list",
		Aliases: []string{"ls"},
		Args:    cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			projectID, err := cmd.Flags().GetString(projectFlagName)
			if err != nil {
				return err
			}
			emulatorHost, err := cmd.Flags().GetString(hostFlagName)
			if err != nil {
				return err
			}
			gcpCredentialFilePath, err := cmd.Flags().GetString(credFileFlagName)
			if err != nil {
				return err
			}
			pubsubClient, err := pkg.NewPubSubClient(cmd.Context(), projectID, emulatorHost, gcpCredentialFilePath)
			if err != nil {
				return errors.Wrap(err, "initialize pubsub client")
			}
			return listTopics(cmd.Context(), out, pubsubClient)
		},
	}
	command.SetOut(out)
	return command
}

// listTopics lists all Pub/Sub topics in the project
func listTopics(ctx context.Context, out io.Writer, pubsubClient *pkg.PubSubClient) error {
	topics, err := pubsubClient.FindAllTopics(ctx)
	if err != nil {
		return errors.Wrapf(err, "list topics")
	}
	for _, topic := range topics {
		_, _ = colorstring.Fprintf(out, "[green]name: '%s'\n", topic.String())
	}
	return nil
}
