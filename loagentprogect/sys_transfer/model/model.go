package model

import "github.com/shirou/gopsutil/disk"

type Config struct {
	KafkaConf  `ini:"kafka"`
	InfluxConf `ini:"influx"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}
type InfluxConf struct {
	Address      string `ini:"address"`
	Username     string `ini:"username"`
	Password     string `ini:"password"`
	MaxChan      int    `ini:"max_chan_size"`
	GoroutineNum int    `ini:"goroutine_num"`
}
type CpuInfo struct {
	CpuPercent float64 `json:"cpu_percent"`
}
type MenInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
	Buffers     uint64  `json:"buffers"`
	Cached      uint64  `json:"cached"`
}

type DiskInfo struct {
	PartitionUsageStat map[string]*disk.UsageStat
}
type IOStat struct {
	BytesSent       uint64
	BytesRecv       uint64
	PacketsSent     uint64
	PacketsRecv     uint64
	BytesSentRate   float64 `json:"bytesSentRate"`
	BytesRecvRate   float64 `json:"bytesRecvRate"`
	PacketsSentRate float64 `json:"packetsSentRate"`
	PacketsRecvRate float64 `json:"packetsRecvRate"`
}
type NetInfo struct {
	NetIOCountersStat map[string]*IOStat
}
