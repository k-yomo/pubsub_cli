package cmd

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/k-yomo/pubsub_cli/util"
	"github.com/mitchellh/colorstring"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// subscribeCmd represents the command to subscribe messages
var subscribeCmd = &cobra.Command{
	Use:   "subscribe TOPIC_ID ...",
	Short: "subscribe Pub/Sub topic",
	Long:  "create subscription for given Pub/Sub topic and subscribe the topic",
	Args:  cobra.MinimumNArgs(1),
	RunE:  subscribe,
}

func init() {
	rootCmd.AddCommand(subscribeCmd)
}

type subscriber struct {
	topic *pubsub.Topic
	sub   *pubsub.Subscription
}

func subscribe(_ *cobra.Command, args []string) error {
	ctx := context.Background()
	topicIDs := args
	subscribers := make([]*subscriber, len(topicIDs))
	for i, topicID := range topicIDs {
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
		subscribers[i] = &subscriber{topic: topic, sub: sub}
		_, _ = colorstring.Println(fmt.Sprintf("[green][success] created subscription to %s", topic.String()))
	}

	fmt.Println("[start] waiting for publish...")
	wg := sync.WaitGroup{}
	for _, sub := range subscribers {
		wg.Add(1)
		go func(s *subscriber) {
			defer wg.Done()
			err := s.sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
				msg.Ack()
				_, _ = colorstring.Println(fmt.Sprintf("[green][success] got message to %s, id: %s, data: %q", s.topic.ID(), msg.ID, string(msg.Data)))
			})
			if err != nil {
				_, _ = colorstring.Println(fmt.Sprintf("[red][error] %s has got error: %v", s.sub.ID(), err))
			}
		}(sub)
	}
	wg.Wait()
	return nil
}
