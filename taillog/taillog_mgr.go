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
		// 初始化的时候起了多少个tailtask都要记下来
		tailObj := NewTailTask(logEntry.Path, logEntry.Topic)
		mk := fmt.Sprintf("%s_%s", logEntry.Path, logEntry.Topic)
		tskMgr.tskMap[mk] = tailObj
	}
	// 后台开启goroutine取新的配置
	go tskMgr.run()
}

// 监听自己的NewConfChan，有了新的配置做对应处理
func (t *taillogMgr) run() {
	for {
		select {
		case newConf := <-tskMgr.newConfChan:
			for _, conf := range newConf {
				mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
				_, ok := t.tskMap[mk]
				if ok {
					// 如果存在，不做操作
					continue
				} else {
					// 如果不存在，则视为新增
					fmt.Printf("%s_%s, 新增任务\n", conf.Path, conf.Topic)
					tailObj := NewTailTask(conf.Path, conf.Topic)
					tskMgr.tskMap[mk] = tailObj
				}
			}
			// 2.配置删除
			// 找出t.logEntryConf中有的，和newConf中没有的，将其删掉
			for _, c1 := range t.logEntryConf {
				isDelete := true
				for _, c2 := range newConf {
					if c1.Path == c2.Path && c1.Topic == c2.Topic {
						isDelete = false
						continue
					}
				}
				if isDelete {
					// 把c1对应的tailObj删掉
					mk := fmt.Sprintf("%s_%s", c1.Path, c1.Topic)
					t.tskMap[mk].cancelFunc()
					// 从tskMap中也删除该任务
					delete(t.tskMap, mk)
				}
			}
			// 3.配置变更
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
