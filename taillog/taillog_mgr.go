package taillog

import (
	"fmt"
	"logagent/etcd"
	"time"
)

type taillogMgr struct {
	logEntryConf []*etcd.LogEntry
	tskMap       map[string]*TailTask
	newConfChan  chan []*etcd.LogEntry
}

var tskMgr *taillogMgr

// Init ...
func Init(logEntryConf []*etcd.LogEntry) {
	tskMgr = &taillogMgr{
		logEntryConf: logEntryConf,
		tskMap:       make(map[string]*TailTask, 32),
		newConfChan:  make(chan []*etcd.LogEntry),
	}
	for _, logEntry := range tskMgr.logEntryConf {
		fmt.Printf("path:%v topic:%v\n", logEntry.Path, logEntry.Topic)
		NewTailTask(logEntry.Path, logEntry.Topic)
	}
	// 后台开启goroutine取新的配置
	go tskMgr.run()
}

// 监听自己的NewConfChan，有了新的配置做对应处理
func (t *taillogMgr) run() {
	for {
		select {
		case newConf := <-tskMgr.newConfChan:
			fmt.Printf("新配置来了,newConf: %v\n", newConf)
		default:
			time.Sleep(time.Second)
		}
	}
}

// NewConfChan 向外提供一个通道
func NewConfChan() (NewConf chan<- []*etcd.LogEntry) {
	return tskMgr.newConfChan
}
