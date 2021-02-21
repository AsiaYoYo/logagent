package conf

import (
	"time"
)

// AppConf ...
type AppConf struct {
	KafkaConf `ini:"kafka"`
	EtcdConf  `ini:"etcd"`
}

// KafkaConf ...
type KafkaConf struct {
	Address     string `ini:"address"`
	ChanMaxSzie int    `ini:"chan_max_size"`
}

// EtcdConf ...
type EtcdConf struct {
	Address string        `ini:"address"`
	Key     string        `ini:"key"`
	Timeout time.Duration `ini:"timeout"`
}
