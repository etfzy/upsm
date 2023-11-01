package nacosresolver

type Options struct {
	Servers             []string
	NameSpace           string
	TimeOutMs           int
	ListenInterval      int
	NotLoadCacheAtStart bool
	BeatInterval        int
	Username            string
	Password            string
	LogDir              string
	CacheDir            string
	Port                int
}

type OptionFns func(opt *Options)

func WithDailTimeOut(timout int) OptionFns {
	return func(opt *Options) {
		opt.TimeOutMs = timout
	}
}

func WithServers(servers []string) OptionFns {
	return func(opt *Options) {
		opt.Servers = servers
	}

}

func WithPort(port int) OptionFns {
	return func(opt *Options) {
		opt.Port = port
	}
}

func WithNameSpace(nameSpace string) OptionFns {
	return func(opt *Options) {
		opt.NameSpace = nameSpace
	}

}

func WithListenInterval(interval int) OptionFns {
	return func(opt *Options) {
		opt.ListenInterval = interval
	}

}

func WithNotLoadCacheAtStart(bStart bool) OptionFns {
	return func(opt *Options) {
		opt.NotLoadCacheAtStart = bStart
	}

}

func WithBeatInterval(interval int) OptionFns {
	return func(opt *Options) {
		opt.BeatInterval = interval
	}

}

func WithUsername(user string) OptionFns {
	return func(opt *Options) {
		opt.Username = user
	}

}

func WithPwd(pwd string) OptionFns {
	return func(opt *Options) {
		opt.Password = pwd
	}

}

func WithLogDir(dir string) OptionFns {
	return func(opt *Options) {
		opt.LogDir = dir
	}

}

func WithCacheDir(dir string) OptionFns {
	return func(opt *Options) {
		opt.CacheDir = dir
	}

}

func loadOptions(opts interface{}) *Options {

	poolOpt := &Options{
		TimeOutMs:           1000,
		ListenInterval:      1000,
		NotLoadCacheAtStart: true,
		BeatInterval:        1000,
		LogDir:              "./nacoscache/register/log",
		CacheDir:            "./nacoscache/register/cache",
	}
	optslist := opts.([]interface{})
	for _, optFn := range optslist {
		temp := optFn.(OptionFns)
		temp(poolOpt)
	}

	return poolOpt
}
