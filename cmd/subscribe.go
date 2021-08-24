package cmd

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/k-yomo/pubsub_cli/pkg"
	"github.com/mitchellh/colorstring"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"strconv"
	"time"
)

// newSubscribeCmd returns the command to subscribe messages
func newSubscribeCmd(out io.Writer) *cobra.Command {
	command :=  &cobra.Command{
		Use:     "subscribe TOPIC_ID ...",
		Short:   "subscribe Pub/Sub topics",
		Long:    "create temporary subscriptions for given Pub/Sub topics(or you can set 'all' to subscribe all topics) and subscribe the topics",
		Example: "pubsub_cli subscribe test_topic another_topic --create-if-not-exist --host=localhost:8085 --project=test_project",
		Aliases: []string{"s"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicIDs := args
			ackDeadline, _ := cmd.Flags().GetInt(ackDeadlineFlagName)
			ackDeadlineSecond := time.Duration(ackDeadline) * time.Second
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
			createTopicIfNotExist, err := cmd.Flags().GetBool(createTopicIfNotExistFlagName)
			if err != nil {
				return err
			}

			pubsubClient, err := pkg.NewPubSubClient(cmd.Context(), projectID, emulatorHost, gcpCredentialFilePath)
			if err != nil {
				return errors.Wrap(err, "initialize pubsub client")
			}
			return subscribe(cmd.Context(), out, pubsubClient, topicIDs, createTopicIfNotExist, ackDeadlineSecond)
		},
	}
	command.SetOut(out)
	command.PersistentFlags().Bool(createTopicIfNotExistFlagName, false, "create topics if not exist")
	ackDeadlineDefault, _ := strconv.Atoi(os.Getenv("PUBSUB_ACK_DEADLINE"))
	command.PersistentFlags().IntVarP(&ackDeadlineDefault, ackDeadlineFlagName, "a", ackDeadlineDefault, "pubsub ack deadline(unit seconds)")
	return command
}

type subscriber struct {
	topic *pubsub.Topic
	sub   *pubsub.Subscription
}

// subscribe subscribes Pub/Sub messages
func subscribe(ctx context.Context, out io.Writer, pubsubClient *pkg.PubSubClient, topicIDs []string, createTopicIfNotExist bool, ackDeadline time.Duration) error {
	var topics []*pubsub.Topic
	var err error
	// if topic name is "all", subscribe all topics in the project
	if topicIDs[0] == "all" {
		topics, err = pubsubClient.FindAllTopics(ctx)
	} else {
		if createTopicIfNotExist {
			topics, err = pubsubClient.FindOrCreateTopics(ctx, topicIDs)
		} else {
			topics, err = pubsubClient.FindTopics(ctx, topicIDs)
		}
	}
	if err != nil {
		return err
	}
	if len(topics) == 0 {
		return errors.Errorf("topics %v not found", topicIDs)
	}
	if len(topics) != len(topicIDs) {
		topicIDExistMap := map[string]bool{}
		for _, topic := range topics {
			topicIDExistMap[topic.ID()] = true
		}

		for _, topicID := range topicIDs {
			if !topicIDExistMap[topicID] {
				_, _ = colorstring.Fprintf(out, "[yellow][warn] topic %s is not found\n", topicID)
			}
		}
		_, _ = colorstring.Fprintln(out, "[yellow][warn] if you want to create topics if not exist, set --create-if-not-exist flag")
	}

	eg := &errgroup.Group{}
	subscribers := make(chan *subscriber, len(topics))
	for _, topic := range topics {
		topic := topic
		eg.Go(func() error {
			fmt.Fprintf(out, "[start] creating unique subscription to %s...\n", topic.String())
			sub, err := pubsubClient.CreateUniqueSubscription(ctx, topic, ackDeadline)
			if err != nil {
				return errors.Wrapf(err, "create unique subscription to %s", topic.String())
			}
			subscribers <- &subscriber{topic: topic, sub: sub}
			_, _ = colorstring.Fprintf(out, "[green][success] created subscription '%s' to %s\n",  sub.ID(), topic.String())
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	close(subscribers)

	_, _ = fmt.Fprintln(out, "[start] waiting for publish...")
	for s := range subscribers {
		s := s
		eg.Go(func() error {
			err := s.sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
				msg.Ack()
				_, _ = colorstring.Fprintf(out, "[green][success] got message published to %s, id: %s, data: %q\n", s.topic.ID(), msg.ID, string(msg.Data))
			})
			return errors.Wrapf(err, "receive message published to %s through %s subscription", s.topic.ID(), s.sub.ID())
		})
	}
	return eg.Wait()
}
