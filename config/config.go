package config

var (
	MYSQL_CONN_CAP           int    = 2048
	MYSQL_CONN_STR           string = "root:123456@tcp(localhost:3306)/heartbeat" // mysql连接字符串
	MYSQL_MAX_ALLOWED_PACKET int    = 1048576                                     // mysql通信缓冲区的最大长度

	REDIS_HOST        string = "localhost:6379"
	REDIS_PASSWORD    string = ""
	REDIS_DB          int    = 8
	API_PROT          int    = 8081
	Expect_Status     []int  = []int{200, 201, 202, 203, 204, 205, 206}
	Send_Timeout      int64  = 60        //回调超时时间
	Overdue_Timestamp int64  = 3600 * 24 //失效时长
)
