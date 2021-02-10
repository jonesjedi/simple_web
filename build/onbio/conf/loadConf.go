package conf

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

var config = NewConf()

type MysqlConf struct {
	MysqlHost     string `json:"mysql_host"`
	MysqlPort     int    `json:"mysql_port"`
	MysqlUserName string `json:"mysql_user_name"`
	MysqlPwd      string `json:"mysql_pwd"`
	MysqlEncoding string `json:"mysql_encoding"`
	MysqlDBName   string `json:"mysql_db_name"`
}

type LogConf struct {
	FilePath string `json:"file_path"`
	FileName string `json:"file_name"`
}

type RedisConf struct {
	RedisHost string `json:"redis_host"`
	RedisPort int    `json:"redis_port"`
	RedisPwd  string `json:"redis_pwd"`
}

type Conf struct {
	LogConfDetail   LogConf   `json:"log_conf"`
	MysqlConfDetail MysqlConf `json:"mysql_conf"`
	RedisConfDetail RedisConf `json:"redis_conf"`
}

func NewConf() *Conf {

	return &Conf{}
}

func LoadConf(fileName string) error {

	//判断路径是否为空
	if fileName == "" {
		log.Fatal("conf file is empty")
	}
	file, error := os.Open(fileName)
	if error != nil {
		return error
	}
	defer file.Close()
	fileContent, error := ioutil.ReadAll(file)
	conf := string(fileContent)
	error = json.Unmarshal([]byte(conf), config)
	if error != nil {
		log.Fatal("conf file is not json format: ", conf, ",err:", error)
		return errors.New("conf file is not json format")
	}
	return nil
}
func GetLogConfig() (conf LogConf) {
	return config.LogConfDetail
}

func GetMysqlConfig() (conf MysqlConf) {
	return config.MysqlConfDetail
}

func GetRedisConfig() (conf RedisConf) {
	return config.RedisConfDetail
}
