package cmd

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/k-yomo/pubsub_cli/pkg"
	"github.com/mitchellh/colorstring"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
)

// newPublishCmd returns the command to publish message
func newPublishCmd(out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:     "publish TOPIC_ID DATA",
		Short:   "publish Pub/Sub message",
		Long:    "publish new message to given topic with given data",
		Example: "pubsub_cli publish test_topic '{\"key\":\"value\"}' --host=localhost:8085 --project=test_project",
		Aliases: []string{"p"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID := args[0]
			data := args[1]
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
			return publish(cmd.Context(), out, pubsubClient, topicID, data)
		},
	}
}

// publish publishes Pub/Sub message
func publish(ctx context.Context, out io.Writer, pubsubClient *pkg.PubSubClient, topicID, data string) error {
	topic, err := pubsubClient.FindOrCreateTopic(ctx, topicID)
	if err != nil {
		return errors.Wrapf(err, "find or create topic %s", topicID)
	}

	_, _ = colorstring.Fprintf(out, "[start] publishing message to %s...\n", topic.String())
	messageID, err := topic.Publish(ctx, &pubsub.Message{Data: []byte(data)}).Get(ctx)
	if err != nil {
		return errors.Wrapf(err, "publish message with data = %s", data)
	}
	_, _ = colorstring.Fprintf(out, "[green][success] published message to %s successfully, message ID = %s\n", topic.String(), messageID)
	return nil
}
