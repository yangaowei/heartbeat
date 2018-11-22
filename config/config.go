package config

// mysqlconnstring       string = "root:@tcp(127.0.0.1:3306)" // mysql连接字符串
//     mysqlconncap          int    = 2048                        // mysql连接池容量
//     mysqlmaxallowedpacket int    = 1048576                     //mysql通信缓冲区的最大长度，单位B，默认1MB
var (
	MYSQL_CONN_CAP           int    = 2048
	MYSQL_CONN_STR           string = "root:123456@tcp(localhost:3306)/heartbeat" // mysql连接字符串
	MYSQL_MAX_ALLOWED_PACKET int    = 1048576                                     // mysql通信缓冲区的最大长度
)
