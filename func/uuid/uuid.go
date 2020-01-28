package uuid

import (
	"github.com/rs/xid"
)

const IDLen = 20

// NewUUID Generates a new unique id
func NewUUID() string {
	id := xid.New()
	return id.String()
}
