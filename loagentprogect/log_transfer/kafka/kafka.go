package kafka

import (
	"LogProject/log_transfer/es"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
)

func Init(address []string, topic string) (err error) {
	//创建一个consumer
	consumer, err := sarama.NewConsumer(address, nil)
	if err != nil {
		fmt.Println("New Consumer err", err)
		return
	}
	//对consumer进行分区
	partitionsLiat, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Println("file get to partition list err:", err)
		return
	}
	//遍历分区
	for partition := range partitionsLiat {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Println("consumer for partition err:", err)
			panic(err)
		}
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				var m1 map[string]interface{}
				err = json.Unmarshal(msg.Value, &m1)
				if err != nil {
					fmt.Println("Unmarshal msg file,err:", err)
					continue
				}
				es.PutLogData(m1)
			}
		}(pc)
	}
	return
}
