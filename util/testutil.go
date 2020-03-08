package util

import (
	"context"
	"testing"
)

func NewTestPubSubClient(_ *testing.T) (*PubSubClient, error) {
	return NewPubSubClient(context.Background(), "test", "localhost:8085", "")
}
