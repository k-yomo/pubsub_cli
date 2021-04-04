package cmd

import (
	"bytes"
	"context"
	"github.com/k-yomo/pubsub_cli/pkg"
	"github.com/spf13/cobra"
	"testing"
)

func Test_createSubscription(t *testing.T) {
	pubsubClient, err := pkg.NewTestPubSubClient(t)
	if err != nil {
		t.Fatal(err)
	}
	rootCmd := newTestRootCmd(t)

	subscriptionID := pkg.UUID()

	type args struct {
		rootCmd *cobra.Command
		args    []string
	}
	tests := []struct {
		name               string
		mockSubscriptionID string
		args               args
		check              func()
		wantErr            bool
	}{
		{
			name:               "subscription is created successfully",
			mockSubscriptionID: "test",
			args:               args{rootCmd: rootCmd, args: []string{"create_subscription", "test_topic", subscriptionID}},
			check: func() {
				sub := pubsubClient.Subscription(subscriptionID)
				subConfig, err := sub.Config(context.Background())
				if err != nil {
					t.Fatal(err)
				}
				topic := "test_topic"
				// check if topic is collect
				if subConfig.Topic.ID() != topic {
					t.Errorf("createSubscription() got topic = %v, want %v", subConfig.Topic.String(), topic)
				}
				sub.Delete(context.Background())
			},
		},
		{
			name:    "push subscription with invalid topic name causes error",
			args:    args{rootCmd: rootCmd, args: []string{"create_subscription", "a", "test_topic_sub"}},
			check:   func() {},
			wantErr: true,
		},
		{
			name: "parent cmd without projectFlag causes error",
			args: args{rootCmd: func() *cobra.Command {
				cmd := &cobra.Command{}
				cmd.PersistentFlags().String(hostFlagName, "host", "")
				cmd.PersistentFlags().String(credFileFlagName, "cred.json", "")
				return cmd
			}(), args: []string{"create_subscription", "test_topic", "test_topic_sub"}},
			check:   func() {},
			wantErr: true,
		},
		{
			name: "parent cmd without hostFlag causes error",
			args: args{rootCmd: func() *cobra.Command {
				cmd := &cobra.Command{}
				cmd.PersistentFlags().String(projectFlagName, "project", "")
				cmd.PersistentFlags().String(credFileFlagName, "cred.json", "")
				return cmd
			}(), args: []string{"create_subscription", "test_topic", "test_topic_sub"}},
			check:   func() {},
			wantErr: true,
		},
		{
			name: "parent cmd without credFileFlag causes error",
			args: args{rootCmd: func() *cobra.Command {
				cmd := &cobra.Command{}
				cmd.PersistentFlags().String(projectFlagName, "project", "")
				cmd.PersistentFlags().String(hostFlagName, "host", "")
				return cmd
			}(), args: []string{"create_subscription", "test_topic", "test_topic_sub"}},
			check:   func() {},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clear := pkg.SetMockUUID(t, tt.mockSubscriptionID)
			defer clear()

			out := &bytes.Buffer{}
			cmd := newCreateSubscriptionCmd(out)
			tt.args.rootCmd.SetArgs(tt.args.args)
			tt.args.rootCmd.AddCommand(cmd)

			err := tt.args.rootCmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("createSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.check()
		})
	}
}
