package handler_test

import (
	"testing"

	"github.com/Ferriem/Distributed_schedule_platform/common/models"
	"github.com/Ferriem/Distributed_schedule_platform/node/internal/handler"
)

func TestCmd(t *testing.T) {
	jobs := []handler.Job{
		{
			Job: &models.Job{
				Name:    "test",
				Command: "/home/Ferriem/hello.sh",
				Cmd:     []string{"/home/Ferriem/hello.sh"},
				Type:    models.JobTypeCmd,
				Timeout: 0,
			},
		},
		//{
		//	Name:          "post",
		//	HttpUrl:       "",
		//	HttpMethod:    models.HTTPMethodPost,
		//	Timeout:       3000,
		//},
	}
	var cmd handler.CMDHandler
	for i := 0; i < len(jobs); i++ {
		rsp, err := cmd.Run(&jobs[i])
		if err != nil {
			t.Error(err)
		}
		t.Log(rsp)
	}
}
