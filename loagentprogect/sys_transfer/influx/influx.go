package influx

import (
	"LogProject/sys_transfer/model"
	"encoding/json"
	"github.com/Shopify/sarama"
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"time"
)

var (
	cli client.Client
)

func InitConnInflux(addr, username, password string) (err error) {
	cli, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: username,
		Password: password,
	})
	if err != nil {
		return err
	}
	log.Println("connect influx success")
	return
}

//将信息发往 influx

func ProcessCpuData(message *sarama.ConsumerMessage) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	//var data model.CpuInfo
	//err = json.Unmarshal(message.Value, &data)
	if err != nil {
		log.Println("Error Unmarshal cpu data:", err)
		return
	}
	tags := map[string]string{"cpu": "cpu0"}
	fields := map[string]interface{}{
		"cpu_percent": message.Value,
	}

	pt, err := client.NewPoint("cpu_percent", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)
	err = cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert cpu info success")
}
func ProcessMemData(message *sarama.ConsumerMessage) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	var data model.MenInfo
	err = json.Unmarshal(message.Value, &data)
	if err != nil {
		log.Println("Error Unmarshal mem data:", err)
		return
	}
	tags := map[string]string{"mem": "mem"}
	fields := map[string]interface{}{
		"total":        int64(data.Total),
		"available":    int64(data.Available),
		"used":         int64(data.Used),
		"used_percent": data.UsedPercent,
	}

	pt, err := client.NewPoint("memory", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)
	err = cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert mem info success")
}
func ProcessDiskData(message *sarama.ConsumerMessage) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	var data model.DiskInfo
	err = json.Unmarshal(message.Value, &data)
	if err != nil {
		log.Println("Error Unmarshal disk data:", err)
		return
	}
	//根据传入的数据类型插入数据
	for k, v := range data.PartitionUsageStat {
		tags := map[string]string{"path": k}
		fields := map[string]interface{}{
			"total":             int64(v.Total),
			"free":              int64(v.Free),
			"used":              int64(v.Used),
			"used_percent":      v.UsedPercent,
			"inodesTotal":       int64(v.InodesTotal),
			"inodesUsed":        int64(v.InodesUsed),
			"InodesFree":        int64(v.InodesFree),
			"InodesUsedPercent": v.InodesUsedPercent,
		}
		pt, err := client.NewPoint("disk", tags, fields, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}
	err = cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert Disk info success")
}
func ProcessNetData(message *sarama.ConsumerMessage) {
	var data model.NetInfo
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(message.Value, &data)
	if err != nil {
		log.Println("Error unmarshaling Net data:", err)
		return
	}
	//根据传入的数据类型插入数据
	for k, v := range data.NetIOCountersStat {
		tags := map[string]string{"name": k}
		fields := map[string]interface{}{
			"bytesSentRate":   v.BytesSentRate,
			"bytesRecvRate":   v.BytesRecvRate,
			"packetsSentRate": v.PacketsSentRate,
			"packetsRecvRate": v.PacketsRecvRate,
		}
		pt, err := client.NewPoint("net", tags, fields, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}
	err = cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert Net info success")
}
