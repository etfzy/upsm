package example

import (
	"fmt"
	"testing"

	upstreams "github.com/xpizy2020/upsm"
	"github.com/xpizy2020/upsm/loadbalancer"
	nr "github.com/xpizy2020/upsm/resolver/nacosresolver"
)

func TestNacosUps(t *testing.T) {
	t.Run("wrr load balance test", func(t *testing.T) {
		ups := upstreams.NewUpstreams()

		//构建loadbalancer
		wrrlb := loadbalancer.NewWrrLoadBalancer(loadbalancer.DefaultWrrBuilder)
		ups.SetLoadBalancer(wrrlb)

		//创建一个resolver
		inst := nr.CreateWrrNacosResolver()
		inst.Init(nr.WithUsername("nacos"), nr.WithPwd("nacos"), nr.WithServers([]string{""}), nr.WithPort(), nr.WithNameSpace(""))

		inst.Subscribe(&nr.Watcher{
			Server:  "",
			Group:   "",
			Cluster: "",
		}, ups.Publish)

		addr, _, _ := ups.NextWithoutConnect()
		fmt.Println("xxxx:", addr)

		/*
			m := make(map[string]int)
			lock := sync.Mutex{}
			sg := sync.WaitGroup{}

			for i := 0; i < 4; i++ {
				sg.Add(1)
				go func() {
					defer sg.Done()
					addr, _, _ := ups.NextWithConnect()
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
			}*/
	})

}
