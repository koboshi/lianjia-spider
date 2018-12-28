package context

import (
	"errors"
	"github.com/koboshi/mole/database"
	"github.com/koboshi/mole/work"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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
var ErrInitMysql = errors.New("could not connect mysql")
var ErrInitLog = errors.New("could not init logger")
var ErrInitPool = errors.New("could not init net concurrent pool")

var conf Config
var Db *database.Database
var NetRoutinePool *work.Pool
var TraceLogger *log.Logger
var InfoLogger *log.Logger
var WarnLogger *log.Logger
var ErrorLogger *log.Logger

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
	//初始化
	initLog()
	initMysql()
	initNet()
	initSpider()
}

func initLog() {
	log.SetPrefix("LOG:")
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	if conf.LogOn < 1 {
		return
	}
	//初始化日志记录器
	var logDir string
	if len(conf.LogDir) > 0 {
		logDir = conf.LogDir
	}else {
		logDir = filepath.Dir(os.Args[0]) + "/logs/"
	}
	var logLevelMap = make(map[string]int)
	logLevelMap["trace"] = 1
	logLevelMap["info"] = 2
	logLevelMap["warn"] = 3
	logLevelMap["error"] = 4
	destLevel, ok := logLevelMap[strings.ToLower(conf.LogLevel)]
	if !ok {
		destLevel = 4
	}

	if destLevel <= logLevelMap["trace"] {
		traceFile, err := os.OpenFile(logDir + "/trace.log", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
		if err != nil {
			panic(ErrInitLog)
		}
		TraceLogger = log.New(traceFile, "Trace:", log.Ldate | log.Ltime | log.Llongfile)
	}else {
		TraceLogger = log.New(ioutil.Discard, "Trace:", log.Ldate | log.Ltime | log.Llongfile)
	}

	if destLevel <= logLevelMap["info"] {
		infoFile, err := os.OpenFile(logDir + "/info.log", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
		if err != nil {
			panic(ErrInitLog)
		}
		InfoLogger = log.New(infoFile, "Info:", log.Ldate | log.Ltime | log.Llongfile)
	}else {
		InfoLogger = log.New(ioutil.Discard, "Info:", log.Ldate | log.Ltime | log.Llongfile)
	}

	if destLevel <= logLevelMap["warn"] {
		warnFile, err := os.OpenFile(logDir + "/warn.log", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
		if err != nil {
			panic(ErrInitLog)
		}
		WarnLogger = log.New(warnFile, "Warn:", log.Ldate | log.Ltime | log.Llongfile)
	}else {
		WarnLogger = log.New(ioutil.Discard, "Warn:", log.Ldate | log.Ltime | log.Llongfile)
	}

	if destLevel <= logLevelMap["error"] {
		errorFile, err := os.OpenFile(logDir + "/error.log", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
		if err != nil {
			panic(ErrInitLog)
		}
		ErrorLogger = log.New(errorFile, "Error:", log.Ldate | log.Ltime | log.Llongfile)
	}else {
		ErrorLogger = log.New(ioutil.Discard, "Error:", log.Ldate | log.Ltime | log.Llongfile)
	}

}

func initMysql() {
	host := conf.MysqlHost
	username := conf.MysqlUser
	password := conf.MysqlPsw
	schema := conf.MysqlSchema
	charset := conf.MysqlCharset

	var err error
	Db, err = database.New(host, username, password, schema, charset)
	if err != nil {
		panic(ErrInitMysql)
	}
}

func initNet() {
	var err error
	NetRoutinePool, err = work.New(conf.NetConcurrent)
	if err != nil {
		panic(ErrInitPool)
	}
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