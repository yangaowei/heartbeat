package orm

import (
	"fmt"
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

func (p *Project) TableName() string {

	return "vtc_heartbeat_project"
}

func TestORM(t *testing.T) {
	RegisterDataBase("default")
	fmt.Println(Orm.DbInfo())
	project := &Project{Id: 1}
	p := Orm.Get(project)
	fmt.Println(p)
}
