package cmd

import (
	"bytes"
	"context"
	"github.com/k-yomo/pubsub_cli/pkg"
	"github.com/spf13/cobra"
	"testing"
)

func Test_createTopic(t *testing.T) {
	t.Parallel()

	pubsubClient, err := pkg.NewTestPubSubClient(t)
	if err != nil {
		t.Fatal(err)
	}

	topicID := pkg.UUID()

	type args struct {
		rootCmd *cobra.Command
		args    []string
	}
	tests := []struct {
		name               string
		args               args
		check              func()
		wantErr            bool
	}{
		{
			name:               "topic is created successfully",
			args:               args{
				rootCmd: newTestRootCmd(t),
				args: []string{"create_topic", topicID},
			},
			check: func() {
				if _, err := pubsubClient.FindTopic(context.Background(), topicID); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name:               "skip existing topic",
			args:               args{
				rootCmd: newTestRootCmd(t),
				args: []string{"create_topic", topicID, topicID},
			},
			check: func() {},
		},
		{
			name:    "invalid topic name causes error",
			args:    args{rootCmd: newTestRootCmd(t), args: []string{"create_topic", "a"}},
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
			}(), args: []string{"create_topic", "test_topic"}},
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
			}(), args: []string{"create_topic", "test_topic"}},
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
			}(), args: []string{"create_topic", "test_topic"}},
			check:   func() {},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cmd := newCreateTopicCmd(out)
			tt.args.rootCmd.SetArgs(tt.args.args)
			tt.args.rootCmd.AddCommand(cmd)

			err := tt.args.rootCmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("createTopic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.check()
		})
	}
}
