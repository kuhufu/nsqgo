package nsqgo

import (
	"github.com/nsqio/go-nsq"
	"time"
)

type Options struct {
	mode              string        //nsqd, nsqlookupd
	addrs             []string      //mode=nsqlookupd则为nsqlookupd地址，mode=nsqd则为nsqd地址
	nsqConfig         *nsq.Config   //nsq配置
	reconnectInterval time.Duration //重连间隔
	selectStrategy    Selector      //选择策略
}

func NewOptions() *Options {
	return &Options{
		nsqConfig:         nsq.NewConfig(),
		reconnectInterval: time.Second * 5,
		selectStrategy:    Random(),
	}
}

func WithNSQConfig(cfg *nsq.Config) Option {
	return func(opts *Options) {
		opts.nsqConfig = cfg
	}
}

func WithNSQLookupd(addrs ...string) Option {
	return func(opts *Options) {
		opts.mode = "nsqlookupd"
		opts.addrs = addrs
	}
}

func WithNSQDs(addrs ...string) Option {
	return func(opts *Options) {
		opts.mode = "nsqd"
		opts.addrs = addrs
	}
}

func WithSelector(selector Selector) Option {
	return func(opts *Options) {
		opts.selectStrategy = selector
	}
}

func WithReconnectInterval(duration time.Duration) Option {
	return func(opts *Options) {
		opts.reconnectInterval = duration
	}
}
