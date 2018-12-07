package rediscli

import (
	"../../config"
	"../../logs"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
	"sync"
	"time"
)

var (
	once     sync.Once
	RedisCli *redis.Client
)

func RedisConn() {
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     config.REDIS_HOST,
		Password: "",              // no password set
		DB:       config.REDIS_DB, // use default DB
	})
	pong, err := RedisCli.Ping().Result()
	logs.Log.Debug("pong:%v, err: %v", pong, err)
}

func init() {
	once.Do(func() {
		RedisConn()
	})
}

func Info() (info map[string]interface{}) {
	info = make(map[string]interface{})
	infoString := RedisCli.Info().String()
	for _, v := range strings.Split(infoString, "\r\n\r\n#") {
		item := strings.Split(v, "\r\n")
		tmp := make(map[string]interface{})
		for _, k := range item[1:] {
			k_v := strings.Split(k, ":")
			if len(k_v) > 1 {
				tmp[k_v[0]] = k_v[1]
			}
		}
		info[item[0]] = tmp
	}
	return
}

func Len(args ...string) (lm map[string]int64) {
	lm = make(map[string]int64)
	for _, k := range args {
		lm[k] = RedisCli.LLen(k).Val()
	}
	return
}

func ListKey(prefix string) (keys []string) {
	var cursor uint64
	var n int
	for {
		var err error
		var key []string
		prefix = fmt.Sprintf("%v*", prefix)
		key, cursor, err = RedisCli.Scan(cursor, prefix, 10).Result()
		if err != nil {
			panic(err)
		}
		keys = append(keys, key...)
		n += len(key)
		if cursor == 0 {
			break
		}
	}
	return keys
}

func LPush(key string, values ...interface{}) int64 {
	intCmd := RedisCli.LPush(key, values...)
	return intCmd.Val()
}

func BRPop(timeout time.Duration, keys ...string) (result []string) {
	stringSliceCmd := RedisCli.BRPop(timeout, keys...)
	result = stringSliceCmd.Val()
	return
}
