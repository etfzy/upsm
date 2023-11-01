package loadbalancer

import "errors"

const (
	WRR = "wrr"
)

var (
	ErrNoServerFound = errors.New("server find null error")
)

type LoadBalancer interface {
	Picker(servers []AddrInfos) string
	BuildAddr(info interface{}) AddrInfos
}

type AddrInfos interface {
	GetAddr() string
}

type Builder func(info interface{}) AddrInfos
