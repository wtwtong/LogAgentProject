package main

import (
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		fmt.Println("client New err:", err)
		return
	}
	watchCh := cli.Watch(context.Background(), "client_log_192.168.3.134_conf")
	for wtch := range watchCh {
		fmt.Println("...")
		for _, vn := range wtch.Events {
			fmt.Printf("Type:%s,key:%s,value:%s\n", vn.Type, vn.Kv.Key, vn.Kv.Value)
		}
	}
}
