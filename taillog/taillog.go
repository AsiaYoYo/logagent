package taillog

import (
	"fmt"
	"logagent/kafka"
	"time"

	"github.com/hpcloud/tail"
)

// 专门从日志文件收集日志的模块

var (
	tailTask *TailTask
)

// TailTask 存放每个tailObj的path topic instance
type TailTask struct {
	path     string
	topic    string
	instance *tail.Tail
}

// NewTailTask TailTask结构体构造函数
func NewTailTask(path, topic string) (tailTask *TailTask) {
	tailTask = &TailTask{
		path:  path,
		topic: topic,
	}
	tailTask.init()
	return
}

// Init 打开日志文件
func (t *TailTask) init() (err error) {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		fmt.Printf("tail file failed, err:%v\n", err)
		return
	}
	// 开启一个goroutine循环取日志
	go t.run()
	return
}

// run ...
func (t *TailTask) run() {
	for {
		select {
		case line := <-t.instance.Lines:
			// 收集日志发往kafka
			// kafka.SendToKafka(t.topic, line.Text) // 这里存在函数调函数，可能影响程序效率
			kafka.SendToChan(t.topic, line.Text)
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}
