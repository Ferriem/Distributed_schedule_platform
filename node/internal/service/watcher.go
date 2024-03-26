package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Ferriem/Distributed_schedule_platform/common/models"
	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/etcdclient"
	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/logger"
	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/utils"
	"github.com/Ferriem/Distributed_schedule_platform/node/internal/handler"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

func (srv *NodeServer) watchJobs() {
	rch := handler.WatchJobs(srv.UUID)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch {
			case ev.IsCreate():
				var job handler.Job
				if err := json.Unmarshal(ev.Kv.Value, &job); err != nil {
					err = fmt.Errorf("watch job[%s] create json unmarshal err: %s", string(ev.Kv.Key), err.Error())
					continue
				}
				srv.jobs[job.ID] = &job
				job.InitNodeInfo(models.JobStatusAssigned, srv.UUID, srv.Hostname, srv.IP)
				srv.addJob(&job)
			case ev.IsModify():
				var job handler.Job
				if err := json.Unmarshal(ev.Kv.Value, &job); err != nil {
					err = fmt.Errorf("watch job[%s] modify json unmarshal err: %s", string(ev.Kv.Key), err.Error())
					continue
				}
				job.InitNodeInfo(models.JobStatusAssigned, srv.UUID, srv.Hostname, srv.IP)
				srv.modifyJob(&job)
			case ev.Type == mvccpb.DELETE:
				srv.deleteJob(handler.GetJobIDFromKey(string(ev.Kv.Key)))
			default:
				logger.GetLogger().Warn(fmt.Sprintf("watch job unknown event type[%v] from job[%s]", ev.Type, string(ev.Kv.Key)))
			}
		}
	}
}

func (srv *NodeServer) watchSystemInfo() {
	rch := handler.WatchSystem(srv.UUID)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch {
			case ev.IsCreate() || ev.IsModify():
				key := string(ev.Kv.Key)
				if string(ev.Kv.Value) != models.NodeSystemInfoSwitch || srv.Node.UUID != getUUID(key) {
					logger.GetLogger().Error(fmt.Sprintf("get system info from node[%s] ,switch is not alive", srv.Node.UUID))
					continue
				}
				s, err := utils.GetServerInfo()
				if err != nil {
					logger.GetLogger().Error(fmt.Sprintf("get system info from node[%s] error:%s", srv.Node.UUID, err.Error()))
					continue
				}
				b, err := json.Marshal(s)
				if err != nil {
					logger.GetLogger().Error(fmt.Sprintf("get system info from node[%s] json marshal error:%s", srv.Node.UUID, err.Error()))
					continue
				}
				_, err = etcdclient.PutWithTtl(fmt.Sprintf(etcdclient.KeyEtcdSystemGet, getUUID(key)), string(b), 5*60)
				if err != nil {
					logger.GetLogger().Error(fmt.Sprintf("get system info from node[%s] etcd put error:%s", srv.UUID, err.Error()))
					continue
				}
			}
		}
	}
}

func getUUID(key string) string {
	index := strings.LastIndex(key, "/")
	if index == -1 {
		return ""
	}
	return key[index+1:]
}

func (srv *NodeServer) watchOnce() {
	rch := handler.WatchOnce()
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch {
			case ev.IsCreate(), ev.IsModify():
				// do not execute in this node
				if len(ev.Kv.Value) != 0 && string(ev.Kv.Value) != srv.UUID {
					continue
				}
				j, ok := srv.jobs[handler.GetJobIDFromKey(string(ev.Kv.Key))]
				if !ok {
					continue
				}
				go j.RunWithRecovery()
			}
		}
	}
}
