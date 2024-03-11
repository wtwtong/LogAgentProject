package tailfile

import (
	"context"
	"fmt"
	"github.com/nxadm/tail"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/sirupsen/logrus"
	"sysagentv2/model"
	"time"
)

var (
	lastNetIoStatTimeStamp int64
	lastNetInfo            *model.NetInfo
)

type tailTask struct {
	path    string
	topic   string
	TailObj *tail.Tail
	ctx     context.Context
	cancel  context.CancelFunc
}

func newTailTask(path, topic string) *tailTask {
	ctx, cancel := context.WithCancel(context.Background())
	tt := &tailTask{
		path:   path,
		topic:  topic,
		ctx:    ctx,
		cancel: cancel,
	}
	return tt
}
func (t *tailTask) Init() (err error) {
	//配置
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	//打开文件开始读取数据
	t.TailObj, err = tail.TailFile(t.path, config)
	if err != nil {
		logrus.Error("tail init path %v  err:%v", t.TailObj, err)
		return
	}
	return
}
func (t *tailTask) run() {
	//读取日志发往kafka
	logrus.Infof("collect for path:%s is running...", t.path)
	for {
		runSysMsg(time.Second)

	}
}

func getCpuInfo() {
	var cpuInfo = new(model.CpuInfo)
	// CPU使用率
	percent, _ := cpu.Percent(time.Second, false)
	//	写入到kafka中
	cpuInfo.CpuPercent = percent[0]
	writesCpuPoints(cpuInfo)
}

func getMemInfo() {
	var memInfo = new(model.MenInfo)
	// CPU使用率
	info, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("get mem info filed, err:", err)
		return
	}
	//	写入到kafka中
	memInfo.Total = info.Total
	memInfo.Used = info.Used
	memInfo.Available = info.Available
	memInfo.UsedPercent = info.UsedPercent
	memInfo.Buffers = info.Buffers
	memInfo.Cached = info.Cached
	writesMemPoints(memInfo)
}

func getDiskInfo() {
	var diskInfo = &model.DiskInfo{
		PartitionUsageStat: make(map[string]*disk.UsageStat, 16),
	}
	parts, _ := disk.Partitions(true)
	for _, part := range parts {
		usageStat, err := disk.Usage(part.Mountpoint)
		if err != nil {
			fmt.Printf("get Partitions failed, err:%v\n", err)
			continue
		}
		diskInfo.PartitionUsageStat[part.Mountpoint] = usageStat
	}
	writesDiskPoints(diskInfo)
}
func getNetInfo() {
	var netInfo = &model.NetInfo{
		NetIOCountersStat: make(map[string]*model.IOStat, 8),
	}
	currentTimeStamp := time.Now().Unix()
	info, _ := net.IOCounters(true)
	for _, netIo := range info {
		iostat := new(model.IOStat)
		iostat.BytesSent = netIo.BytesSent
		iostat.BytesRecv = netIo.BytesRecv
		iostat.PacketsSent = netIo.PacketsSent
		iostat.PacketsRecv = netIo.PacketsRecv
		//保存数据
		netInfo.NetIOCountersStat[netIo.Name] = iostat
		//	判断
		if lastNetIoStatTimeStamp == 0 || lastNetInfo == nil {
			continue
		}
		//	计算速率
		interval := currentTimeStamp - lastNetIoStatTimeStamp
		iostat.BytesSentRate = (float64(iostat.BytesSent) - float64(lastNetInfo.NetIOCountersStat[netIo.Name].BytesSent)) / float64(interval)
		iostat.BytesRecvRate = (float64(iostat.BytesRecv) - float64(lastNetInfo.NetIOCountersStat[netIo.Name].BytesRecv)) / float64(interval)
		iostat.PacketsSentRate = (float64(iostat.PacketsSent) - float64(lastNetInfo.NetIOCountersStat[netIo.Name].PacketsSent)) / float64(interval)
		iostat.PacketsRecvRate = (float64(iostat.PacketsRecv) - float64(lastNetInfo.NetIOCountersStat[netIo.Name].PacketsRecv)) / float64(interval)
	}
	//更新全局记录的上一次采集的时间点和网卡数据
	lastNetInfo = netInfo
	lastNetIoStatTimeStamp = currentTimeStamp
	writesNetPoints(netInfo)
}

func runSysMsg(interval time.Duration) {
	ticker := time.Tick(interval)
	for _ = range ticker {
		getCpuInfo()
		getMemInfo()
		getDiskInfo()
		getNetInfo()
	}
}
