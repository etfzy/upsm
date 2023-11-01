package loadbalancer

type WrrAddr struct {
	Addr      string
	Weight    int
	CurWeight int
}

func (w *WrrAddr) GetAddr() string {
	return w.Addr
}

func (w *WrrAddr) BuildAddr() AddrInfos {
	return w
}

func CreateWrrServer(addr string, weight int) *WrrAddr {
	return &WrrAddr{
		Addr:   addr,
		Weight: weight,
	}
}

type Wrr struct {
	builder Builder
}

func DefaultWrrBuilder(info interface{}) AddrInfos {
	temp := info.(*WrrAddr)
	return temp
}

func NewWrrLoadBalancer(builder Builder) LoadBalancer {
	wrrlb := &Wrr{
		builder: builder,
	}
	return wrrlb
}

func (w *Wrr) BuildAddr(info interface{}) AddrInfos {
	return w.builder(info)
}

func (w *Wrr) Picker(servers []AddrInfos) string {

	if servers == nil {
		return ""
	}

	if len(servers) == 0 {
		return ""
	}

	var best *WrrAddr
	var totalWeight int = 0

	for _, temp := range servers {

		server := temp.(*WrrAddr)
		totalWeight += server.Weight

		server.CurWeight += server.Weight

		if best == nil {
			best = server
		} else if server.CurWeight > best.CurWeight {
			best = server
		}

	}

	if best == nil {
		return ""

	}

	best.CurWeight -= totalWeight

	return best.Addr
}
