package example

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	upstreams "github.com/xpizy2020/upsm"
	"github.com/xpizy2020/upsm/connspool"
	"github.com/xpizy2020/upsm/loadbalancer"
)

func connbuilder(addr string) (interface{}, error) {
	// 创建http.Client对象
	httpTransport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     false,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	httpClient := &http.Client{
		Transport: httpTransport,
		Timeout:   time.Second * 5, //超时时间
	}
	return httpClient, nil
}

func TestConnsUps(t *testing.T) {
	t.Run("http connections pool test", func(t *testing.T) {
		ups := upstreams.NewUpstreams()
		//构建loadbalancer
		wrrlb := loadbalancer.NewWrrLoadBalancer(loadbalancer.DefaultWrrBuilder)
		ups.SetLoadBalancer(wrrlb)

		//构建conn pool
		connpool := connspool.CreateDefaultPool(connbuilder)
		//创建connection builder,传入builer 的
		ups.SetConnsPool(connpool)

		//创建一个server
		ups.AddServer(loadbalancer.CreateWrrServer("http://baidu.com", 10))
		addr, conninst, _ := ups.NextWithConnect()

		if addr != "http://baidu.com" {
			t.Errorf("error addr %s", addr)
		}

		client := conninst.(*http.Client)
		// 创建HTTP请求
		req, err := http.NewRequest("GET", "http://baidu.com", nil)
		if err != nil {
			t.Errorf("error %s", err.Error())
		}

		// 发送HTTP请求
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("error %s", err.Error())
		}

		// 处理HTTP响应
		defer resp.Body.Close()

		fmt.Println(resp)
		if resp.StatusCode != 200 {
			t.Errorf("error %d", resp.StatusCode)
		}

	})
}
