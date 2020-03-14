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
	"time"
)

// newRegisterPushCmd returns the command to register an endpoint for subscribing
func newRegisterPushCmd(out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:     "register_push TOPIC_ID ENDPOINT",
		Short:   "register Pub/Sub push endpoint",
		Long:    "register new endpoint for  push http request from Pub/Sub",
		Example: "pubsub_cli register_push test_topic http://localhost:1323/subscribe --host=localhost:8085 --project=test_project",
		Aliases: []string{"r"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			pubsubClient, err := util.NewPubSubClient(context.Background(), projectID, emulatorHost, gcpCredentialFilePath)
			if err != nil {
				return errors.Wrap(err, "initialize pubsub client")
			}
			return registerPush(cmd, out, pubsubClient, args)
		},
	}
}

// registerPush registers new push endpoint
func registerPush(_ *cobra.Command, out io.Writer, pubsubClient *util.PubSubClient, args []string) error {
	ctx := context.Background()
	topicID := args[0]
	endpoint := args[1]

	topic, err := pubsubClient.FindOrCreateTopic(ctx, topicID)
	if err != nil {
		return errors.Wrapf(err, "[error]find or create topic %s", topicID)
	}

	_, _ = colorstring.Fprintln(out, fmt.Sprintf("[start] registering push endpoint for %s...", topic.String()))
	subscriptionConfig := pubsub.SubscriptionConfig{
		Topic:            topic,
		ExpirationPolicy: time.Hour * 24,
		PushConfig: pubsub.PushConfig{
			Endpoint:             endpoint,
			Attributes:           nil,
			AuthenticationMethod: nil,
		},
	}
	if _, err := pubsubClient.CreateSubscription(context.Background(), util.UUID(), subscriptionConfig); err != nil {
		return errors.Wrapf(err, "register push endpoint for = %s", topic.String())
	}
	_, _ = colorstring.Fprintln(out, fmt.Sprintf("[green][success] registered %s as an endpoint for %s", endpoint, topic.String()))
	return nil
}
