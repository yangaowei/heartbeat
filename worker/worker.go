package worker

import (
	"../common"
	"../common/rediscli"
	"../config"
	"../logs"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	tr          *http.Transport
	client      *http.Client
	SHOULD_EXIT bool
	wg          sync.WaitGroup
	keys        []string
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

func callBackMessage(messge *common.Message) (flag bool, err error) {
	var req *http.Request
	var resp *http.Response
	var eMsg string
	now := time.Now().Unix()
	if messge.OverdueTimestamp < now {
		logs.Log.Debug("OverdueTimestamp lt now , discarded。 %v", messge.OverdueTimestamp)
		return true, errors.New("OverdueTimestamp lt now , discarded")
	}
	if (messge.LastSendTime + config.RETRY_WAIT_TIME) > now {
		eMsg = fmt.Sprintf("LastSendTime + %v gt now,sleep 1s", messge.LastSendTime, config.RETRY_WAIT_TIME)
		logs.Log.Debug(eMsg)
		time.Sleep(1 * time.Second)
		return false, errors.New(eMsg)
	}
	req, err = http.NewRequest("POST", messge.Callback, bytes.NewBuffer(messge.Body))
	req.Header.Add("Content-Type", messge.ContentType)
	if err == nil {
		resp, err = client.Do(req)
		if err == nil {
			status := resp.StatusCode
			if messge.IsExceptStatus(status) {
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				if messge.ExpectContent == string(body) {
					flag = true
				} else {
					eMsg = fmt.Sprintf("body %v, ExpectContent % v", string(body), messge.ExpectContent)
					err = errors.New(eMsg)
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

func runForever(i int) {
	for {
		reSend := false
		workerName := fmt.Sprintf("worker-%v", i)
		value := rediscli.BRPop(1*time.Second, keys...)
		if value == nil {
			logs.Log.Debug("%v no task , sleep 1s", workerName)
			time.Sleep(1 * time.Second)
		} else {
			message := new(common.Message)
			err := json.Unmarshal([]byte(value[1]), message)
			if err == nil {
				result, e := callBackMessage(message)
				if result {
					logs.Log.Debug("%v callBackMessage success, msg: %v", workerName, e)
				} else {
					logs.Log.Debug("%v callBackMessage err, %v. error : %v", workerName, string(message.ToJson()), e)
					//rediscli.LPush(fmt.Sprintf("project:%s", message.Project), message.ToJson())
					reSend = true
				}
			} else {
				logs.Log.Debug("Unmarshal errr , %v", err)
				continue
			}
			if SHOULD_EXIT || reSend {
				rediscli.LPush(fmt.Sprintf("project:%s", message.Project), message.ToJson())
			}
		}
		if SHOULD_EXIT {
			break
		}
	}
	wg.Done()
}

func Run(wNum int) {
	//fmt.Println(wNum)
	if wNum == 0 {
		wNum = config.Default_Worker_Num
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL)
	//捕获信号
	go func() {
		s := <-c
		logs.Log.Debug("Got signal:, %v", s)
		SHOULD_EXIT = true
	}()
	//更新队列列表
	go func() {
		for {
			keys = rediscli.ListKey("project")
			if len(keys) > 0 {
				time.Sleep(1 * time.Minute)
			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}()
	for i := 1; i < wNum+1; i++ {
		wg.Add(1)
		go runForever(i)
	}
	//time.Sleep(10 * time.Second)
	wg.Wait()
}
