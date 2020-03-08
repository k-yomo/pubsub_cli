package util

import (
	"context"
	"testing"
)

func NewTestPubSubClient(t *testing.T) (*PubSubClient, error) {
	t.Helper()
	return NewPubSubClient(context.Background(), "test", "localhost:8085", "")
}
