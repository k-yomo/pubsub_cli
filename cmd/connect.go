package cmd

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/k-yomo/pubsub_cli/pkg"
	"github.com/mitchellh/colorstring"
	"github.com/pkg/errors"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// newConnectCmd returns the command to connect topic
func newConnectCmd(out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "connect PROJECT_ID TOPIC_ID",
		Short: "connect remote topic to local topic",
		Long: `Connect subscribes Pub/Sub topic on GCP and publish got data to local topic on Pub/Sub emulator.
This command is useful when you want to make local push subscription subscribe Pub/Sub topic on GCP.
You need to be authenticated to subscribe the topic on GCP in some way listed in README and also need to set local emulator host either from env variable or from --host option.
`,
		Example: "pubsub_cli connect gcp_project test_topic --host=localhost:8085 --project=dev",
		Aliases: []string{"c"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			remoteProjectID := args[0]
			topicID := args[1]
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

			localPubsubClient, err := pkg.NewPubSubClient(ctx, projectID, emulatorHost, "")
			if err != nil {
				return errors.Wrap(err, "initialize pubsub emulator client")
			}
			// To avoid to connect local emulator, unset PUBSUB_EMULATOR_HOST explicitly
			if err := os.Unsetenv("PUBSUB_EMULATOR_HOST"); err != nil {
				return errors.Wrap(err, "unset PUBSUB_EMULATOR_HOST")
			}
			remotePubsubClient, err := pkg.NewPubSubClient(ctx, remoteProjectID, "", gcpCredentialFilePath)
			if err != nil {
				return errors.Wrap(err, "initialize pubsub client")
			}

			return connect(ctx, out, remotePubsubClient, localPubsubClient, topicID)
		},
	}
}

// connect connects remote Pub/Sub topic to local Pub/Sub topic
func connect(ctx context.Context, out io.Writer, remotePubsubClient, localPubsubClient *pkg.PubSubClient, topicID string) error {
	remoteTopic, err := remotePubsubClient.FindOrCreateTopic(ctx, topicID)
	if err != nil {
		return errors.Wrapf(err, "find or create remote topic %s", topicID)
	}
	localTopic, err := localPubsubClient.FindOrCreateTopic(ctx, topicID)
	if err != nil {
		return errors.Wrapf(err, "find or create local topic %s", topicID)
	}

	fmt.Println(fmt.Sprintf("[start]creating unique subscription to %s...", remoteTopic.String()))
	sub, err := remotePubsubClient.CreateUniqueSubscription(ctx, remoteTopic)
	if err != nil {
		return errors.Wrapf(err, "create unique subscription to %s", remoteTopic.String())
	}
	_, _ = colorstring.Fprintln(out, fmt.Sprintf("[green][success] topic %s is now connected!\n", topicID))

	fmt.Println("[start] waiting for publish...")
	err = sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
		_, _ = colorstring.Println(fmt.Sprintf("[green][success] Got message: %q", string(msg.Data)))
		messageID, err := localTopic.Publish(ctx, msg).Get(ctx)
		if err != nil {
			_, _ = colorstring.Fprintln(out, fmt.Sprintf("[red][error] publish message with data = %s", msg.Data))
			return
		}
		_, _ = colorstring.Fprintln(out, fmt.Sprintf("[green][success] published message to %s successfully, message ID = %s\n", localTopic.String(), messageID))
		msg.Ack()
	})
	if err != nil {
		return errors.Wrapf(err, "subscribe %s failed", remoteTopic.String())
	}
	return nil
}
