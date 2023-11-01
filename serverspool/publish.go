package serverspool

import "github.com/xpizy2020/upsm/loadbalancer"

type Pool interface {
	Publish(info []loadbalancer.AddrInfos)
	Refresh(info loadbalancer.AddrInfos)
	Remove(addr string)
	Next(lb loadbalancer.LoadBalancer) string
	Reset()
	Init()
}
