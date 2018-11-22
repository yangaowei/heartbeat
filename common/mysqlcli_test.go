package common

import (
	"../logs"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
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

func SetDefault(f reflect.StructField, v reflect.Value) {
	default_value := f.Tag.Get("default")
	parseQueryColumn(v, default_value)
}

func orm(data, cols []string) (project *Project) {
	fmt.Println(data, cols)
	project = &Project{}
	getType := reflect.TypeOf(project).Elem()
	getValue := reflect.ValueOf(project).Elem()
	for i := 0; i < getValue.NumField(); i++ {
		//fmt.Println(getType.Field(i))
		//getType.Field(i)
		v := getValue.Field(i)
		SetDefault(getType.Field(i), v)
		for index, col := range cols {
			Fname := getType.Field(i).Name
			if strings.ToUpper(col[:1])+col[1:] == Fname {
				//v := getValue.Field(i)
				fmt.Println(getType.Field(i).Tag)
				parseQueryColumn(v, data[index])
			}
		}
	}
	return
}

func TestConn(t *testing.T) {

	defer func() {
		db.Close()
	}()
	logs.Log.Debug(" *Start: %v", "baidu.com")
	logs.Log.Debug(" *Start: %v", "baidu.com")

	db, err := DB()
	logs.Log.Debug(" *conn: %v", db)
	logs.Log.Debug(" *err: %v", err)

	rows, err := db.Query("SELECT queue FROM vtc_heartbeat_project order by id")
	logs.Log.Debug(" *err: %v", err)
	cols, _ := rows.Columns()
	buff := make([]interface{}, len(cols)) // 临时slice
	data := make([]string, len(cols))      // 存数据slice
	for i, _ := range buff {
		buff[i] = &data[i]
	}
	for rows.Next() {
		rows.Scan(buff...) // ...是必须的
		project := orm(data, cols)
		fmt.Println(project)
		// for k, values := range data {
		// 	//根据 colum获取字段名称
		// 	fmt.Println(cols[k], values)
		// 	//break
		// }
		break
	}

}

type User struct {
	Name string
	Age  int
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

func TestRe(t *testing.T) {
	user := User{Name: "yan"}
	fmt.Println(user)
	getType := reflect.TypeOf(user)
	fmt.Println(getType)
	getValue := reflect.ValueOf(user)
	fmt.Println(getValue)
	typ := getValue.Type().Name()
	fmt.Println(typ)
	// for i := 0; i < getType.NumField(); i++ {
	// 	field := getType.Field(i)
	// 	value := getValue.Field(i)
	// 	//v = getValue.Field(i)
	// 	fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
	// 	fmt.Println(value.Kind())
	// }
}
