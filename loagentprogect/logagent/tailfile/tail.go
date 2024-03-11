package tailfile

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/nxadm/tail"
	"github.com/sirupsen/logrus"
	"logagent/kafka"
	"strings"
	"time"
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
		select {
		case <-t.ctx.Done(): //只要调用tail.cancel就会接收到信号
			logrus.Infof("path %s is stopping...", t.path)
			return
		//循环读数据
		case line, ok := <-t.TailObj.Lines:
			if !ok {
				logrus.Warn("tail file close reopen, filename:%s\n", t.TailObj.Filename)
				time.Sleep(time.Second) // 读取出错等1秒
				continue
			}
			//去掉空行
			if len(strings.Trim(line.Text, "\r")) == 0 {
				logrus.Info("出现空行，跳过")
				continue
			}
			//利用通道将同步的代码改为异步
			//把读出来的一行日志包装成kafka里面的msg类型
			msg := &sarama.ProducerMessage{}
			msg.Topic = t.topic //每个tailObj自己的topic
			msg.Value = sarama.StringEncoder(line.Text)
			//丢到通道中
			kafka.MsgChan(msg)
		}
	}
}
