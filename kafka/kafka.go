package kafka

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

var (
	client  sarama.SyncProducer // 声明一个全局的连接kafka的生产者client
	logChan chan *logData
)

// logData ...
type logData struct {
	topic string
	data  string
}

// Init 初始化client
func Init(addrs []string, maxSize int) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	// 连接kafka
	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	// 初始化channel
	logChan = make(chan *logData, maxSize)
	// 开启一个后台的goroutine循环往kafka发送消息
	go sendToKafka()
	return
}

// SendToChan 给外部暴露一个函数，将消息发送到内部的通道中
func SendToChan(topic, data string) {
	msg := &logData{
		topic: topic,
		data:  data,
	}
	logChan <- msg
}

// SendToKafka 将消息从通道中取出后发送到kafka
func sendToKafka() {
	for {
		select {
		case ld := <-logChan:
			// 构造一个消息
			msg := &sarama.ProducerMessage{}
			msg.Topic = ld.topic
			msg.Value = sarama.StringEncoder(ld.data)
			// 发送到Kafka
			ptn, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println("send msg failed, err:", err)
				return
			}
			fmt.Printf("ptn:%v offset:%v\n", ptn, offset)
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}

}
