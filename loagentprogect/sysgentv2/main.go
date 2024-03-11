package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"sysagentv2/common"
	"sysagentv2/etcd"
	"sysagentv2/kafkas"
	"sysagentv2/model"
	"sysagentv2/tailfile"
)

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
	configObj := new(model.Config)
	//0 读取配置文件 `go-ini`
	err = ini.MapTo(configObj, "./conf/config.ini")
	if err != nil {
		logrus.Error(" file MapTo err:", err)
		return
	}
	//fmt.Printf("%#v\n", configObj)

	//1 初始化kafka
	err = kafkas.Init([]string{configObj.KafkaConfig.Address}, configObj.KafkaConfig.ChanSize)
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
