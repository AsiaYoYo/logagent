package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

var (
	cli *clientv3.Client
)

// LogEntry ...
type LogEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

// Init 初始化etcd的函数
func Init(addr string, timeout time.Duration) (err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: timeout,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	return
}

// GetConf 从etcd中取配置
func GetConf(key string) (LogEntryConf []*LogEntry, err error) {
	// get
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		// fmt.Printf("%s:%s\n", ev.Key, ev.Value)
		err = json.Unmarshal(ev.Value, &LogEntryConf)
		if err != nil {
			fmt.Printf("Unmarshal etcd value failed, err:%v\n", err)
			return
		}
	}
	return
}

// WatchConf 监听未来配置文件的改动
func WatchConf(key string, newConfChan chan<- []*LogEntry) {
	// watch key change
	rch := cli.Watch(context.Background(), key) // <-chan WatchResponse
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			if ev.Type != clientv3.EventTypeDelete {
				// fmt.Printf("ev.Type: %v, ev.Kv.Value: %v\n", ev.Type, ev.Kv.Value)
				var newConf []*LogEntry
				json.Unmarshal(ev.Kv.Value, &newConf)
				newConfChan <- newConf
			}
		}
	}
}
