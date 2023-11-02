package resolver

import upstreams "github.com/xpizy2020/upsm"

type Resolver interface {
	Init(options ...interface{}) error
	Subscribe(keys interface{}, ups upstreams.Upstreams) error
	UnSubscribe(keys interface{}) error
}
