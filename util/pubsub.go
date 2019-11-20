package util

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"time"
)

type PubSubClient struct {
	*pubsub.Client
}

func NewPubSubClient(projectID string, pubsubEmulatorHost string) (*PubSubClient, error) {
	if projectID == "" {
		return nil, errors.New("GCP Project ID must be set with either env varibale 'GCP_PROJECT_ID' or --project flag")
	}
	if pubsubEmulatorHost == "" {
		return nil, errors.New("Emulator host must be set with either env varibale 'PUBSUB_EMULATOR_HOST' or --host flag")
	}

	conn, err := grpc.Dial(pubsubEmulatorHost, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "grpc.Dial")
	}
	client, err := pubsub.NewClient(context.Background(), projectID, option.WithGRPCConn(conn))
	if err != nil {
		return nil, errors.Wrap(err, "intialize new pubsub client failed")
	}
	return &PubSubClient{client}, nil
}

func (pc *PubSubClient) FindOrCreateTopic(ctx context.Context, topicID string) (*pubsub.Topic, error) {
	topic := pc.Topic(topicID)

	exists, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	} else if exists {
		return topic, nil
	}

	topic, err = pc.CreateTopic(ctx, topicID)
	if err != nil {
		return nil, err
	}
	return topic, nil
}

func (pc *PubSubClient) CreateUniqueSubscription(topic *pubsub.Topic) (*pubsub.Subscription, error) {
	subscriptionConfig := pubsub.SubscriptionConfig{
		Topic:            topic,
		ExpirationPolicy: time.Hour * 24,
	}
	sub, err := pc.CreateSubscription(context.Background(), xid.New().String(), subscriptionConfig)
	if err != nil {
		return nil, err
	}
	return sub, err

}
