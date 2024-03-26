package notify

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/httpclient"
	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/logger"
)

type WebHook struct {
	Kind string
	Url  string
}

var _defaultWebHook *WebHook

func (w *WebHook) SendMsg(msg *Message) {
	switch _defaultWebHook.Kind {
	case "feishu":
		var sendData = feishuTemplateCard
		sendData = strings.Replace(sendData, "timeSlot", msg.OccurTime, 1)
		sendData = strings.Replace(sendData, "ipSlot", msg.IP, 1)

		userSlot := ""
		for _, to := range msg.To {
			userSlot += fmt.Sprintf("<at email='' >%s</at>", to)
		}
		sendData = strings.Replace(sendData, "userSlot", userSlot, 1)
		sendData = strings.Replace(sendData, "msgSlot", msg.Body, 1)
		sendData = strings.Replace(sendData, "subjectSlot", msg.Subject, 1)
		_, err := httpclient.PostJson(_defaultWebHook.Url, sendData, 0)
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("send feishu msg[%+v] err: %s", msg, err.Error()))
		}
	default:
		b, err := json.Marshal(msg)
		if err != nil {
			return
		}
		_, err = httpclient.PostJson(_defaultWebHook.Url, string(b), 0)
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("send webhook msg[%+v] err: %s", msg, err.Error()))
		}
	}
}
