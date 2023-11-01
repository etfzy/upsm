package nacosresolver

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	upstreams "github.com/xpizy2020/upsm"
	"github.com/xpizy2020/upsm/loadbalancer"
	"github.com/xpizy2020/upsm/resolver"
)

type Watcher struct {
	Server  string
	Group   string
	Cluster string
}

func CreateWrrNacosResolver() resolver.Resolver {
	return &NacosResolver{}
}

type NacosResolver struct {
	nacosClient naming_client.INamingClient
}

func (c *NacosResolver) Init(optFns ...interface{}) error {
	fmt.Println(optFns)
	options := loadOptions(optFns)
	clientConfig := constant.ClientConfig{
		NamespaceId:         options.NameSpace,
		TimeoutMs:           uint64(options.TimeOutMs),
		ListenInterval:      uint64(options.ListenInterval),
		NotLoadCacheAtStart: options.NotLoadCacheAtStart,
		BeatInterval:        int64(options.BeatInterval),
		Username:            options.Username,
		Password:            options.Password,
		LogDir:              options.LogDir,
		CacheDir:            options.CacheDir,
		LogLevel:            "error",
		OpenKMS:             false,
	}

	serverConfigs := []constant.ServerConfig{}
	for _, v := range options.Servers {
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			Scheme:      "http",
			ContextPath: "/nacos",
			IpAddr:      v,
			Port:        uint64(options.Port),
		})
	}

	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	if err != nil {
		return err
	}
	c.nacosClient = namingClient
	return nil
}

func (c *NacosResolver) Subscribe(key interface{}, publish upstreams.Publish) error {

	conf := key.(*Watcher)
	err := c.nacosClient.Subscribe(&vo.SubscribeParam{
		ServiceName: conf.Server,
		GroupName:   conf.Group,             // default value is DEFAULT_GROUP
		Clusters:    []string{conf.Cluster}, // default value is DEFAULT
		SubscribeCallback: func(services []model.SubscribeService, err error) {

			if len(services) == 0 {
				return
			}
			nodes := make([]loadbalancer.AddrInfos, 0, len(services))
			for _, v := range services {
				if v.Healthy == false || v.Weight == 0 {
					continue
				}
				temp := &loadbalancer.WrrAddr{
					Weight: int(v.Weight),
					Addr:   fmt.Sprintf("%s:%d", v.Ip, v.Port),
				}

				nodes = append(nodes, temp)
			}

			publish(nodes)
		},
	})

	return err
}
func (c *NacosResolver) UnSubscribe(key interface{}) error {
	conf := key.(*Watcher)
	c.nacosClient.Unsubscribe(&vo.SubscribeParam{
		ServiceName: conf.Server,
		GroupName:   conf.Group,             // default value is DEFAULT_GROUP
		Clusters:    []string{conf.Cluster}, // default value is DEFAULT
	})
	return nil
}
