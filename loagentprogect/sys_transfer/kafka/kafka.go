package kafka

import (
	"LogProject/sys_transfer/influx"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
)

func Init(address []string) (err error) {
	client, err := sarama.NewClient(address, nil)
	if err != nil {
		fmt.Println("Error creating client:", err)
		return
	}
	defer client.Close()

	topics, err := client.Topics()
	if err != nil {
		log.Fatal(err)
	}

	for _, topic := range topics {
		//创建一个consumer

		consumer, err := sarama.NewConsumer(address, nil)
		if err != nil {
			log.Fatal(err)
		}
		//对consumer进行分区
		partitionsLiat, err := consumer.Partitions(topic)
		if err != nil {
			log.Fatal(err)

		}

		//遍历分区
		for partition := range partitionsLiat {
			pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
			if err != nil {
				fmt.Println("consumer for partition err:", err)
				panic(err)
			}
			// 开启一个协程用于消费消息并发送到InfluxDB
			go func(pc sarama.PartitionConsumer) {
				// 初始化 channel 来传递消息
				for message := range pc.Messages() {
					fmt.Println("topic:", message.Topic)
					// 根据消息的不同topic执行不同的处理逻辑
					switch message.Topic {
					case "sysCpu":
						influx.ProcessCpuData(message)
					case "sysMem":
						influx.ProcessMemData(message)
					case "sysDisk":
						influx.ProcessDiskData(message)
					case "sysNet":
						influx.ProcessNetData(message)
					default:
						fmt.Println("Unknown topic:", message.Topic)
					}
				}
			}(pc)
		}
	}

	////对consumer进行分区
	//topics := []string{"sysCpu", "sysMem", "sysDisk", "sysNet"}
	//for _, topic := range topics {
	//	fmt.Println("topic1:", topic)
	//	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	//	if err != nil {
	//		log.Fatalf("Error consuming partition for topic %s: %v", topic, err)
	//	}
	//
	//	// 开启一个协程用于消费消息并发送到InfluxDB
	//	go func(pc sarama.PartitionConsumer) {
	//		// 初始化 channel 来传递消息
	//		messageChannel := make(chan *sarama.ConsumerMessage)
	//		fmt.Println("message:", messageChannel)
	//		for message := range messageChannel {
	//			fmt.Println("topic:", message.Topic)
	//			// 根据消息的不同topic执行不同的处理逻辑
	//			switch message.Topic {
	//			case "sysCpu":
	//				influx.ProcessCpuData(message)
	//			case "sysMem":
	//				influx.ProcessMemData(message)
	//			case "sysDisk":
	//				influx.ProcessDiskData(message)
	//			case "sysNet":
	//				influx.ProcessNetData(message)
	//			default:
	//				fmt.Println("Unknown topic:", message.Topic)
	//			}
	//		}
	//	}(partitionConsumer)
	//}
	return
}
