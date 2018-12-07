package common

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNewMessage(t *testing.T) {
	header := make(http.Header)
	var body []byte
	project := "test"
	callback := "http://www.baidu.com"
	// header.Add("Expect-Status", "301,302")
	// header.Add("Expect-Content", "1")
	// header.Add("Send-Timeout", "20")
	message := NewMessage(&header, body, project, callback)
	fmt.Println(message)

	fmt.Println(string(message.ToJson()))
}
