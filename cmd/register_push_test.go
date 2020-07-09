package cmd

import (
	"bytes"
	"context"
	"github.com/k-yomo/pubsub_cli/pkg"
	"github.com/spf13/cobra"
	"testing"
	"time"
)

func Test_registerPush(t *testing.T) {
	pubsubClient, err := pkg.NewTestPubSubClient(t)
	if err != nil {
		t.Fatal(err)
	}
	rootCmd := newTestRootCmd(t)

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
			name:               "push subscription is registered successfully",
			mockSubscriptionID: "test",
			args:               args{rootCmd: rootCmd, args: []string{"register_push", "test_topic", "http://localhost:9000"}},
			check: func() {
				sub := pubsubClient.Subscription("test")
				subConfig, err := sub.Config(context.Background())
				if err != nil {
					t.Fatal(err)
				}
				topic := "test_topic"
				// check if topic is collect
				if subConfig.Topic.ID() != topic {
					t.Errorf("registerPush() got topic = %v, want %v", subConfig.Topic.String(), topic)
				}
				// check if endpoint is collect
				if subConfig.PushConfig.Endpoint != "http://localhost:9000" {
					t.Errorf("registerPush() got endpoint = %v, want %v", subConfig.PushConfig.Endpoint, "http://localhost:9000")
				}
				// check if expirationPolicy is set to 24 hours
				if subConfig.ExpirationPolicy != 24*time.Hour {
					t.Errorf("registerPush() got expirationPolicy = %v, want %v", subConfig.ExpirationPolicy, 24*time.Hour)
				}
				sub.Delete(context.Background())
			},
		},
		{
			name:    "push subscription with invalid topic name causes error",
			args:    args{rootCmd: rootCmd, args: []string{"register_push", "1", "http://localhost:9000"}},
			check:   func() {},
			wantErr: true,
		},
		{
			name:    "push subscription with invalid endpoint causes error",
			args:    args{rootCmd: rootCmd, args: []string{"register_push", "test_topic", "invalid"}},
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
			}(), args: []string{"publish", "test_topic", "hello"}},
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
			}(), args: []string{"publish", "test_topic", "hello"}},
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
			}(), args: []string{"publish", "test_topic", "hello"}},
			check:   func() {},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clear := pkg.SetMockUUID(t, tt.mockSubscriptionID)
			defer clear()

			out := &bytes.Buffer{}
			cmd := newRegisterPushCmd(out)
			tt.args.rootCmd.SetArgs(tt.args.args)
			tt.args.rootCmd.AddCommand(cmd)

			err := tt.args.rootCmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("registerPush() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.check()
		})
	}
}
