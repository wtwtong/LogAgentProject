package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

/*
	func main() {
		//1.生产者配置
		config := sarama.NewConfig()
		config.Producer.RequiredAcks = sarama.WaitForAll          //ACk
		config.Producer.Partitioner = sarama.NewRandomPartitioner //分区
		config.Producer.Return.Successes = true                   //确认

		//2.链接kafka
		clint, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
		if err != nil {
			fmt.Println("Producer closed err:", err)
			return
		}
		defer clint.Close()

		//3.封装消息
		msg := &sarama.ProducerMessage{}
		msg.Topic = "web_log"
		msg.Value = sarama.StringEncoder("2023.11.13 s4 very happy--!")

		//4.发送消息
		pid, offset, err := clint.SendMessage(msg)
		if err != nil {
			fmt.Println("send message err:", err)
			return
		}
		fmt.Printf("pid:%v  offset:%v\n", pid, offset)

}
*/
func main() {
	//配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          //ACK
	config.Producer.Partitioner = sarama.NewRandomPartitioner //分区
	config.Producer.Return.Successes = true                   //确认
	//连接kafka
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		fmt.Println("Producer closed err:", err)
		return
	}
	defer client.Close()
	//封装消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = "web"
	msg.Value = sarama.StringEncoder("2024.1.19 wang tong!")
	//	发送消息
	paten, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send message err:", err)
		return
	}
	fmt.Printf("paten:%v  offset:%v\n", paten, offset)

}
