package example

import (
	"fmt"
	"sync"
	"testing"

	upstreams "github.com/xpizy2020/upsm"
	"github.com/xpizy2020/upsm/loadbalancer"
)

func TestWrrUps(t *testing.T) {
	t.Run("wrr load balance test", func(t *testing.T) {
		ups := upstreams.NewUpstreams()

		//构建loadbalancer
		wrrlb := loadbalancer.NewWrrLoadBalancer(loadbalancer.DefaultWrrBuilder)
		ups.SetLoadBalancer(wrrlb)

		//创建一个server
		ups.AddServer(loadbalancer.CreateWrrServer("1.1.1.1", 10))
		ups.AddServer(loadbalancer.CreateWrrServer("1.1.1.2", 10))
		ups.AddServer(loadbalancer.CreateWrrServer("1.1.1.3", 20))

		m := make(map[string]int)
		lock := sync.Mutex{}
		sg := sync.WaitGroup{}

		for i := 0; i < 4; i++ {
			sg.Add(1)
			go func() {
				defer sg.Done()
				addr, _, _ := ups.NextWithoutConnect()
				lock.Lock()
				_, ok := m[addr]
				if ok {
					m[addr] = m[addr] + 1
				} else {
					m[addr] = 1
				}
				lock.Unlock()
			}()

		}
		sg.Wait()

		v, ok := m["1.1.1.1"]
		if !ok || v != 1 {
			fmt.Println(m)
			t.Errorf("1.1.1.1 %d", v)
		}

		v, ok = m["1.1.1.2"]
		if !ok || v != 1 {
			fmt.Println(m)
			t.Errorf("1.1.1.2 %d", v)
		}

		v, ok = m["1.1.1.3"]
		if !ok || v != 2 {
			fmt.Println(m)
			t.Errorf("1.1.1.3 %d", v)
		}
	})

	t.Run("resfresh wrr weight load balance test", func(t *testing.T) {
		ups := upstreams.NewUpstreams()
		//构建loadbalancer
		wrrlb := loadbalancer.NewWrrLoadBalancer(loadbalancer.DefaultWrrBuilder)
		ups.SetLoadBalancer(wrrlb)

		//创建一个server
		ups.AddServer(loadbalancer.CreateWrrServer("1.1.1.1", 10))
		ups.AddServer(loadbalancer.CreateWrrServer("1.1.1.2", 10))
		ups.AddServer(loadbalancer.CreateWrrServer("1.1.1.3", 20))
		ups.AddServer(loadbalancer.CreateWrrServer("1.1.1.3", 10))
		m := make(map[string]int)
		lock := sync.Mutex{}
		sg := sync.WaitGroup{}

		for i := 0; i < 4; i++ {
			sg.Add(1)
			go func() {
				defer sg.Done()
				addr, _, _ := ups.NextWithoutConnect()
				lock.Lock()
				_, ok := m[addr]
				if ok {
					m[addr] = m[addr] + 1
				} else {
					m[addr] = 1
				}
				lock.Unlock()
			}()

		}
		sg.Wait()

		v, ok := m["1.1.1.1"]
		if !ok || v != 2 {
			t.Errorf("1.1.1.1 %d", v)
		}

		v, ok = m["1.1.1.2"]
		if !ok || v != 1 {
			fmt.Println(m)
			t.Errorf("1.1.1.2 %d", v)
		}

		v, ok = m["1.1.1.3"]
		if !ok || v != 1 {
			fmt.Println(m)
			t.Errorf("1.1.1.3 %d", v)
		}
	})

	t.Run("remove wrr load balance test", func(t *testing.T) {
		ups := upstreams.NewUpstreams()

		//构建loadbalancer
		wrrlb := loadbalancer.NewWrrLoadBalancer(loadbalancer.DefaultWrrBuilder)
		ups.SetLoadBalancer(wrrlb)

		//创建一个server
		ups.AddServer(loadbalancer.CreateWrrServer("1.1.1.1", 10))
		ups.AddServer(loadbalancer.CreateWrrServer("1.1.1.2", 10))
		ups.AddServer(loadbalancer.CreateWrrServer("1.1.1.3", 20))
		ups.AddServer(loadbalancer.CreateWrrServer("1.1.1.3", 10))
		ups.RemoveServer("1.1.1.3")
		m := make(map[string]int)
		lock := sync.Mutex{}
		sg := sync.WaitGroup{}

		for i := 0; i < 4; i++ {
			sg.Add(1)
			go func() {
				defer sg.Done()
				addr, _, _ := ups.NextWithoutConnect()
				lock.Lock()
				_, ok := m[addr]
				if ok {
					m[addr] = m[addr] + 1
				} else {
					m[addr] = 1
				}
				lock.Unlock()
			}()

		}
		sg.Wait()

		v, ok := m["1.1.1.1"]
		if !ok || v != 2 {
			t.Errorf("1.1.1.1 %d", v)
		}

		v, ok = m["1.1.1.2"]
		if !ok || v != 2 {
			fmt.Println(m)
			t.Errorf("1.1.1.2 %d", v)
		}

		v, ok = m["1.1.1.3"]
		if ok {
			fmt.Println(m)
			t.Errorf("1.1.1.3 %d", v)
		}
	})
}
