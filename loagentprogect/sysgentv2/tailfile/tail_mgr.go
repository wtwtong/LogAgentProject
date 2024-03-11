package tailfile

import (
	"github.com/sirupsen/logrus"
	"sysagentv2/common"
)

type tailTaskMgr struct {
	tailTaskMap      map[string]*tailTask
	collectEntryList []common.CollectEntry
	confChan         chan []common.CollectEntry
}

var (
	ttMgr *tailTaskMgr
)

func Init(allConf []common.CollectEntry) (err error) {
	//allConf里面存了若干个日志的收集项
	//针对每一个日志收集项创建一个tailObj
	ttMgr = &tailTaskMgr{
		tailTaskMap:      make(map[string]*tailTask, 20),
		collectEntryList: allConf,
		confChan:         make(chan []common.CollectEntry),
	}
	for _, conf := range allConf {
		//创建一个新的tailTask任务
		tt := newTailTask(conf.Path, conf.Topic) //创建一个日志收集项
		err = tt.Init()                          //去打开文件准备读
		if err != nil {
			logrus.Errorf("create a tail task for path:%s failed err:%v", conf.Path, err)
			continue
		}
		logrus.Infof("create a tail task for path:%s success", conf.Path)
		ttMgr.tailTaskMap[tt.path] = tt
		//起一个后台的goroutine去收集日志
		go tt.run()
	}
	go ttMgr.watch() //在后台等新的配置来
	return
}

func (t *tailTaskMgr) watch() {
	for {
		//派一个小弟等待新配置
		newConf := <-t.confChan //取到值说明新的配置来了
		//新配置来了后应该管理一下我之前启动的那些tailTask
		logrus.Infof("get new conf from etcd，conf:%v", newConf)
		for _, conf := range newConf {
			//	1. 原来已经存在的任务就不用动
			if t.isExist(conf) {
				continue
			}
			//	2. 原来没有的我新创建一个tailTask任务
			tt := newTailTask(conf.Path, conf.Topic) //创建一个日志收集项
			err := tt.Init()                         //去打开文件准备读
			if err != nil {
				logrus.Errorf("create a tail task for path:%s failed err:%v", conf.Path, err)
				continue
			}
			logrus.Infof("create a tail task for path:%s success", conf.Path)
			t.tailTaskMap[tt.path] = tt
			//起一个后台的goroutine去收集日志
			go tt.run()
		}
		//	3. 原来有的现在没有
		for key, task := range t.tailTaskMap {
			found := false
			for _, conf := range newConf {
				if key == conf.Path {
					found = true
					break
				}
			}
			if !found {
				//这个tailTask要停了
				logrus.Infof("need to stop tailTask path:%s", task.path)
				delete(t.tailTaskMap, key)
				task.cancel()
			}
		}
	}
}

func (t *tailTaskMgr) isExist(conf common.CollectEntry) bool {
	_, ok := t.tailTaskMap[conf.Path]
	return ok
}
func SendNewConf(newConf []common.CollectEntry) {
	ttMgr.confChan <- newConf
}
