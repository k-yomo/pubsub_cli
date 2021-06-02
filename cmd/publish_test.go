package cmd

import (
	"bytes"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/k-yomo/pubsub_cli/pkg"
	"github.com/spf13/cobra"
	"sync"
	"testing"
)

func Test_publish(t *testing.T) {
	t.Parallel()

	pubsubClient, err := pkg.NewTestPubSubClient(t)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		rootCmd *cobra.Command
		args    []string
	}
	tests := []struct {
		name     string
		args     args
		before   func() *pubsub.Subscription
		wantData string
		wantAttributes map[string]string
		wantErr  bool
	}{
		{
			name: "message is expected to be published successfully",
			args: args{rootCmd: newTestRootCmd(t), args: []string{"publish", "publish_topic1", "hello"}},
			before: func() *pubsub.Subscription {
				topic, err := pubsubClient.FindOrCreateTopic(context.Background(), "publish_topic1")
				if err != nil {
					t.Fatal(err)
				}
				sub, err := pubsubClient.CreateUniqueSubscription(context.Background(), topic)
				if err != nil {
					t.Fatal(err)
				}
				return sub
			},
			wantData: "hello",
		},
		{
			name: "message with attributes is expected to be published successfully",
			args: args{rootCmd: newTestRootCmd(t), args: []string{"publish", "publish_topic2", "hello", "--attribute=k1=v1", "-a k2=v2"}},
			before: func() *pubsub.Subscription {
				topic, err := pubsubClient.FindOrCreateTopic(context.Background(), "publish_topic2")
				if err != nil {
					t.Fatal(err)
				}
				sub, err := pubsubClient.CreateUniqueSubscription(context.Background(), topic)
				if err != nil {
					t.Fatal(err)
				}
				return sub
			},
			wantData: "hello",
			wantAttributes: map[string]string{
				"k1": "v1",
				"k2": "v2",
			},
		},
		{
			name:    "publish to topic with invalid attribute causes error",
			args:    args{rootCmd: newTestRootCmd(t), args: []string{"publish", "publish_topic3", "hello", "-a k1:v1"}},
			before:  func() *pubsub.Subscription { return &pubsub.Subscription{} },
			wantErr: true,
		},
		{
			name:    "publish to topic with invalid name causes error",
			args:    args{rootCmd: newTestRootCmd(t), args: []string{"publish", "1", "hello"}},
			before:  func() *pubsub.Subscription { return &pubsub.Subscription{} },
			wantErr: true,
		},
		{
			name:    "publish empty message causes error",
			args:    args{rootCmd: newTestRootCmd(t), args: []string{"publish", "publish_topic", ""}},
			before:  func() *pubsub.Subscription { return &pubsub.Subscription{} },
			wantErr: true,
		},
		{
			name: "parent cmd without projectFlag causes error",
			args: args{rootCmd: func() *cobra.Command {
				cmd := &cobra.Command{}
				cmd.PersistentFlags().String(hostFlagName, "host", "")
				cmd.PersistentFlags().String(credFileFlagName, "cred.json", "")
				return cmd
			}(), args: []string{"publish", "publish_topic", "hello"}},
			before:  func() *pubsub.Subscription { return &pubsub.Subscription{} },
			wantErr: true,
		},
		{
			name: "parent cmd without hostFlag causes error",
			args: args{rootCmd: func() *cobra.Command {
				cmd := &cobra.Command{}
				cmd.PersistentFlags().String(projectFlagName, "project", "")
				cmd.PersistentFlags().String(credFileFlagName, "cred.json", "")
				return cmd
			}(), args: []string{"publish", "publish_topic", "hello"}},
			before:  func() *pubsub.Subscription { return &pubsub.Subscription{} },
			wantErr: true,
		},
		{
			name: "parent cmd without credFileFlag causes error",
			args: args{rootCmd: func() *cobra.Command {
				cmd := &cobra.Command{}
				cmd.PersistentFlags().String(projectFlagName, "project", "")
				cmd.PersistentFlags().String(hostFlagName, "host", "")
				return cmd
			}(), args: []string{"publish", "publish_topic", "hello"}},
			before:  func() *pubsub.Subscription { return &pubsub.Subscription{} },
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sub := tt.before()
			out := &bytes.Buffer{}
			cmd := newPublishCmd(out)
			tt.args.rootCmd.SetArgs(tt.args.args)
			tt.args.rootCmd.AddCommand(cmd)

			err := tt.args.rootCmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("publish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				ctx, cancel := context.WithCancel(context.Background())
				err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
					defer wg.Done()
					defer cancel()
					msg.Ack()
					if string(msg.Data) != tt.wantData {
						t.Errorf("publish() gotData = %v, want %v", string(msg.Data), tt.wantData)
					}
					if diff := cmp.Diff(msg.Attributes, tt.wantAttributes); diff != "" {
						t.Errorf("publish() attributes diff %v", diff)
					}
				})
				if err != nil && ctx.Err() == nil {
					t.Fatal(err)
				}
			}()
			wg.Wait()
		})
	}
}
