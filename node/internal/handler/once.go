package handler

import (
	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/etcdclient"
	"github.com/coreos/etcd/clientv3"
)

func WatchOnce() clientv3.WatchChan {
	return etcdclient.Watch(etcdclient.KeyEtcdJobProfile, clientv3.WithPrefix())
}
