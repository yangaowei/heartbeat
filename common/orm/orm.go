package orm

import (
	"../../config"
	"../../logs"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var (
	Tables map[string]string  = make(map[string]string)
	Conn   map[string]*sql.DB = make(map[string]*sql.DB)
	once   sync.Once
	db     *sql.DB
	err    error
	mu     sync.Mutex
)

type orm struct {
}

var Orm = new(orm)

func RegisterDataBase(dbname string) {
	once.Do(func() {
		db, err = sql.Open("mysql", config.MYSQL_CONN_STR+"?charset=utf8")
		if err != nil {
			logs.Log.Error("Mysql：%v\n", err)
			return
		}
		db.SetMaxOpenConns(config.MYSQL_CONN_CAP)
		db.SetMaxIdleConns(config.MYSQL_CONN_CAP)
	})
	if err = db.Ping(); err != nil {
		//logs.Log.Error("Mysql：%v\n", err)
		logs.Log.Debug("url: %v", err)
	}
	//logs.Log.Debug("db: %v", db)
	Conn[dbname] = db
}

func RegisterTable(table string) error {
	mu.Lock()
	if _, ok := Tables[table]; !ok {
		Tables[table] = table
	}
	mu.Unlock()
	return nil
}

func getTableName(val reflect.Value) string {
	if fun := val.MethodByName("TableName"); fun.IsValid() {
		vals := fun.Call([]reflect.Value{})
		// has return and the first val is string
		if len(vals) > 0 && vals[0].Kind() == reflect.String {
			return vals[0].String()
		}
	}
	return reflect.Indirect(val).Type().Name()
}

func genkv(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	getType := reflect.TypeOf(obj).Elem()
	getValue := reflect.ValueOf(obj).Elem()
	fmt.Println(getType, getValue)
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.Field(i)
		var v interface{}
		switch value.Kind() {
		case reflect.String:
			if len(value.String()) > 0 {
				v = value.String()
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value.Int() > 0 {
				v = value.Int()
			}
		default:
		}
		if v != nil {
			result[strings.ToLower(field.Name[:1])+field.Name[1:]] = v
		}
	}
	return result
}

func genSql(data map[string]interface{}, table, action string) (sql string) {
	switch action {
	case "select":
		sql = fmt.Sprintf("select * from %s", table)
		var where string
		for k, v := range data {
			where += fmt.Sprintf(" %v='%v' ", k, v)
		}
		if len(where) > 0 {
			sql = fmt.Sprintf("%s where %s", sql, where)
		}
	}

	return sql
}

func (o *orm) List(obj interface{}) (objs []interface{}) {
	getValue := reflect.ValueOf(obj)
	getType := reflect.TypeOf(obj).Elem()
	table := getTableName(getValue)
	kv := genkv(obj)
	sql := genSql(kv, table, "select")
	logs.Log.Debug(" *sql: %v", sql)
	rows, _ := db.Query(sql)
	cols, _ := rows.Columns()
	buff := make([]interface{}, len(cols)) // 临时slice
	data := make([]string, len(cols))      // 存数据slice
	for i, _ := range buff {
		buff[i] = &data[i]
	}
	tmp := reflect.New(getType).Interface()
	// logs.Log.Debug(" *tmp: %v", tmp)
	// logs.Log.Debug(" *tmp: %v", obj)
	for rows.Next() {
		rows.Scan(buff...) // ...是必须的
		genObj(data, cols, tmp)
		objs = append(objs, tmp)
	}
	return
}

func genObj(data, cols []string, obj interface{}) {
	getType := reflect.TypeOf(obj).Elem()
	getValue := reflect.ValueOf(obj).Elem()
	for i := 0; i < getValue.NumField(); i++ {
		//fmt.Println(getType.Field(i))
		//getType.Field(i)
		v := getValue.Field(i)
		//SetDefault(getType.Field(i), v)
		for index, col := range cols {
			Fname := getType.Field(i).Name
			if strings.ToUpper(col[:1])+col[1:] == Fname {
				//v := getValue.Field(i)
				//fmt.Println(getType.Field(i).Tag)
				parseQueryColumn(v, data[index])
			}
		}
	}
}

//单个属性赋值
func parseQueryColumn(field reflect.Value, s string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(s)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, _ := strconv.ParseUint(s, 10, 0)
		field.SetUint(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, _ := strconv.ParseInt(s, 10, 0)
		field.SetInt(v)
	case reflect.Float32:
		v, _ := strconv.ParseFloat(s, 32)
		field.SetFloat(v)
	case reflect.Float64:
		v, _ := strconv.ParseFloat(s, 64)
		field.SetFloat(v)
	default:

	}
}

func (o *orm) Get(obj interface{}) interface{} {
	objs := o.List(obj)
	if len(objs) > 0 {
		return objs[0]
	} else {
		return nil
	}
}

func (o *orm) DbInfo() map[string]*sql.DB {
	return Conn
}
