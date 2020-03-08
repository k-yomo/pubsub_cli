package cmd

import (
	"bytes"
	"context"
	"github.com/k-yomo/pubsub_cli/util"
	"github.com/spf13/cobra"
	"testing"
	"time"
)

func Test_registerPush(t *testing.T) {
	pubsubClient, err := util.NewTestPubSubClient(t)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		in0          *cobra.Command
		pubsubClient *util.PubSubClient
		args         []string
	}
	tests := []struct {
		name               string
		mockSubscriptionID string
		args               args
		check              func()
		wantErr            bool
	}{
		{
			name:               "push subscription is expected to be registered successfully",
			mockSubscriptionID: "test",
			args:               args{pubsubClient: pubsubClient, args: []string{"test_topic", "http://localhost:9000"}},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clear := util.SetMockUUID(tt.mockSubscriptionID)
			defer clear()

			out := &bytes.Buffer{}
			err := registerPush(tt.args.in0, out, tt.args.pubsubClient, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerPush() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.check()
		})
	}
}
