package cmd

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/k-yomo/pubsub_cli/util"
	"github.com/mitchellh/colorstring"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/spf13/cobra"
	"time"
)

// newRegisterPushCmd returns the command to register an endpoint for subscribing
func newRegisterPushCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "register_push TOPIC_ID ENDPOINT",
		Short: "register Pub/Sub push endpoint",
		Long:  "register new endpoint for  push http request from Pub/Sub",
		Args:  cobra.ExactArgs(2),
		RunE:  registerPush,
	}
}

// registerPush registers new push endpoint
func registerPush(_ *cobra.Command, args []string) error {
	ctx := context.Background()
	topicID := args[0]
	endpoint := args[1]

	client, err := util.NewPubSubClient(ctx, projectID, emulatorHost, gcpCredentialFilePath)
	if err != nil {
		return errors.Wrap(err, "[error]initialize pubsub client")
	}
	topic, err := client.FindOrCreateTopic(ctx, topicID)
	if err != nil {
		return errors.Wrapf(err, "[error]find or create topic %s", topicID)
	}

	_, _ = colorstring.Println(fmt.Sprintf("[start] registering push endpoint for %s...", topic.String()))
	subscriptionConfig := pubsub.SubscriptionConfig{
		Topic:            topic,
		ExpirationPolicy: time.Hour * 24,
		PushConfig: pubsub.PushConfig{
			Endpoint:             endpoint,
			Attributes:           nil,
			AuthenticationMethod: nil,
		},
	}
	if _, err := client.CreateSubscription(context.Background(), xid.New().String(), subscriptionConfig); err != nil {
		return errors.Wrapf(err, "register push endpoint for = %s", topic.String())
	}
	_, _ = colorstring.Println(fmt.Sprintf("[green][success] registered %s as an endpoint for %s", endpoint, topic.String()))
	return nil
}
