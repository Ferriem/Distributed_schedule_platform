package handler_test

import (
	"testing"

	"github.com/Ferriem/Distributed_schedule_platform/common/models"
	"github.com/Ferriem/Distributed_schedule_platform/node/internal/handler"
)

func TestHttpCall(t *testing.T) {
	jobs := []handler.Job{
		{
			Job: &models.Job{
				Name:       "get",
				Command:    "http://localhost:8089/ping",
				HttpMethod: models.HTTPMethodGet,
				Timeout:    3000,
			},
		},
		//{
		//	Name:          "post",
		//	HttpUrl:       "",
		//	HttpMethod:    models.HTTPMethodPost,
		//	Timeout:       3000,
		//},
	}
	var http handler.HTTPHandler
	for i := 0; i < len(jobs); i++ {
		rsp, err := http.Run(&jobs[i])
		if err != nil {
			t.Error(err)
		}
		t.Log(rsp)
	}
}
