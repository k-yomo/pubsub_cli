package cmd

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/k-yomo/pubsub_cli/util"
	"github.com/mitchellh/colorstring"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// subscribeCmd represents the command to subscribe messages
var subscribeCmd = &cobra.Command{
	Use:   "subscribe TOPIC_ID",
	Short: "subscribe Pub/Sub topic",
	Long:  "create subscription for given Pub/Sub topic and subscribe the topic",
	Args:  cobra.ExactArgs(1),
	RunE:  subscribe,
}

func init() {
	rootCmd.AddCommand(subscribeCmd)
}

func subscribe(_ *cobra.Command, args []string) error {
	ctx := context.Background()
	topicID := args[0]

	client, err := util.NewPubSubClient(ctx, projectID, emulatorHost, gcpCredentialFilePath)
	if err != nil {
		return errors.Wrap(err, "initialize pubsub client")
	}
	topic, err := client.FindOrCreateTopic(ctx, topicID)
	if err != nil {
		return errors.Wrapf(err, "find or create topic %s")
	}

	fmt.Println(fmt.Sprintf("[start]creating unique subscription to %s...", topic.String()))
	sub, err := client.CreateUniqueSubscription(ctx, topic)
	if err != nil {
		return errors.Wrapf(err, "create unique subscription to %s", topic.String())
	}
	_, _ = colorstring.Println(fmt.Sprintf("[green][success] created subscription to %s", topic.String()))

	fmt.Println("[start] waiting for publish...")
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		_, _ = colorstring.Println(fmt.Sprintf("[green][success] Got message id: %s, data: %q\n", msg.ID, string(msg.Data)))
	})
	if err != nil {
		return errors.Wrapf(err, "subscribe %s failed", topic.String())
	}
	return nil
}
