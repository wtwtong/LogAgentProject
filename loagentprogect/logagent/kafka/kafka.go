package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	client  sarama.SyncProducer
	msgChan chan *sarama.ProducerMessage
)

func Init(address []string, chanSize int64) (err error) {
	//配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          //ACK
	config.Producer.Partitioner = sarama.NewRandomPartitioner //分区
	config.Producer.Return.Successes = true                   //确认
	//连接kafka
	client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		logrus.Error("kafka:producer closed err", err)
		return
	}
	//发送消息，用一个管道，将读日志和发送日志改为异步执行
	msgChan = make(chan *sarama.ProducerMessage, chanSize)
	//后台启动一个goroutine从MsgChan读取数据
	go sendMsg()
	return
}

// 从MsgChan中取出数据，发送到kafka
func sendMsg() {
	for {
		select {
		case msg := <-msgChan:
			paten, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println("send message err:", err)
				return
			}
			logrus.Infof("send message paten %v , offset %v", paten, offset)
			fmt.Println(msg.Value, msg.Topic)
		}
	}

}
func MsgChan(msg *sarama.ProducerMessage) {
	msgChan <- msg
}
