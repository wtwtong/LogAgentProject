package model

import "github.com/shirou/gopsutil/disk"

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
