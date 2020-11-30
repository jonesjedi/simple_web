package mysql

import (
	"fmt"
	"strconv"
	"sync"

	"onbio/conf"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DBConns sync.Map
	DBConn  *gorm.DB
)

const (
	connect_timeout = 1
	read_timeout    = 3
	write_timeout   = 3
	max_idle_conn   = 200
	max_open_conn   = 3000
	charset         = "utf8"
)

func getDSN(instance string, args ...string) (string, error) {
	temp := "%s:%s@tcp(%s:%s)/%s?charset=%s&loc=Local"
	host, port, database, username, password := getDbConfig(instance)
	dsn := fmt.Sprintf(temp, username, password,
		host, strconv.Itoa(port), database, charset)
	fmt.Println("init mysql:", dsn)
	return dsn, nil
}

func NewMySQL(instance string, args ...string) error {
	dsn, err := getDSN(instance, args...)
	if err != nil {
		return err
	}

	DBConn, err = gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	DBConn.DB().SetMaxIdleConns(max_idle_conn)
	DBConn.DB().SetMaxOpenConns(max_open_conn)
	err = DBConn.DB().Ping()

	if err != nil {
		return err
	}

	return nil
}

func Init(instance string) {

	con, ok := DBConns.Load(instance)
	if ok == false {
		err := NewMySQL(instance)
		if err != nil {
			panic(err.Error())
		}
		DBConns.Store(instance, DBConn)
	} else {
		//切换到当前实例
		DBConn = con.(*gorm.DB)
	}

}

// 根据实例名称，获取对应的连接池
func GetDBConn(instance string) *gorm.DB {
	con, ok := DBConns.Load(instance)
	if ok {
		return con.(*gorm.DB)
	}
	return DBConn
}

func getDbConfig(instance string) (string, int, string, string, string) {
	DBConf := conf.GetMysqlConfig()
	return DBConf.MysqlHost,
		int(DBConf.MysqlPort),
		DBConf.MysqlDBName,
		DBConf.MysqlUserName,
		DBConf.MysqlPwd
}
