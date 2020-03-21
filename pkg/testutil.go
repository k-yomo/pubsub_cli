package pkg

import (
	"context"
	"testing"

	"github.com/rs/xid"
)

// NewTestPubSubClient initializes new pubsub client for local pubsub emulator
func NewTestPubSubClient(t *testing.T) (*PubSubClient, error) {
	t.Helper()
	return NewPubSubClient(context.Background(), "test", "localhost:8085", "")
}

// SetMockUUID mocks uuid with given string
func SetMockUUID(t *testing.T, uuid string) (clear func()) {
	t.Helper()
	m := &mockIDGen{uuid: uuid}
	idgen = m
	return func() {
		idgen = xid.New()
	}
}

type mockIDGen struct {
	uuid string
}

func (m *mockIDGen) String() string {
	return m.uuid
}
