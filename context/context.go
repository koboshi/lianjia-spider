package context

import (
	"errors"
	"github.com/koboshi/mole/database"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	SpiderType string `ini:"spider_type"`
	SpiderRole string `ini:"spider_role"`

	LogOn int `ini:"log_on"`
	LogLevel string `ini:"log_level"`
	LogDir string `ini:"log_dir"`

	MysqlHost string `ini:"mysql_host"`
	MysqlUser string `ini:"mysql_user"`
	MysqlPsw string `ini:"mysql_psw"`
	MysqlSchema string `ini:"mysql_schema"`
	MysqlCharset string `ini:"mysql_charset"`
	MysqlMaxConn int `ini:"mysql_max_conn"`
	MysqlIdleConn int `ini:"mysql_idle_conn"`

	NetConcurrent int `ini:"net_concurrent"`
	NetInterface string `ini:"net_interface"`
}

var ErrLoadConf = errors.New("could not load config file")
var ErrConnectMysql = errors.New("could not connect mysql")
var ErrInitPool = errors.New("could not init net concurrent pool")

var conf Config
var db *database.Database
var logger *log.Logger


func init() {
	//读取配置文件
	//默认读取 conf/spider.conf
	//支持传入指定路径的配置文件
	var confPath string
	if len(os.Args) > 1 {
		confPath = os.Args[1]
	}else {
		confPath = filepath.Dir(os.Args[0]) + "/conf/spider.conf"
	}
	var err error
	conf, err = load(confPath)
	if err != nil {
		panic(ErrLoadConf)
	}

	//初始化日志
	log.SetPrefix("LOG:")
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	//初始化数据库连节

	//初始化爬虫配置

	//初始化网络配置
}

func initLog() {
	//TODO
}

func initMysql() {
	//TODO
}

func initNet() {
	//TODO
}

func initSpider() {
	//TODO
}

func load(path string) (Config, error) {
	var config Config
	conf, err := ini.Load(path)   //加载配置文件
	if err != nil {
		return config, err
	}
	conf.BlockMode = false
	err = conf.MapTo(&config)   //解析成结构体
	if err != nil {
		return config, err
	}
	return config, nil
}