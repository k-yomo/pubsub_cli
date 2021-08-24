package cmd

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/k-yomo/pubsub_cli/pkg"
	"github.com/mitchellh/colorstring"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"strconv"
	"time"
)

// newRegisterPushCmd returns the command to register an endpoint for subscribing
func newRegisterPushCmd(out io.Writer) *cobra.Command {
	command := &cobra.Command{
		Use:     "register_push TOPIC_ID ENDPOINT",
		Short:   "register Pub/Sub push endpoint",
		Long:    "register new endpoint for  push http request from Pub/Sub",
		Example: "pubsub_cli register_push test_topic http://localhost:1323/createSubscription --host=localhost:8085 --project=test_project",
		Aliases: []string{"r"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID := args[0]
			endpoint := args[1]
			ackDeadline, err := cmd.Flags().GetString("timeout")
			ackDeadlineNum, err := strconv.ParseInt(ackDeadline, 10, 64)
			projectID, err := cmd.Flags().GetString("project")
			emulatorHost, err := cmd.Flags().GetString("host")
			gcpCredentialFilePath, err := cmd.Flags().GetString("cred-file")

			pubsubClient, err := pkg.NewPubSubClient(cmd.Context(), projectID, emulatorHost, gcpCredentialFilePath)
			if err != nil {
				return errors.Wrap(err, "initialize pubsub client")
			}
			return registerPush(cmd.Context(), out, pubsubClient, topicID, endpoint, ackDeadlineNum)
		},
	}
	command.SetOut(out)
	return command
}

// registerPush registers new push endpoint
func registerPush(ctx context.Context, out io.Writer, pubsubClient *pkg.PubSubClient, topicID, endpoint string, ackDeadline int64) error {
	topic, err := pubsubClient.FindOrCreateTopic(ctx, topicID)
	if err != nil {
		return errors.Wrapf(err, "[error]find or create topic %s", topicID)
	}

	_, _ = colorstring.Fprintf(out, "[start] registering push endpoint for %s...\n", topic.String())
	subscriptionConfig := pubsub.SubscriptionConfig{
		Topic:            topic,
		AckDeadline: time.Duration(ackDeadline) * time.Second,
		//RetryPolicy: &pubsub.RetryPolicy{
		//	MinimumBackoff: time.Duration(ackDeadline) * time.Second,
		//	MaximumBackoff: time.Duration(ackDeadline) * time.Second,
		//},
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
