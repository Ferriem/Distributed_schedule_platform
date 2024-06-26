package notify

import (
	"strings"
	"time"

	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/utils"
)

type Noticer interface {
	SendMsg(*Message)
}

type Message struct {
	Type      int
	IP        string
	Subject   string
	Body      string
	To        []string
	OccurTime string
}

var msgQueue chan *Message

func Init(mail *Mail, web *WebHook) {
	_defaultMail = &Mail{
		Port:     mail.Port,
		From:     mail.From,
		Host:     mail.Host,
		Secret:   mail.Secret,
		Nickname: mail.Nickname,
	}
	_defaultWebHook = &WebHook{
		Kind: web.Kind,
		Url:  web.Url,
	}
	msgQueue = make(chan *Message, 64)
}

func Send(msg *Message) {
	msgQueue <- msg
}

func Serve() {
	for {
		select {
		case msg := <-msgQueue:
			if msg == nil {

			}
			switch msg.Type {
			case 1:
				//Mail
				msg.Check()
				_defaultMail.SendMsg(msg)
			case 2:
				//webhook
				msg.Check()
				go _defaultMail.SendMsg(msg)
			}
		}
	}
}

func (m *Message) Check() {
	if m.OccurTime == "" {
		m.OccurTime = time.Now().Format(utils.TimeFormatSecond)
	}
	//Remove th transfer character
	m.Body = strings.Replace(m.Body, "\"", "'", -1)
	m.Body = strings.Replace(m.Body, "\n", "", -1)
}
