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

// subscribeCmd represents the subscribe command
var subscribeCmd = &cobra.Command{
	Use:   "subscribe TOPIC_ID",
	Short: "subscribe Pub/Sub topic",
	Long:  "create subscription for given Pub/Sub topic and subscribe the topic",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		topicID := args[0]
		client, err := util.NewPubSubClient(projectID, emulatorHost)
		if err != nil {
			return errors.Wrap(err, "initialize pubsub client failed")
		}
		topic, err := client.FindOrCreateTopic(context.Background(), topicID)
		if err != nil {
			return errors.Wrapf(err, "find or create topic %s failed")
		}
		fmt.Println(fmt.Sprintf("[start]create unique subscription to %s", topic.String()))
		sub, err := client.CreateUniqueSubscription(topic)
		if err != nil {
			return errors.Wrapf(err, "create unique subscription to %s failed", topic.String())
		}
		fmt.Println(fmt.Sprintf("create unique subscription to %s", topic.String()))
		fmt.Println("[start] waiting for publish...")
		err = sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
			msg.Ack()
			_, _ = colorstring.Println(fmt.Sprintf("[green][success] Got message: %q\n", string(msg.Data)))
		})
		if err != nil {
			return errors.Wrapf(err, "subscribe %s failed", topic.String())
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(subscribeCmd)
}
