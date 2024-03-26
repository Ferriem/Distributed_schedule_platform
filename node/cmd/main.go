package main

import (
	"fmt"
	"os"

	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/dbclient"
	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/logger"
	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/notify"
	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/server"
	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/utils/event"
	"github.com/Ferriem/Distributed_schedule_platform/node/internal/service"
)

const ServerName = "node"

func main() {
	if _, err := server.InitNodeServer(ServerName); err != nil {
		fmt.Println("init node server error:", err.Error())
		os.Exit(1)
	}
	nodeServer, err := service.NewNodeServer()
	if err != nil {
		fmt.Println("init node server error:", err.Error())
		os.Exit(1)
	}
	service.RegisterTables(dbclient.GetMysqlDB())
	if err = nodeServer.Register(); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("register node into etcd error:%s", err.Error()))
		os.Exit(1)
	}
	if err = nodeServer.Run(); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("node run error: %s", err.Error()))
		os.Exit(1)
	}
	go notify.Serve()
	logger.GetLogger().Info(fmt.Sprintf("crony node %s service started, Ctrl+C or send kill sign to exit", nodeServer.String()))
	event.OnEvent(event.EXIT, nodeServer.Stop)
	event.WaitEvent()
	event.EmitEvent(event.EXIT, nil)
	logger.GetLogger().Info("exit success")
}
