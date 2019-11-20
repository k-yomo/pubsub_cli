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

// publishCmd represents the command to publish message
var publishCmd = &cobra.Command{
	Use:   "publish TOPIC_ID DATA",
	Short: "publish Pub/Sub message",
	Long:  "publish new message to given topic with given data",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		topicID := args[0]
		data := args[1]
		client, err := util.NewPubSubClient(projectID, emulatorHost)
		if err != nil {
			return errors.Wrap(err, "initialize pubsub client failed")
		}
		topic, err := client.FindOrCreateTopic(context.Background(), topicID)
		if err != nil {
			return errors.Wrapf(err, "find or create topic %s failed")
		}
		_, _ = colorstring.Println(fmt.Sprintf("[start] publish message to %s", topic.String()))
		messageID, err := topic.Publish(context.Background(), &pubsub.Message{Data: []byte(data)}).Get(context.Background())
		if err != nil {
			return errors.Wrapf(err, "publish message with data = %s failed", data)
		}
		_, _ = colorstring.Println(fmt.Sprintf("[green][success] published message to %s successfully, message ID = %s", topic.String(), messageID))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)
}
