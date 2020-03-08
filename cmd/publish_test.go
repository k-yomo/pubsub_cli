package cmd

import (
	"bytes"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/k-yomo/pubsub_cli/util"
	"github.com/spf13/cobra"
	"sync"
	"testing"
)

func Test_publish(t *testing.T) {
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
		name     string
		args     args
		before   func() *pubsub.Subscription
		wantData string
		wantErr  bool
	}{
		{
			name: "message is expected to be published successfully",
			args: args{pubsubClient: pubsubClient, args: []string{"test_topic", "hello"}},
			before: func() *pubsub.Subscription {
				topic, err := pubsubClient.FindOrCreateTopic(context.Background(), "test_topic")
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := tt.before()
			out := &bytes.Buffer{}
			err := publish(tt.args.in0, out, tt.args.pubsubClient, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("publish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				err := sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
					defer wg.Done()
					msg.Ack()
					if string(msg.Data) != tt.wantData {
						t.Errorf("publish() gotData = %v, want %v", string(msg.Data), tt.wantData)
					}
				})
				if err != nil {
					t.Fatal(err)
				}
			}()
			wg.Wait()
		})
	}
}
