package handler

import (
	"fmt"

	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/etcdclient"
	"github.com/coreos/etcd/clientv3"
)

func WatchSystem(nodeUUID string) clientv3.WatchChan {
	return etcdclient.Watch(fmt.Sprintf(etcdclient.KeyEtcdSystemSwitch, nodeUUID), clientv3.WithPrefix())
}
