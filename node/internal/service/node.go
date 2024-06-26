package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/Ferriem/Distributed_schedule_platform/common/models"
	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/config"
	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/dbclient"
	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/etcdclient"
	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/logger"
	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/utils"
	"github.com/Ferriem/Distributed_schedule_platform/node/internal/handler"
	"github.com/jakecoffman/cron"
	"github.com/ouqiang/goutil"
)

type NodeServer struct {
	*etcdclient.ServerReg
	*models.Node
	*cron.Cron

	jobs handler.Jobs
}

func NewNodeServer() (*NodeServer, error) {
	uuid, err := utils.UUID()
	if err != nil {
		return nil, err
	}
	ip, err := utils.LocalIP()
	if err != nil {
		return nil, err
	}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = uuid
		err = nil
	}
	return &NodeServer{
		Node: &models.Node{
			UUID:     uuid,
			PID:      strconv.Itoa(os.Getpid()),
			IP:       ip.String(),
			Hostname: hostname,
			UpTime:   time.Now().Unix(),
			Status:   models.NodeConnSuccess,
			Version:  config.GetConfigModels().System.Version,
		},
		Cron:      cron.New(),
		jobs:      make(handler.Jobs, 8),
		ServerReg: etcdclient.NewServerReg(config.GetConfigModels().System.NodeTtl),
	}, nil
}

// Check whether the node is registered with ETCD
// If yes, PID is returned. If no, -1 is returned
func (srv *NodeServer) exist(nodeUUID string) (pid int, err error) {
	resp, err := etcdclient.Get(fmt.Sprintf(etcdclient.KeyEtcdNode, nodeUUID))
	if err != nil {
		return
	}
	if len(resp.Kvs) == 0 {
		return -1, nil
	}

	if pid, err = strconv.Atoi(string(resp.Kvs[0].Value)); err != nil {
		if _, err = etcdclient.Delete(fmt.Sprintf(etcdclient.KeyEtcdNode, nodeUUID)); err != nil {
			return
		}
		return -1, nil
	}

	p, err := os.FindProcess(pid)
	if err != nil {
		return -1, nil
	}

	if p != nil && p.Signal(syscall.Signal(0)) == nil {
		return
	}
	return -1, nil
}

// Register into ETCD with /crony/node/<node_id>
func (srv *NodeServer) Register() error {
	pid, err := srv.exist(srv.UUID)
	if err != nil {
		return err
	}
	if pid != -1 {
		return fmt.Errorf("node[%s] with pid[%d] exist", srv.UUID, pid)
	}
	b, err := json.Marshal(&srv.Node)
	if err != nil {
		return fmt.Errorf("node[%s] with pid[%d] json error%s", srv.UUID, pid, err.Error())
	}
	if err := srv.ServerReg.Register(fmt.Sprintf(etcdclient.KeyEtcdNode, srv.UUID), string(b)); err != nil {
		return err
	}
	return nil
}

func (srv *NodeServer) Stop(i interface{}) {
	srv.Down()
	_, err := etcdclient.Delete(fmt.Sprintf(etcdclient.KeyEtcdNode, srv.UUID))
	if err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("failed to delete etcd node[%s] key error%s", srv.UUID, err.Error()))
	}
	_, err = etcdclient.Delete(fmt.Sprintf(etcdclient.KeyEtcdSystemGet, srv.UUID))
	if err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("failed to delete system etcd node[%s] key error%s", srv.UUID, err.Error()))
	}

	_ = srv.Client.Close()
	srv.Cron.Stop()
}

func (srv *NodeServer) Down() {
	srv.Status = models.NodeConnFail
	srv.DownTime = time.Now().Unix()
	err := srv.Node.Update()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("failed to update node[%s] down error:%s", srv.UUID, err.Error()))
	}
	err = dbclient.GetMysqlDB().Table(models.CronyJobTableName).Select("status").Where("run_on = ?", srv.UUID).Updates(models.Job{
		Status: models.JobStatusNotAssigned,
	}).Error
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("failed to update job on node[%s] down error%s", srv.UUID, err.Error()))
	}
}

func (srv *NodeServer) Run() (err error) {
	defer func() {
		if err != nil {
			srv.Stop(err)
		}
	}()

	if err = srv.loadJobs(); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("node[%s] failed to load job error%s", srv.UUID, err.Error()))
		return
	}

	insertId, err := srv.Node.Insert()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("failed to create ndoe[%s] into db error%s", srv.UUID, err.Error()))
		return
	}
	srv.Node.ID = insertId
	srv.Cron.Start()
	go srv.watchJobs()
	go srv.watchOnce()
	go srv.watchSystemInfo()
	return
}

func (srv *NodeServer) loadJobs() (err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.GetLogger().Error(fmt.Sprintf("load jobs panic:%v", r))
		}
	}()
	jobs, err := handler.GetJobs(srv.UUID)
	if err != nil {
		return
	}

	if len(jobs) == 0 {
		return
	}
	srv.jobs = jobs

	for _, j := range jobs {
		j.InitNodeInfo(models.JobStatusAssigned, srv.UUID, srv.Hostname, srv.IP)
		srv.addJob(j)
	}
	return
}

func (srv *NodeServer) addJob(j *handler.Job) {
	if err := j.Check(); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("job[%d] check error :%s", j.ID, err.Error()))
		return
	}
	//run script
	if j.Type == models.JobTypeCmd {
		for _, id := range j.ScriptIDArray {
			script := &models.Script{ID: id}
			err := script.FindByID()
			if err != nil {
				logger.GetLogger().Warn(fmt.Sprintf("job[%d] find script[%d] error:%s", j.ID, id, err.Error()))
				continue
			}
			err = script.Check()
			if err != nil {
				logger.GetLogger().Warn(fmt.Sprintf("job[%d] check script[%d] error:%s", j.ID, id, err.Error()))
				continue
			}
			result, err := handler.RunPresetScript(script)
			if err != nil {
				logger.GetLogger().Warn(fmt.Sprintf("job[%d] run script[%d] error:%s", j.ID, id, err.Error()))
			}
			logger.GetLogger().Info(fmt.Sprintf("job[%d] run script[%d] result:%s", j.ID, id, result))
		}
	}
	taskFunc := handler.CreateJob(j)
	if taskFunc == nil {
		logger.GetLogger().Warn(fmt.Sprintf("Failed to create a task to process the Job. The task protocol was not supported%s", j.Type))
		return
	}
	err := goutil.PanicToError(func() {
		srv.Cron.AddFunc(j.Spec, taskFunc, srv.jobCronName(j.ID))
	})
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("Failed to add a task to the scheduler#%v", err.Error()))
	}
	return
}

func (srv *NodeServer) jobCronName(jobId int) string {
	return fmt.Sprintf(srv.UUID+"/%d", jobId)
}

func (srv *NodeServer) modifyJob(j *handler.Job) {
	oldJob, ok := srv.jobs[j.ID]
	if !ok {
		srv.addJob(j)
		return
	}
	srv.deleteJob(oldJob.ID)
	srv.addJob(j)
	return
}

func (srv *NodeServer) deleteJob(jobId int) {
	if _, ok := srv.jobs[jobId]; !ok {
		srv.Cron.RemoveJob(srv.jobCronName(jobId))
		delete(srv.jobs, jobId)
		return
	}
	return
}
