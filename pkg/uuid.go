package pkg

import (
	"github.com/rs/xid"
)

// UUIDGenerator represents uuid generator
type UUIDGenerator interface {
	String() string
}

var idgen UUIDGenerator = xid.New()

// UUID generates uuid
func UUID() string {
	return idgen.String()
}
