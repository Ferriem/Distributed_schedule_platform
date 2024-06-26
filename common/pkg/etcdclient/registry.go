package etcdclient

import (
	"context"
	"time"

	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/logger"
	"github.com/coreos/etcd/clientv3"
)

type ServerReg struct {
	Client        *Client
	stop          chan error
	leaseId       clientv3.LeaseID
	cancelFunc    func()
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	Ttl           int64
}

func NewServerReg(ttl int64) *ServerReg {
	return &ServerReg{
		Client: _defaultEtcd,
		Ttl:    ttl,
		stop:   make(chan error),
	}
}

func (s *ServerReg) Register(key, value string) error {
	if err := s.setLease(s.Ttl); err != nil {
		return err
	}
	go s.keepAlive()
	if err := s.putService(key, value); err != nil {
		return err
	}
	return nil
}

func (s *ServerReg) setLease(ttl int64) error {
	leaseResp, err := Grant(ttl)
	if err != nil {
		return nil
	}
	ctx, cancelFunc := context.WithCancel(context.TODO())
	leaseRespChan, err := s.Client.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		return err
	}
	s.leaseId = leaseResp.ID
	s.cancelFunc = cancelFunc
	s.keepAliveChan = leaseRespChan
	return nil
}

func (s *ServerReg) Stop() {
	s.stop <- nil
}

func (s *ServerReg) keepAlive() {
	for {
		select {
		case <-s.stop:
			return
		case leaseKeepResp := <-s.keepAliveChan:
			if leaseKeepResp == nil {
				logger.GetLogger().Info("the lease renewal function has been turned off\n")
				return
			}
		}
	}
}

func (s *ServerReg) putService(key, val string) error {
	kv := clientv3.NewKV(s.Client.Client)
	_, err := kv.Put(context.TODO(), key, val, clientv3.WithLease(s.leaseId))
	return err
}

func (s *ServerReg) RevokeLease() error {
	s.cancelFunc()
	time.Sleep(2 * time.Second)
	_, err := Revoke(s.leaseId)
	return err
}
