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
		fmt.Println("client.New err:", err)
		return
	}
	defer cli.Close()
	//	put
	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	//val := `[{"path":"E:/go_tool/logs/web.log","topic":"web"}]`
	val := `[{"path":"E:/go_tool/logs/sysCpu.log","topic":"sysCpu"},{"path":"E:/go_tool/logs/sysMem.log","topic":"sysMem"},{"path":"E:/go_tool/logs/sysDisk.log","topic":"sysDisk"},{"path":"E:/go_tool/logs/sysNet.log","topic":"sysNet"},{"path":"E:/go_tool/logs/web.log","topic":"web"}]`
	key := "client_log_192.168.3.134_conf"
	_, err = cli.Put(ctx, key, val)
	if err != nil {
		fmt.Println("put filed err:", err)
		return
	}
	cancle()
	//	get
	ctx, cancle = context.WithTimeout(context.Background(), time.Second)
	rp, err := cli.Get(ctx, key)
	if err != nil {
		fmt.Println("get filed err:", err)
		return
	}

	for _, ev := range rp.Kvs {
		fmt.Printf("key:%s,value:%s\n", ev.Key, ev.Value)
	}
	cancle()
}
