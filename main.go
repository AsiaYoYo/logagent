package main

import (
	"fmt"
	"time"

	"gopkg.in/ini.v1"

	"logagent/conf"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/taillog"
)

var cfg = new(conf.AppConf)

func run() {
	select {}
}

func main() {
	// 0.加载配置文件
	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Printf("init config failed, err:%v\n", err)
		return
	}
	// 1.初始化kafka连接
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.ChanMaxSzie)
	if err != nil {
		fmt.Printf("init kafka failed, err:%v\n", err)
		return
	}
	fmt.Println("init Kafka success")
	// 2.初始化etcd
	err = etcd.Init(cfg.EtcdConf.Address, cfg.EtcdConf.Timeout*time.Second)
	if err != nil {
		fmt.Printf("init kafka failed, err:%v\n", err)
		return
	}
	fmt.Println("init etcd success")
	// 2.1 从etcd中获取日志收集的配置信息
	logEntryConf, err := etcd.GetConf(cfg.EtcdConf.Key)
	if err != nil {
		fmt.Printf("etcd.GetConf failed, err:%v\n", err)
		return
	}
	fmt.Printf("get conf from etcd success, %v\n", logEntryConf)
	// 2.2 派一个哨兵监听配置

	// 3.收集日志发往kafka
	// 3.1 循环每一个日志收集项
	taillog.Init(logEntryConf)

	run()
}
