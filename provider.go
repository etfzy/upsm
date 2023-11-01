package upstreams

import (
	"errors"

	"github.com/xpizy2020/upsm/connspool"
	"github.com/xpizy2020/upsm/loadbalancer"

	"github.com/xpizy2020/upsm/serverspool"
)

type Publish func(info []loadbalancer.AddrInfos)

type Upstreams interface {
	//构造方法
	SetLoadBalancer(lb loadbalancer.LoadBalancer) Upstreams
	SetConnsPool(pool connspool.ConnsPool) Upstreams

	//功能方法
	NextWithoutConnect() (string, interface{}, error)
	NextWithConnect() (string, interface{}, error)
	AddServer(info loadbalancer.AddrInfos) error
	UpdateServer(info loadbalancer.AddrInfos)
	RemoveServer(addr string)

	Publish(info []loadbalancer.AddrInfos)
}

type upstreamsInfos struct {
	servers serverspool.Pool
	lb      loadbalancer.LoadBalancer
	conns   connspool.ConnsPool
	err     error
}

func (p *upstreamsInfos) Publish(info []loadbalancer.AddrInfos) {
	p.servers.Publish(info)
	if p.conns != nil {
		p.conns.Publish(info)
	}
}

// get next addr by using loadbalancer
func (p *upstreamsInfos) NextWithoutConnect() (string, interface{}, error) {
	addr := p.servers.Next(p.lb)

	if p.conns == nil {
		return addr, nil, errors.New("connections pool is not init!")
	}
	return addr, p.conns.GetConnsWithoutConnect(addr), nil
}

// get next addr by using loadbalancer
func (p *upstreamsInfos) NextWithConnect() (string, interface{}, error) {
	addr := p.servers.Next(p.lb)

	if p.conns == nil {
		return addr, nil, errors.New("connections pool is not init!")
	}
	conn, err := p.conns.GetConnsWithConnect(addr)
	return addr, conn, err
}

// set load balance mode,
//
//eg:wrr、rr
func (p *upstreamsInfos) SetLoadBalancer(lb loadbalancer.LoadBalancer) Upstreams {
	p.lb = lb
	return p
}

func (p *upstreamsInfos) SetConnsPool(pool connspool.ConnsPool) Upstreams {
	p.conns = pool
	return p
}

// add a server
func (p *upstreamsInfos) AddServer(info loadbalancer.AddrInfos) error {
	p.servers.Refresh(info)
	return nil
}

// update a server
func (p *upstreamsInfos) UpdateServer(info loadbalancer.AddrInfos) {
	p.servers.Refresh(info)
	return
}

// remove a server
func (p *upstreamsInfos) RemoveServer(addr string) {
	if p.conns != nil {
		p.conns.RemoveConns(addr)
	}

	p.servers.Remove(addr)
	return
}

// create a upstreams
func NewUpstreams() Upstreams {
	temp := &upstreamsInfos{
		servers: &serverspool.ServersPool{},
	}

	temp.servers.Reset()

	return temp
}
