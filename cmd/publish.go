package cmd

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/k-yomo/pubsub_cli/util"
	"github.com/mitchellh/colorstring"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
)

// newPublishCmd returns the command to publish message
func newPublishCmd(out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "publish TOPIC_ID DATA",
		Short: "publish Pub/Sub message",
		Long:  "publish new message to given topic with given data",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			pubsubClient, err := util.NewPubSubClient(context.Background(), projectID, emulatorHost, gcpCredentialFilePath)
			if err != nil {
				return errors.Wrap(err, "initialize pubsub client")
			}
			return publish(cmd, out, pubsubClient, args)
		},
	}
}

// publish publishes Pub/Sub message
func publish(_ *cobra.Command, out io.Writer, pubsubClient *util.PubSubClient, args []string) error {
	ctx := context.Background()
	topicID := args[0]
	data := args[1]

	topic, err := pubsubClient.FindOrCreateTopic(ctx, topicID)
	if err != nil {
		return errors.Wrapf(err, "find or create topic %s", topicID)
	}

	_, _ = colorstring.Fprintln(out, fmt.Sprintf("[start] publishing message to %s...", topic.String()))
	messageID, err := topic.Publish(ctx, &pubsub.Message{Data: []byte(data)}).Get(ctx)
	if err != nil {
		return errors.Wrapf(err, "publish message with data = %s", data)
	}
	_, _ = colorstring.Fprintln(out, fmt.Sprintf("[green][success] published message to %s successfully, message ID = %s", topic.String(), messageID))
	return nil
}
