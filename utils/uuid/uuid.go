package uuid

import (
	"github.com/rs/xid"
)

//UUID generates uuid
func UUID() string {
	id := xid.New()
	return id.String()
}
