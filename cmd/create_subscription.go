package cmd

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/k-yomo/pubsub_cli/pkg"
	"github.com/mitchellh/colorstring"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
)

// newCreateSubscriptionCmd returns the command to createSubscription messages
func newCreateSubscriptionCmd(out io.Writer) *cobra.Command {
	command := &cobra.Command{
		Use:     "create_subscription TOPIC_ID SUBSCRIPTION_ID",
		Short:   "create Pub/Sub subscription",
		Long:    "create subscription",
		Example: "pubsub_cli createSubscription test_topic test_topic_sub --create-if-not-exist --host=localhost:8085 --project=test_project",
		Aliases: []string{"cs"},
		Args:    cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID := args[0]
			subscriptionID := args[1]
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
			return createSubscription(cmd.Context(), out, pubsubClient, topicID, subscriptionID, createTopicIfNotExist)
		},
	}
	command.SetOut(out)
	command.PersistentFlags().Bool(createTopicIfNotExistFlagName, false, "create topics if not exist")
	// TODO: add flags below
	//   optional flags may be  --ack-deadline | --dead-letter-topic |
	//                         --dead-letter-topic-project |
	//                         --enable-message-ordering | --expiration-period |
	//                         --help | --labels | --max-delivery-attempts |
	//                         --max-retry-delay | --message-retention-duration |
	//                         --min-retry-delay | --push-auth-service-account |
	//                         --push-auth-token-audience | --push-endpoint |
	//                         --retain-acked-messages | --topic-project
	return command
}

// createSubscription create Pub/Sub subscription
func createSubscription(ctx context.Context, out io.Writer, pubsubClient *pkg.PubSubClient, topicID, subscriptionID string, createTopicIfNotExist bool) error {
	var topic *pubsub.Topic
	var err error
	if createTopicIfNotExist {
		topic, err = pubsubClient.FindOrCreateTopic(ctx, topicID)
	} else {
		topic, err = pubsubClient.FindTopic(ctx, topicID)
	}
	if err != nil {
		return err
	}
	if topic == nil {
		return errors.Errorf("topic %v not found", topicID)
	}

	ok, err := pubsubClient.Subscription(subscriptionID).Exists(ctx)
	if err != nil {
		return err
	}
	if ok {
		_, _ = colorstring.Fprintf(out, "[cyan][skip] subscription '%s' already exists\n", subscriptionID)
		return nil
	}

	fmt.Fprintf(out, "[start] creating subscription to %s...\n", topic.String())
	sub, err := pubsubClient.CreateSubscription(ctx, subscriptionID, pubsub.SubscriptionConfig{Topic: topic})
	if err != nil {
		return errors.Wrapf(err, "create unique subscription to %s", topic.String())
	}
	_, _ = colorstring.Fprintf(out, "[green][success] created subscription '%s' to the topic '%s'\n", sub.ID(), topic.String())
	return nil
}
