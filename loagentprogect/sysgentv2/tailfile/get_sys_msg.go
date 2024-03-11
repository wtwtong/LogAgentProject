package tailfile

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"sysagentv2/kafkas"
	"sysagentv2/model"
)

// insert
// 将CPU信息写入kafka中
func writesCpuPoints(cpuInfo *model.CpuInfo) {
	byts, err := json.Marshal(cpuInfo.CpuPercent)
	if err != nil {
		panic(err)
	}
	//封装消息，异步处理
	msg := &sarama.ProducerMessage{}
	msg.Topic = "sysCpu"
	msg.Value = sarama.StringEncoder(byts)
	fmt.Println("sysCpu发送消息")
	kafkas.MsgChan(msg)
}

// 将内存信息写入kafka中
func writesMemPoints(memInfo *model.MenInfo) {
	byts, err := json.Marshal(memInfo)
	if err != nil {
		panic(err)
	}
	//将数据发送到kafka
	//利用通道将同步的代码改为异步
	//把读出来的一行日志包装成kafka里面的msg类型
	msg := &sarama.ProducerMessage{}
	msg.Topic = "sysMem" //每个tailObj自己的topic
	msg.Value = sarama.StringEncoder(byts)
	fmt.Println("sysMem发送消息")
	//丢到通道中
	kafkas.MsgChan(msg)

}

// 将磁盘信息写入kafka中
func writesDiskPoints(diskInfo *model.DiskInfo) {
	byts, err := json.Marshal(diskInfo)
	if err != nil {
		panic(err)
	}
	//将数据发送到kafka
	//利用通道将同步的代码改为异步
	//把读出来的一行日志包装成kafka里面的msg类型
	msg := &sarama.ProducerMessage{}
	msg.Topic = "sysDisk" //每个tailObj自己的topic
	msg.Value = sarama.StringEncoder(byts)
	//丢到通道中
	fmt.Println("sysDisk发送消息")
	kafkas.MsgChan(msg)
}

// 将网卡信息写入kafka中
func writesNetPoints(netInfo *model.NetInfo) {
	byts, _ := json.Marshal(netInfo)
	//将数据发送到kafka
	//利用通道将同步的代码改为异步
	//把读出来的一行日志包装成kafka里面的msg类型
	msg := &sarama.ProducerMessage{}
	msg.Topic = "sysNet" //每个tailObj自己的topic
	msg.Value = sarama.StringEncoder(byts)
	fmt.Println("sysNet发送消息")
	//丢到通道中
	kafkas.MsgChan(msg)
}
