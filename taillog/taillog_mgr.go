package taillog

import (
	"fmt"
	"logagent/etcd"
)

type taillogMgr struct {
	logEntryConf []*etcd.LogEntry
}

var tskMgr *taillogMgr

// Init ...
func Init(logEntryConf []*etcd.LogEntry) {
	tskMgr = &taillogMgr{
		logEntryConf: logEntryConf,
	}
	for _, logEntry := range tskMgr.logEntryConf {
		fmt.Printf("path:%v topic:%v\n", logEntry.Path, logEntry.Topic)
		NewTailTask(logEntry.Path, logEntry.Topic)
	}
}
