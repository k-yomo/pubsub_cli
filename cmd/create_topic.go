package cmd

import (
	"context"
	"github.com/k-yomo/pubsub_cli/pkg"
	"github.com/mitchellh/colorstring"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
)

// newCreateTopicCmd returns the command to create topics
func newCreateTopicCmd(out io.Writer) *cobra.Command {
	command := &cobra.Command{
		Use:     "create_topic TOPIC_ID",
		Short:   "create topic",
		Long:    "create Pub/Sub topic",
		Example: "pubsub_cli create_topic topic_1 topic_2 --host=localhost:8085 --project=test_project",
		Aliases: []string{"ct"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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
			return createTopic(cmd.Context(), out, pubsubClient, args)
		},
	}
	command.SetOut(out)
	return command
}

// createTopic finds or creates Pub/Sub topics
func createTopic(ctx context.Context, out io.Writer, pubsubClient *pkg.PubSubClient, topicIDs []string) error {
	for _, topicID := range topicIDs {
		topic, created, err := pubsubClient.FindOrCreateTopic(ctx, topicID)
		if err != nil {
			return errors.Wrapf(err, "find or create topic %s", topicID)
		}
		if created {
			_, _ = colorstring.Fprintf(out, "[green][success] topic '%s' created \n", topic.String())
		} else {
			_, _ = colorstring.Fprintf(out, "[cyan][skip] topic '%s' already exists \n", topic.String())
		}
		return nil
	}
	return nil
}
