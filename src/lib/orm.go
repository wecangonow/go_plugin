package lib

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // import your used driver
)

var (
	O  orm.Ormer
)


func InitOrm() {
	orm.RegisterDataBase("default", "mysql", AppConfig.DbConnectstr, 30)
	O = orm.NewOrm()
	orm.Debug = true
}
