package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/client/v3"
	"sysagentv2/common"
	"sysagentv2/tailfile"
	"time"
)

var (
	client *clientv3.Client
)

func Init(address []string) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		logrus.Error("client.New err:", err)
		return
	}
	return
}
func GetConf(key string) (collectEntryList []common.CollectEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	grs, err := client.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get conf from etcd key:%s fild,err:%s", key, err)
		return
	}
	if len(grs.Kvs) == 0 {
		logrus.Errorf("len0: get conf from etcd key:%s fild,err:%s", key, err)
		return
	}
	ret := grs.Kvs[0]
	err = json.Unmarshal(ret.Value, &collectEntryList)
	if err != nil {
		logrus.Errorf("json Unmarshal err:%s", err)
		return
	}
	return
}
func WatchConf(key string) {
	for {
		watchChan := client.Watch(context.Background(), key)
		for wresp := range watchChan {
			logrus.Info("get new conf from etcd!")
			for _, evt := range wresp.Events {
				fmt.Printf("type:%s key:%s value:%s\n", evt.Type, evt.Kv.Key, evt.Kv.Value)
				var newConf []common.CollectEntry
				if evt.Type == clientv3.EventTypeDelete {
					//	如果是删除
					logrus.Warning("FBI waring:etcd delete the key!!!")
					tailfile.SendNewConf(newConf) //没有接收就是堵塞
					continue
				}
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					logrus.Errorf("json Unmarshal err:%s", err)
					return
				}
				//告诉tailFile这个模块该启用新的配置了
				tailfile.SendNewConf(newConf) //没有接收就是堵塞
			}
		}
	}
}
