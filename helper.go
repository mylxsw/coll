package coll

import (
	"github.com/mylxsw/coll"
)

func Map(origin interface{}, dest interface{}, mapper interface{}) error {
	return coll.MustNew(origin).Map(mapper).All(dest)
}
