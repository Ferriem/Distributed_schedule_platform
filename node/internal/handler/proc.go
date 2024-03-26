package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/Ferriem/Distributed_schedule_platform/common/models"
	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/config"
	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/etcdclient"
	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/logger"
	"github.com/coreos/etcd/clientv3"
)

// Information about the current task in execution
// key: /crony/proc/<node_uuid>/<job_id>/pid</job_id></node_uuid>
// value: indicates the start execution time
// The key expires automatically to prevent the key from being cleared after the process exits unexpectedly. The expiration time can be configured
type JobProc struct {
	*models.JobProc
}

func GetProcFromKey(key string) (proc *JobProc, err error) {
	ss := strings.Split(key, "/")
	var sslen = len(ss)
	if sslen < 5 {
		err = fmt.Errorf("invalid proc key [%s]", key)
	}
	id, err := strconv.Atoi(ss[sslen-1])
	if err != nil {
		return
	}
	jobId, err := strconv.Atoi(ss[sslen-2])
	if err != nil {
		return
	}
	proc = &JobProc{
		JobProc: &models.JobProc{
			ID:       id,
			JobID:    jobId,
			NodeUUID: ss[sslen-3],
		},
	}
	return
}

func (p *JobProc) Key() string {
	return fmt.Sprintf(etcdclient.KeyEtcdProc, p.NodeUUID, p.JobID, p.ID)
}

func (p *JobProc) del() error {
	_, err := etcdclient.Delete(p.Key())
	return err
}

func (p *JobProc) Stop() {
	if p == nil {
		return
	}
	if !atomic.CompareAndSwapInt32(&p.Running, 1, 0) {
		return
	}
	p.Wg.Wait()

	if err := p.del(); err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("proc del[%s] err: %s", p.Key(), err.Error()))
	}
}

func WatchProc(nodeUUID string) clientv3.WatchChan {
	return etcdclient.Watch(fmt.Sprintf(etcdclient.KeyEtcdProc, nodeUUID), clientv3.WithPrefix())
}

func (p *JobProc) Start() error {
	if !atomic.CompareAndSwapInt32(&p.Running, 0, 1) {
		return nil
	}

	p.Wg.Add(1)
	b, err := json.Marshal(p.JobProcVal)
	if err != nil {
		return err
	}
	_, err = etcdclient.PutWithTtl(p.Key(), string(b), config.GetConfigModels().System.JobProcTtl)
	if err != nil {
		return err
	}
	p.Wg.Done()
	return nil
}
