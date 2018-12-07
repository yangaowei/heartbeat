package rediscli

import (
	"fmt"
	"testing"
	"time"
)

func TestRe(t *testing.T) {
	fmt.Println(RedisCli)
}

func TestLPush(t *testing.T) {
	intCmd := RedisCli.LPush("project:test", []interface{}{"test1", "test2"}...)
	fmt.Println(intCmd.Result())
	intCmd = RedisCli.LPush("project:2", "10")
	fmt.Println(intCmd.Result())
}

//BRPop(timeout time.Duration, keys ...string) *StringSliceCmd
func TestBRPop(t *testing.T) {
	// stringSliceCmd := RedisCli.BRPop(1*time.Second, []string{"project:test", "project:2"}...)
	// fmt.Println(stringSliceCmd.Result())
	// for _, v := range stringSliceCmd.Val() {
	// 	fmt.Println("v", v)
	// }

	fmt.Println(BRPop(1*time.Second, "project:test"), "BRPop result")
	// var cursor uint64
	// keys, cursor, err := RedisCli.Scan(cursor, "project:*", 1).Result()
	// fmt.Println(keys, cursor, err)

	// llen := RedisCli.LLen("project:test")
	// fmt.Println(llen.Result())
	// fmt.Println(stringSliceCmd.Val())
	// fmt.Println(stringSliceCmd.Name())
	// fmt.Println(stringSliceCmd.String())
}

func TestInfo(t *testing.T) {
	info := Info()
	fmt.Println(info)
}

func TestLen(t *testing.T) {
	info := Len([]string{"project:test", "project:2"}...)
	fmt.Println(info)
}
func TestListKey(t *testing.T) {
	info := ListKey("project")
	fmt.Println(info)
}
