package main

import (
	"LogProject/log_transfer/es"
	"LogProject/log_transfer/kafka"
	"LogProject/log_transfer/model"
	"fmt"
	"gopkg.in/ini.v1"
)

// log transfer
// 从kafka中消费数据，发往ES
func main() {
	//	初始化配置文件
	configObj := new(model.Config)
	err := ini.MapTo(configObj, "./conf/config.ini")
	if err != nil {
		fmt.Println("conf init filed, err:", err)
		panic(err)
	}
	fmt.Println(configObj)
	fmt.Println("conf init success")
	//连接ES
	err = es.Init(configObj.EsConf.Address, configObj.EsConf.Index, configObj.EsConf.MaxChan, configObj.EsConf.GoroutineNum)
	if err != nil {
		fmt.Println("ES init filed, err:", err)
		panic(err)
	}
	fmt.Println("ES init success")
	//连接kafka
	err = kafka.Init([]string{configObj.KafkaConf.Address}, configObj.KafkaConf.Topic)
	if err != nil {
		fmt.Println("kafka init filed, err:", err)
		panic(err)
	}
	fmt.Println("kafka init success")

	for {
		select {}
	}
}
