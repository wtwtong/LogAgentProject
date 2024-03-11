package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"sync"
)

func main() {
	//创建一个consumer
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		fmt.Println("New Consumer err", err)
		return
	}
	//对consumer进行分区
	partitionsLiat, err := consumer.Partitions("web")
	if err != nil {
		fmt.Println("file get to partition list err:", err)
		return
	}
	var wg sync.WaitGroup
	//遍历分区
	for partition := range partitionsLiat {

		pc, err := consumer.ConsumePartition("web", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Println("consumer for partition err", err)
			return
		}
		wg.Add(1)
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d,Topic:%s,Key:%s,Value:%s\n", msg.Partition, msg.Topic, msg.Key, msg.Value)
			}
		}(pc)
	}
	wg.Wait()
}
