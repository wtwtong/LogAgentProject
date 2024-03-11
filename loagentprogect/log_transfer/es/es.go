package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

type ESClient struct {
	client      *elastic.Client
	index       string
	LogDataChan chan interface{}
}

var (
	esClient *ESClient
)

func Init(addr, index string, maxsize, gorNum int) (err error) {
	client, err := elastic.NewClient(elastic.SetURL(addr))
	if err != nil {
		// Handle error
		panic(err)
	}
	esClient = &ESClient{
		client:      client,
		index:       index,
		LogDataChan: make(chan interface{}, maxsize),
	}
	fmt.Println("connect to es success")
	for i := 0; i < gorNum; i++ {
		//发送数据到ES
		go sendToES()
	}

	return
}
func sendToES() {
	for m1 := range esClient.LogDataChan {
		put1, err := esClient.client.Index().
			Index(esClient.index).
			BodyJson(m1).
			Do(context.Background())
		if err != nil {
			// Handle error
			panic(err)
		}
		fmt.Printf("Indexed user %s to index %s,type %s\n", put1.Id, put1.Index, put1.Type)
	}
}

// 拿到包外的msg发送到channel中

func PutLogData(msg map[string]interface{}) {
	esClient.LogDataChan <- msg
}
