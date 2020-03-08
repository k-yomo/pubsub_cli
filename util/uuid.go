package util

import "github.com/rs/xid"

type UUIDGenerator interface {
	String() string
}

var idgen UUIDGenerator = xid.New()

func UUID() string {
	return idgen.String()
}

func SetMockUUID(uuid string) (clear func()) {
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
