package cmd

import (
	"context"
	"fmt"
	"io"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/k-yomo/pubsub_cli/util"
	"github.com/mitchellh/colorstring"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// newSubscribeCmd returns the command to subscribe messages
func newSubscribeCmd(out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:     "subscribe TOPIC_ID ...",
		Short:   "subscribe Pub/Sub topic",
		Long:    "create subscription for given Pub/Sub topic and subscribe the topic",
		Example: "pubsub_cli subscribe test_topic another_topic --host=localhost:8085 --project=test_project",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pubsubClient, err := util.NewPubSubClient(context.Background(), projectID, emulatorHost, gcpCredentialFilePath)
			if err != nil {
				return errors.Wrap(err, "initialize pubsub client")
			}
			return subscribe(cmd, out, pubsubClient, args)
		},
	}
}

type subscriber struct {
	topic *pubsub.Topic
	sub   *pubsub.Subscription
}

// subscribe subscribes Pub/Sub messages
func subscribe(_ *cobra.Command, out io.Writer, pubsubClient *util.PubSubClient, args []string) error {
	ctx := context.Background()
	topicIDs := args
	subscribers := make([]*subscriber, len(topicIDs))
	for i, topicID := range topicIDs {
		topic, err := pubsubClient.FindOrCreateTopic(ctx, topicID)
		if err != nil {
			return errors.Wrapf(err, "find or create topic %s", topicID)
		}

		fmt.Println(fmt.Sprintf("[start]creating unique subscription to %s...", topic.String()))
		sub, err := pubsubClient.CreateUniqueSubscription(ctx, topic)
		if err != nil {
			return errors.Wrapf(err, "create unique subscription to %s", topic.String())
		}
		subscribers[i] = &subscriber{topic: topic, sub: sub}
		_, _ = colorstring.Fprintln(out, fmt.Sprintf("[green][success] created subscription to %s", topic.String()))
	}

	_, _ = fmt.Fprintln(out, "[start] waiting for publish...")
	wg := sync.WaitGroup{}
	for _, sub := range subscribers {
		wg.Add(1)
		go func(s *subscriber) {
			defer wg.Done()
			err := s.sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
				msg.Ack()
				_, _ = colorstring.Fprintln(out, fmt.Sprintf("[green][success] got message to %s, id: %s, data: %q", s.topic.ID(), msg.ID, string(msg.Data)))
			})
			if err != nil {
				_, _ = colorstring.Fprintln(out, fmt.Sprintf("[red][error] %s has got error: %v", s.sub.ID(), err))
			}
		}(sub)
	}
	wg.Wait()
	return nil
}
