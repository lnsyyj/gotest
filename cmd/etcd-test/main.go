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
	ctx := context.Background()

	// put
	putResp, err := kv.Put(ctx, "/samba/global/k1", "value")
	if err != nil {
		fmt.Errorf("put to etcd failed, err:%v, ", err)
	}
	fmt.Println("Revision:", putResp.Header.Revision)

	putResp, err = kv.Put(ctx, "/samba/global/k2", "value")
	if err != nil {
		fmt.Errorf("put to etcd failed, err:%v, ", err)
	}
	fmt.Println("Revision:", putResp.Header.Revision)

	// get
	getResp, err := kv.Get(ctx, "/samba/global/k1")
	if err != nil {
		fmt.Errorf("get etcd failed, err:%v, ", err)
	}
	fmt.Println("Revision:", getResp.Header.Revision)
	fmt.Println("Kvs:", getResp.Kvs)

	// get WithPrefix
	getResp, err = kv.Get(ctx, "/samba/global/", clientv3.WithPrefix())
	if err != nil {
		fmt.Errorf("get etcd failed, err:%v, ", err)
	}
	fmt.Println("Revision:", getResp.Header.Revision)
	fmt.Println("Kvs:", getResp.Kvs)

	// lease
	lease := clientv3.NewLease(cli)
	grantResp, err := lease.Grant(ctx, 10)
	fmt.Println("grantRestp:", grantResp.ID)
	if err != nil {
		fmt.Errorf("grant lease failed, err:%v, ", err)
	}
	kv.Put(ctx, "/samba/global/k3", "value", clientv3.WithLease(grantResp.ID))
	keepResp, err := lease.KeepAliveOnce(ctx, grantResp.ID)
	if err != nil {
		fmt.Errorf("keep alive lease failed, err:%v, ", err)
	}
	fmt.Println("Revision:", keepResp.Revision)

	// OP
	ops := []clientv3.Op{
		clientv3.OpPut("/samba/global/k4", "value"),
		clientv3.OpGet("/samba/global/k4")}
	for _, op := range ops {
		_, err := cli.Do(ctx, op)
		if err != nil {
			fmt.Errorf("do etcd failed, err:%v, ", err)
		}
	}

	// Txn 事务
	txn := kv.Txn(ctx)
	txn.If(
		clientv3.Compare(clientv3.CreateRevision("/samba/global/k5"), "=", 0),
	).Then(
		clientv3.OpPut("/samba/global/k5", "value"),
	).Else(
		clientv3.OpGet("/samba/global/k5"),
	).Commit()

	// watch
	watchChan := cli.Watch(ctx, "/samba/global/k1")
	for wresp := range watchChan {
		for _, ev := range wresp.Events {
			fmt.Printf("Type:%v key:%v value:%v\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}

}
