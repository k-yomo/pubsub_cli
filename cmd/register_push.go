package cmd

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/k-yomo/pubsub_cli/pkg"
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
			topicID := args[0]
			endpoint := args[1]
			projectID, err := cmd.Flags().GetString("project")
			emulatorHost, err := cmd.Flags().GetString("host")
			gcpCredentialFilePath, err := cmd.Flags().GetString("cred-file")

			pubsubClient, err := pkg.NewPubSubClient(cmd.Context(), projectID, emulatorHost, gcpCredentialFilePath)
			if err != nil {
				return errors.Wrap(err, "initialize pubsub client")
			}
			return registerPush(cmd.Context(), out, pubsubClient, topicID, endpoint)
		},
	}
}

// registerPush registers new push endpoint
func registerPush(ctx context.Context, out io.Writer, pubsubClient *pkg.PubSubClient, topicID, endpoint string) error {
	topic, err := pubsubClient.FindOrCreateTopic(ctx, topicID)
	if err != nil {
		return errors.Wrapf(err, "[error]find or create topic %s", topicID)
	}

	_, _ = colorstring.Fprintf(out, "[start] registering push endpoint for %s...\n", topic.String())
	subscriptionConfig := pubsub.SubscriptionConfig{
		Topic:            topic,
		ExpirationPolicy: 24 * time.Hour,
		PushConfig: pubsub.PushConfig{
			Endpoint:             endpoint,
			Attributes:           nil,
			AuthenticationMethod: nil,
		},
	}
	if _, err := pubsubClient.CreateSubscription(context.Background(), pkg.UUID(), subscriptionConfig); err != nil {
		return errors.Wrapf(err, "register push endpoint for = %s", topic.String())
	}
	_, _ = colorstring.Fprintf(out, "[green][success] registered %s as an endpoint for %s\n", endpoint, topic.String())
	return nil
}
