package pkg

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"google.golang.org/api/option"
)

// PubSubClient represents extended pubsub client
type PubSubClient struct {
	*pubsub.Client
}

// NewPubSubClient initializes new pubsub client
func NewPubSubClient(ctx context.Context, projectID, pubsubEmulatorHost, gcpCredFilePath string) (*PubSubClient, error) {
	if projectID == "" {
		return nil, errors.New("GCP Project ID must be set from either env variable 'GCP_PROJECT_ID' or --project flag")
	}

	var opts []option.ClientOption
	if pubsubEmulatorHost != "" {
		conn, err := grpc.DialContext(ctx, pubsubEmulatorHost, grpc.WithInsecure())
		if err != nil {
			return nil, errors.Wrap(err, "grpc.Dial")
		}
		opts = append(opts, option.WithGRPCConn(conn))
	} else {
		opts = append(opts, option.WithCredentialsFile(gcpCredFilePath))
	}

	client, err := pubsub.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "create new pubsub client")
	}
	return &PubSubClient{client}, nil
}

// FindOrCreateTopic finds the topic or create if not exists.
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

// CreateUniqueSubscription creates an unique subscription to given topic
func (pc *PubSubClient) CreateUniqueSubscription(ctx context.Context, topic *pubsub.Topic) (*pubsub.Subscription, error) {
	subscriptionConfig := pubsub.SubscriptionConfig{
		Topic:            topic,
		ExpirationPolicy: time.Hour * 24,
	}
	sub, err := pc.CreateSubscription(ctx, fmt.Sprintf("pubsub_cli_%s", xid.New().String()), subscriptionConfig)
	if err != nil {
		return nil, err
	}
	return sub, err
}
