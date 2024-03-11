package main

import (
	"fmt"
	"github.com/nxadm/tail"
	"time"
)

// tail包读取日志文件
func main() {
	filename := `E:/Go_Code/kafka_logs/web-0/00000000000000000000.log`
	//配置
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	//打开文件开始读取数据
	tails, err := tail.TailFile(filename, config)
	if err != nil {
		fmt.Printf("tail %s failed,err:%v\n", tails, err)
		return
	}
	//	开始读取数据
	var (
		msg *tail.Line
		ok  bool
	)
	for {
		msg, ok = <-tails.Lines //chan tail.line
		if !ok {
			fmt.Printf("tail file close reopen,filename:%s\n", tails.Filename)
			time.Sleep(time.Second) //读取出错等1秒
			continue
		}
		fmt.Println("msg:", msg.Text)
	}
}
