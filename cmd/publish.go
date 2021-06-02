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
	"strings"
)

const attributeFlagName = "attribute"

// newPublishCmd returns the command to publish message
func newPublishCmd(out io.Writer) *cobra.Command {
	command := &cobra.Command{
		Use:     "publish TOPIC_ID DATA",
		Short:   "publish Pub/Sub message",
		Long:    "publish new message to given topic with given data",
		Example: "pubsub_cli publish test_topic '{\"key\":\"value\"}' -a key=value --host=localhost:8085 --project=test_project",
		Aliases: []string{"p"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID := args[0]
			data := args[1]
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
			attrFlags, err := cmd.Flags().GetStringArray(attributeFlagName)
			if err != nil {
				return err
			}
			attrs := make(map[string]string, len(attrFlags))
			for _, attr := range attrFlags {
				kv := strings.Split(strings.TrimSpace(attr), "=")
				if len(kv) != 2 {
					return fmt.Errorf("attribute '%s' is invalid format. attribute format must be 'key=value'", attr)
				}
				attrs[kv[0]] = kv[1]
			}

			pubsubClient, err := pkg.NewPubSubClient(cmd.Context(), projectID, emulatorHost, gcpCredentialFilePath)
			if err != nil {
				return errors.Wrap(err, "initialize pubsub client")
			}
			m := &pubsub.Message{
				Attributes: attrs,
				Data:       []byte(data),
			}
			return publish(cmd.Context(), out, pubsubClient, topicID, m, createTopicIfNotExist)
		},
	}
	command.PersistentFlags().Bool(createTopicIfNotExistFlagName, false, "create topics if not exist")
	command.PersistentFlags().StringArrayP(attributeFlagName, "a", []string{}, "pubsub attribute")
	return command
}

// publish publishes Pub/Sub message
func publish(ctx context.Context, out io.Writer, pubsubClient *pkg.PubSubClient, topicID string, m *pubsub.Message, createTopicIfNotExist bool) error {
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
		return errors.Errorf("topic %s is not found", topicID)
	}

	_, _ = colorstring.Fprintf(out, "[start] publishing message to %s...\n", topic.String())
	messageID, err := topic.Publish(ctx, m).Get(ctx)
	if err != nil {
		return errors.Wrapf(err, "publish message %+v", m)
	}
	_, _ = colorstring.Fprintf(out, "[green][success] published message to %s successfully, message ID = %s\n", topic.String(), messageID)
	return nil
}
