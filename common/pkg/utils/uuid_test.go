package utils_test

import (
	"testing"

	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/utils"
)

func TestUUID(t *testing.T) {
	uuid, err := utils.UUID()
	if err != nil {
		t.Error(err)
	}
	t.Log(uuid)
}
