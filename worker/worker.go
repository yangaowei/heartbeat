package worker

import (
	"../common"
	"../common/rediscli"
	//"../config"
	"../logs"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	tr     *http.Transport
	client *http.Client
)

func init() {
	tr = &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 100,
	}

	client = &http.Client{
		Transport: tr,
	}
}

func callBackMessage(messge *common.Message) (flag bool) {
	now := time.Now().Unix()
	if messge.OverdueTimestamp < now {
		return true
	}
	if (messge.LastSendTime + 6) > now {
		logs.Log.Debug("LastSendTime %v,sleep", messge.LastSendTime)
		time.Sleep(3 * time.Second)
		return false
	}
	req, err := http.NewRequest("POST", messge.Callback, bytes.NewBuffer(messge.Body))
	if err == nil {
		resp, e := client.Do(req)
		if e == nil {
			status := resp.StatusCode
			if messge.IsExceptStatus(status) {
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				logs.Log.Debug("body %v", string(body))
				if messge.ExpectContent == string(body) {
					flag = true
				}
			} else {
				logs.Log.Debug("ExceptStatus error, ExceptStatus:%v, status:%v", messge.ExpectStatus, status)
			}
		}
	}
	if !flag {
		messge.SendNum += 1
		messge.LastSendTime = time.Now().Unix()
	}
	return
}

func Run() {
	for {
		keys := rediscli.ListKey("project")
		value := rediscli.BRPop(1*time.Second, keys...)
		if value == nil {
			logs.Log.Debug("no task , sleep 1s")
			time.Sleep(1 * time.Second)
		} else {
			message := new(common.Message)
			err := json.Unmarshal([]byte(value[1]), message)
			if err == nil {
				if !callBackMessage(message) {
					logs.Log.Debug("callBackMessage err, %v", string(message.ToJson()))
					rediscli.LPush(fmt.Sprintf("project:%s", message.Project), message.ToJson())
				}
			} else {
				logs.Log.Debug("Unmarshal errr , %v", err)
				continue
			}
		}
	}
}
