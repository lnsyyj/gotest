package main

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://192.168.3.162:2379"},
		DialTimeout: 5 * time.Second,
	})
	defer cli.Close()

	if err != nil {
		// handle error!
		fmt.Errorf("connect to etcd failed, err:%v, ", err)
	}

	kv := clientv3.NewKV(cli)
	putResp, err := kv.Put(context.Background(), "/samba/global/", "value")
	if err != nil {
		fmt.Errorf("put to etcd failed, err:%v, ", err)
	}
	fmt.Println("Revision:", putResp.Header.Revision)

	fmt.Println("Hello, World!")
}
