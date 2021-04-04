package cmd

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/k-yomo/pubsub_cli/pkg"
	"github.com/mitchellh/colorstring"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

// newConnectCmd returns the command to connect topic
func newConnectCmd(out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "connect PROJECT_ID TOPIC_ID ...",
		Short: "connect remote topics to local topics",
		Long: `Connect subscribes Pub/Sub topics(or you can set 'all' to createSubscription all topics) on GCP and publish got data to local topics on Pub/Sub emulator.
This command is useful when you want to make local push subscription createSubscription Pub/Sub topic on GCP.
You need to be authenticated to createSubscription the topic on GCP in some way listed in README and also need to set local emulator host either from env variable or from --host option.
`,
		Example: "pubsub_cli connect gcp_project test_topic --host=localhost:8085 --project=emulator",
		Aliases: []string{"c"},
		Args:    cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			remoteProjectID := args[0]
			topicIDs := args[1:]
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

			return connect(ctx, out, remotePubsubClient, localPubsubClient, topicIDs)
		},
	}
}

type subscriberForConnect struct {
	remoteTopic *pubsub.Topic
	localTopic  *pubsub.Topic
	sub         *pubsub.Subscription
}

// connect connects remote Pub/Sub topics to local Pub/Sub topics
func connect(ctx context.Context, out io.Writer, remotePubsubClient, localPubsubClient *pkg.PubSubClient, topicIDs []string) error {
	var remoteTopics []*pubsub.Topic
	var err error
	if topicIDs[0] == "all" {
		remoteTopics, err = remotePubsubClient.FindAllTopics(ctx)
		if err != nil {
			return errors.Wrap(err, "find all topics %#v")
		}
	} else {
		remoteTopics, err = remotePubsubClient.FindOrCreateTopics(ctx, topicIDs)
		if err != nil {
			return errors.Wrapf(err, "find or create remote topics %#v", topicIDs)
		}
	}

	remoteTopicIDs := make([]string, len(remoteTopics))
	for i, t := range remoteTopics {
		remoteTopicIDs[i] = t.ID()
	}
	localTopics, err := localPubsubClient.FindOrCreateTopics(ctx, remoteTopicIDs)
	if err != nil {
		return errors.Wrapf(err, "find or create local topics %#v", remoteTopicIDs)
	}

	sort.Slice(remoteTopics, func(i, j int) bool {
		return remoteTopics[i].ID() > remoteTopics[j].ID()
	})
	sort.Slice(localTopics, func(i, j int) bool {
		return localTopics[i].ID() > localTopics[j].ID()
	})

	eg := &errgroup.Group{}
	subscribers := make(chan *subscriberForConnect, len(remoteTopics))
	for i, topic := range remoteTopics {
		remoteTopic := topic
		localTopic := localTopics[i]
		eg.Go(func() error {
			fmt.Println(fmt.Sprintf("[start]creating unique subscription to %s...", remoteTopic.String()))
			sub, err := remotePubsubClient.CreateUniqueSubscription(ctx, remoteTopic)
			if err != nil {
				return errors.Wrapf(err, "create unique subscription to %s", remoteTopic.String())
			}
			subscribers <- &subscriberForConnect{remoteTopic: remoteTopic, localTopic: localTopic, sub: sub}
			_, _ = colorstring.Fprintf(out, "[green][success] created subscription to %s\n", remoteTopic.String())
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	close(subscribers)
	_, _ = colorstring.Fprintln(out, "\n[green][success] topics are now connected!\n")

	fmt.Println("[start] waiting for publish...")
	for s := range subscribers {
		s := s
		eg.Go(func() error {
			err := s.sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
				_, _ = colorstring.Println(fmt.Sprintf("[green][success] Got message: %q", string(msg.Data)))
				messageID, err := s.localTopic.Publish(ctx, msg).Get(ctx)
				if err != nil {
					_, _ = colorstring.Fprintln(out, fmt.Sprintf("[red][error] publish message with data = %s", msg.Data))
					return
				}
				_, _ = colorstring.Fprintln(out, fmt.Sprintf("[green][success] published message to %s successfully, message ID = %s\n", s.localTopic.String(), messageID))
				msg.Ack()
			})
			if err != nil {
				return errors.Wrapf(err, "createSubscription %s failed", s.remoteTopic.String())
			}
			return nil
		})
	}
	return eg.Wait()
}
