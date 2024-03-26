package utils_test

import (
	"testing"

	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/utils"
)

func TestMD5(t *testing.T) {
	t.Log(utils.MD5("123456"))
}
