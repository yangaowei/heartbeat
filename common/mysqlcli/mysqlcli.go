package mysqlcli

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
	db     *sql.DB
	once   sync.Once
	err    error
	models map[string]reflect.Type
)

type Project struct {
	Id       int `default:"0"`
	Name     string
	Area     string
	Desc     string
	Api_key  string
	Queue    string
	Callback string
	Created  string `default:"0000-00-00 00:00:00"`
	Modified string `default:"0000-00-00 00:00:00"`
}

func DB() (*sql.DB, error) {
	Refresh()
	return db, err
}

func init() {
	DB()
	models = make(map[string]reflect.Type)
	models["project"] = reflect.TypeOf(new(Project)).Elem()
}

func Refresh() {
	once.Do(func() {
		db, err = sql.Open("mysql", config.MYSQL_CONN_STR+"?charset=utf8")
		if err != nil {
			//logs.Log.Error("Mysql：%v\n", err)
			return
		}
		db.SetMaxOpenConns(config.MYSQL_CONN_CAP)
		db.SetMaxIdleConns(config.MYSQL_CONN_CAP)
	})
	if err = db.Ping(); err != nil {
		fmt.Println(err)
		//logs.Log.Error("Mysql：%v\n", err)
		logs.Log.Debug("url: %v", err)
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

func SetDefault(f reflect.StructField, v reflect.Value) {
	default_value := f.Tag.Get("default")
	parseQueryColumn(v, default_value)
}

func orm(data, cols []string, obj interface{}) interface{} {
	getType := reflect.TypeOf(obj).Elem()
	getValue := reflect.ValueOf(obj).Elem()
	for i := 0; i < getValue.NumField(); i++ {
		//fmt.Println(getType.Field(i))
		//getType.Field(i)
		v := getValue.Field(i)
		SetDefault(getType.Field(i), v)
		for index, col := range cols {
			Fname := getType.Field(i).Name
			if strings.ToUpper(col[:1])+col[1:] == Fname {
				//v := getValue.Field(i)
				// fmt.Println(getType.Field(i).Tag)
				parseQueryColumn(v, data[index])
			}
		}
	}
	return obj
}

func _query(sql string) (rows *sql.Rows, err error) {
	return db.Query(sql)
}

func List(sql, model string) (result []interface{}) {
	//obj = new(obj.(Project))
	v, ok := models[model]
	if !ok {
		return nil
	}
	rows, err := db.Query(sql)
	logs.Log.Debug(" *err: %v,%v", err, rows)
	cols, _ := rows.Columns()
	buff := make([]interface{}, len(cols)) // 临时slice
	data := make([]string, len(cols))      // 存数据slice
	for i, _ := range buff {
		buff[i] = &data[i]
	}
	for rows.Next() {
		obj := reflect.New(v).Interface()
		rows.Scan(buff...) // ...是必须的
		orm(data, cols, obj)
		result = append(result, obj)
	}
	return result
}

func Get(sql, model string) (obj interface{}) {
	objs := List(sql, model)
	if objs != nil {
		obj = objs[0]
	}
	return
}

func Insert(sql, models string, params []interface{}) (obj interface{}) {
	fmt.Println(models)
	stmt, err := db.Prepare(sql)
	res, _ := stmt.Exec(params...)
	id, err := res.LastInsertId()
	fmt.Println(id, err)
	return
}
