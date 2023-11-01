package connspool

import (
	"sync"

	"github.com/xpizy2020/upsm/loadbalancer"
)

type ConnsPool interface {
	GetConnsWithConnect(addr string) (interface{}, error)
	GetConnsWithoutConnect(addr string) interface{}
	RemoveConns(addr string)
	Publish(addrs []loadbalancer.AddrInfos)
}

// interface 的本身可以是一个pool
type DefaultConnsPool struct {
	lock      sync.RWMutex
	mConns    map[string]interface{}
	buildConn ConnsBuilder
}

func CreateDefaultPool(builder ConnsBuilder) ConnsPool {
	return &DefaultConnsPool{
		lock:      sync.RWMutex{},
		mConns:    make(map[string]interface{}),
		buildConn: builder,
	}
}

func (p *DefaultConnsPool) addConnection(addr string, conn interface{}) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.mConns[addr] = conn
}

func (p *DefaultConnsPool) GetConnsWithConnect(addr string) (interface{}, error) {
	conn := p.getConnection(addr)
	if conn == nil {
		inst, err := p.buildConn(addr)
		if err != nil {
			return nil, err
		}

		p.addConnection(addr, inst)
	}

	return p.getConnection(addr), nil
}

func (p *DefaultConnsPool) getConnection(addr string) interface{} {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.mConns[addr]
}

func (p *DefaultConnsPool) GetConnsWithoutConnect(addr string) interface{} {
	return p.getConnection(addr)
}

func (p *DefaultConnsPool) RemoveConns(addr string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	delete(p.mConns, addr)
}

func findConns(addrs []loadbalancer.AddrInfos, addr string) bool {

	for _, v := range addrs {
		if addr == v.GetAddr() {
			return true
		}
	}

	return false
}
func (p *DefaultConnsPool) Publish(addrs []loadbalancer.AddrInfos) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for k, _ := range p.mConns {
		if !findConns(addrs, k) {
			delete(p.mConns, k)
		}
	}
}
