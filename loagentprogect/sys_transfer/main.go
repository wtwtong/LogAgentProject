package main

import (
	"LogProject/sys_transfer/influx"
	"LogProject/sys_transfer/kafka"
	"LogProject/sys_transfer/model"
	"fmt"
	"gopkg.in/ini.v1"
	"log"
)

// sys transfer
// 从kafka中消费数据，发往influxdb
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
	//连接influx
	err = influx.InitConnInflux(configObj.InfluxConf.Address, configObj.InfluxConf.Username, configObj.InfluxConf.Password)
	if err != nil {
		log.Fatal(err)
	}

	//连接kafka
	err = kafka.Init([]string{configObj.KafkaConf.Address})
	if err != nil {
		fmt.Println("kafka init filed, err:", err)
		panic(err)
	}
	fmt.Println("kafka init success")

	for {
		select {}
	}
}
