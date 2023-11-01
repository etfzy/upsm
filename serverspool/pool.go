package serverspool

import (
	"sync"

	"github.com/xpizy2020/upsm/loadbalancer"
)

type ServersPool struct {
	lock       sync.Mutex
	addrsMap   map[string]loadbalancer.AddrInfos //addrs map in pool
	addrsSlice []loadbalancer.AddrInfos          //addrs slices  in pool
}

func (ap *ServersPool) Publish(info []loadbalancer.AddrInfos) {
	ap.lock.Lock()
	defer ap.lock.Unlock()

	if info == nil {
		return
	}

	if len(info) == 0 {
		return
	}

	if ap.addrsMap == nil {
		ap.addrsMap = make(map[string]loadbalancer.AddrInfos)
	}

	if ap.addrsSlice == nil {
		ap.addrsSlice = make([]loadbalancer.AddrInfos, 0, 1024)
	} else {
		ap.addrsSlice = ap.addrsSlice[:0]
	}
	for k, v := range info {
		addr := v.GetAddr()
		ap.addrsMap[addr] = info[k]
		ap.addrsSlice = append(ap.addrsSlice, info[k])
	}

}

func (ap *ServersPool) Refresh(info loadbalancer.AddrInfos) {
	ap.lock.Lock()
	defer ap.lock.Unlock()
	addr := info.GetAddr()
	if _, ok := ap.addrsMap[addr]; !ok {
		ap.addrsSlice = append(ap.addrsSlice, info)
	} else {
		for i, v := range ap.addrsSlice {
			temp := v.GetAddr()
			if temp == addr {
				ap.addrsSlice = append(ap.addrsSlice[:i], ap.addrsSlice[i+1:]...)
				break
			}
		}
		ap.addrsSlice = append(ap.addrsSlice, info)
	}

	ap.addrsMap[addr] = info
}

func (ap *ServersPool) Remove(addr string) {
	ap.lock.Lock()
	defer ap.lock.Unlock()
	if _, ok := ap.addrsMap[addr]; ok {
		delete(ap.addrsMap, addr)
		for i, v := range ap.addrsSlice {
			temp := v.GetAddr()
			if temp == addr {
				ap.addrsSlice = append(ap.addrsSlice[:i], ap.addrsSlice[i+1:]...)
				break
			}
		}

	}
}

func (ap *ServersPool) Next(lb loadbalancer.LoadBalancer) string {
	ap.lock.Lock()
	defer ap.lock.Unlock()
	return lb.Picker(ap.addrsSlice)
}

/**
 * initialize the address pool
 */

func (ap *ServersPool) Init() {
	ap.lock.Lock()
	defer ap.lock.Unlock()

	if ap.addrsMap == nil {
		ap.addrsMap = make(map[string]loadbalancer.AddrInfos)
	}

	if ap.addrsSlice == nil {
		ap.addrsSlice = make([]loadbalancer.AddrInfos, 0, 1024)
	} else {
		ap.addrsSlice = ap.addrsSlice[:0]
	}
}

func (ap *ServersPool) Reset() {
	ap.lock.Lock()
	defer ap.lock.Unlock()

	if ap.addrsMap == nil {
		ap.addrsMap = make(map[string]loadbalancer.AddrInfos)
	}

	if ap.addrsSlice == nil {
		ap.addrsSlice = make([]loadbalancer.AddrInfos, 0, 1024)
	} else {
		ap.addrsSlice = ap.addrsSlice[:0]
	}
}
