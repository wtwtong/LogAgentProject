package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"logagent/common"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/tailfile"
)

type Config struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
	EtcdConfig    `ini:"etcd"`
	InfluxDB      `ini:"influxdb"`
}
type KafkaConfig struct {
	Address  string `ini:"address"`
	ChanSize int64  `ini:"chan_size"`
}
type CollectConfig struct {
	CollectPath string `ini:"log_path"`
}
type EtcdConfig struct {
	Address   string `ini:"address"`
	ClientKey string `ini:"client_key"`
}
type InfluxDB struct {
	Address  string `ini:"address"`
	UserName string `ini:"username"`
	Password string `ini:"password"`
}

/*var (
	ok  bool
	msg *tail.Line
)*/

func run() {
	//TailObj  -->log -->Client --> kafka

	for {
		select {}
	}
}

func main() {
	//获取本机IP,为后续去etcd取配置文件打下基础
	ip, err := common.GetOutboundIP()
	if err != nil {
		logrus.Errorf("get ip faild,err:%v", err)
		return
	}
	configObj := new(Config)
	//0 读取配置文件 `go-ini`
	err = ini.MapTo(configObj, "./conf/config.ini")
	if err != nil {
		logrus.Error(" file MapTo err:", err)
		return
	}
	//fmt.Printf("%#v\n", configObj)

	//1 初始化kafka
	err = kafka.Init([]string{configObj.KafkaConfig.Address}, configObj.KafkaConfig.ChanSize)
	if err != nil {
		logrus.Error(" kafka ini err:", err)
		return
	}
	logrus.Info("client success...")
	//连接etcd
	err = etcd.Init([]string{configObj.EtcdConfig.Address})
	if err != nil {
		logrus.Error("ini etcd flied err:", err)
		return
	}
	//从etcd中拉去要收集日志的配置项
	clientKey := fmt.Sprintf(configObj.EtcdConfig.ClientKey, ip)
	allConf, err := etcd.GetConf(clientKey)
	if err != nil {
		logrus.Error("get etcd conf err:", err)
		return
	}
	//fmt.Println(allConf)
	//派一个小弟去监控etcd中clientKey的变化
	go etcd.WatchConf(clientKey)
	//	2 根据配置中的日志文件初始化tail
	err = tailfile.Init(allConf)
	if err != nil {
		logrus.Error(" init tail filed err:", err)
		return
	}
	logrus.Info("init tail success")
	//	使用sarama往kafka中发数据
	run()

}
