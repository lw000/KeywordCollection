package global

import (
	"KeywordCollection/chrome"
	"KeywordCollection/config"
	"KeywordCollection/constant"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	tyIdWorker "github.com/lw000/gocommon/IdWorker"
	tymysql "github.com/lw000/gocommon/db/mysql"
	tyrdsex "github.com/lw000/gocommon/db/rdsex"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"path"
	"sync"
	"time"
)

var (
	// ProjectConfig 工程配置
	ProjectConfig *config.IniConfig
	// DBSrv 数据库实例
	DBReptiledata *tymysql.Mysql
	// Chromes ...
	Chromes map[string]*chrome.ChromeDriver
	// QueryClients 查询客户端
	QueryClients sync.Map
	GormSql      *gorm.DB
)

var (
	idworker *tyIdWorker.IdWorker
	once     sync.Once
)

// config logrus log to local filesystem, with file rotation
func configLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	baseLogPath := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d_%H%M",
		// rotatelogs.WithLinkName(baseLogPath), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge), // 文件最大保存时间
		// rotatelogs.WithRotationCount(365),  // 最多存365个文件
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)

	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})

	// log.SetLevel(log.ErrorLevel)

	log.SetReportCaller(true)

	log.AddHook(lfHook)
}

func init() {

}

func GetIdWorker() *tyIdWorker.IdWorker {
	once.Do(func() {
		idworker = &tyIdWorker.IdWorker{}
		_ = idworker.Start(1)
	})
	return idworker
}

// LoadGlobalConfig 加载全局配置文件
func LoadGlobalConfig() error {
	var err error
	ProjectConfig, err = config.LoadIniConfig("conf/conf.ini")
	if err != nil {
		log.Error(err)
		return err
	}
	var logName = "KeywordCollection"
	// 日志分割 1按天分割，2按周分割, 3 按月分割，4按年分割
	switch ProjectConfig.SplitLog {
	case 1:
		configLocalFilesystemLogger("log", logName, time.Hour*24*365, time.Hour*24)
	case 2:
		configLocalFilesystemLogger("log", logName, time.Hour*24*365, time.Hour*24*7)
	case 3:
		configLocalFilesystemLogger("log", logName, time.Hour*24*365, time.Hour*24*30)
	case 4:
		configLocalFilesystemLogger("log", logName, time.Hour*24*365, time.Hour*24*365)
	default:
		configLocalFilesystemLogger("log", logName, time.Hour*24*365, time.Hour*24)
	}

	return nil
}

// OpenMysql 打开数据库实例
func OpenMysql(cfg *tymysql.JsonConfig) (*tymysql.Mysql, error) {
	srv := &tymysql.Mysql{}
	if err := srv.OpenWithJsonConfig(cfg); err != nil {
		log.Error(err)
		return nil, err
	}
	log.Info("数据库连接成功")

	return srv, nil
}

func OpenGORM(cfg *tymysql.JsonConfig) (*gorm.DB, error) {
	// "user:password@/dbname?charset=utf8&parseTime=True&loc=Local"
	var err error
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", cfg.Username, cfg.Password, cfg.Host, cfg.Database)
	GormSql, err = gorm.Open("mysql", dns)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return GormSql, nil
}

// OpenRedis 打开redis实例
func OpenRedis(cfg *tyrdsex.JsonConfig) (*tyrdsex.RdsServer, error) {
	srv := &tyrdsex.RdsServer{}
	if err := srv.OpenWithJsonConfig(cfg); err != nil {
		log.Error(err)
		return nil, err
	}
	log.Info("REDIS接成功")
	return srv, nil
}

// LoadChromeServer 加载chrome服务
func LoadChromeServer() error {
	Chromes = make(map[string]*chrome.ChromeDriver)
	Chromes[constant.DevicePc] = &chrome.ChromeDriver{}
	Chromes[constant.DeviceMobile] = &chrome.ChromeDriver{}
	for k, s := range Chromes {
		if k == constant.DevicePc {
			err := s.NewPcService(9950)
			if err != nil {
				log.Error(err)
				return err
			}
		}
		if k == constant.DeviceMobile {
			err := s.NewMobileService(9951)
			if err != nil {
				log.Error(err)
				return err
			}
		}
	}
	return nil
}

// StopChromeServer 管理chrome服务
func StopChromeServer() {
	for _, s := range Chromes {
		s.Stop()
	}
}
