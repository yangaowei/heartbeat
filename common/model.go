package common

import (
	"../config"
	"encoding/json"
	//"fmt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//    expect_status = request.headers.get('Expect-Status')
//    message_id = request.headers.get('Message-Id', uuid.uuid1().hex)
//    expect_content = request.headers.get('Expect-Content')
//    send_timestamp = int(request.headers.get('Send-Timestamp', time.time()))
//    send_timeout = int(request.headers.get('Send-Timeout', 60))
//    overdue_timestamp = int(request.headers.get('Overdue-Timestamp', send_timestamp + 24 * 3600))

type Message struct {
	Project          string `json:"prject"`
	Callback         string
	Body             []byte
	ExpectStatus     []int
	ContentType      string
	MeessgeId        string
	ExpectContent    string
	SendTimestamp    int64
	SendTimeout      int64
	OverdueTimestamp int64
	SendNum          int
	LastSendTime     int64
}

func NewMessage(header *http.Header, body []byte, project, callback string) (message *Message) {
	message = &Message{Project: project, Callback: callback, Body: body, MeessgeId: uuid.New().String()}
	expectStatus := header.Get("Expect-Status")
	if expectStatus == "" {
		message.ExpectStatus = config.Expect_Status
	} else {
		tmp := []int{}
		for _, item := range strings.Split(expectStatus, ",") {
			status, err := strconv.ParseInt(item, 10, 0)
			if err == nil {
				tmp = append(tmp, int(status))
			}
		}
		if len(tmp) > 0 {
			message.ExpectStatus = tmp
		} else {
			message.ExpectStatus = config.Expect_Status
		}
	}
	contentType := header.Get("Content-Type")
	if contentType != "" {
		message.ContentType = contentType
	}
	expectContent := header.Get("Expect-Content")
	message.ExpectContent = expectContent
	sendTimeout := header.Get("Send-Timeout")
	timeout, tRrr := strconv.ParseInt(sendTimeout, 10, 64)
	if tRrr == nil {
		message.SendTimeout = timeout
	} else {
		message.SendTimeout = config.Send_Timeout
	}
	message.SendTimestamp = time.Now().Unix()
	message.OverdueTimestamp = message.SendTimestamp + config.Overdue_Timestamp
	return
}

func (self *Message) ToJson() []byte {
	jsonBytes, err := json.Marshal(self)
	if err == nil {
		return jsonBytes
	} else {
		return nil
	}
}

func (self *Message) IsExceptStatus(status int) bool {
	for _, v := range self.ExpectStatus {
		if v == status {
			return true
		}
	}
	return false
}
