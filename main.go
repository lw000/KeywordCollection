// KeywordCollection project main.go
package main

import (
	"KeywordCollection/global"
	"KeywordCollection/routers/api"
	"KeywordCollection/routers/ws"
	"KeywordCollection/server/dbserver"
	"KeywordCollection/server/parseserver"
	"KeywordCollection/server/queryServer"
	"KeywordCollection/server/taskQueue"
	"fmt"
	_ "github.com/icattlecoder/godaemon"
	"github.com/judwhite/go-svc/svc"
	"github.com/lw000/gocommon/app/gin"
	"github.com/lw000/gocommon/sys"
	"github.com/lw000/gocommon/web/gin/middleware"
	log "github.com/sirupsen/logrus"
	"os"
)

type Program struct {
}

func (p *Program) Init(env svc.Environment) error {
	if env.IsWindowsService() {

	} else {

	}

	var err error
	// 加载全局配置
	if err = global.LoadGlobalConfig(); err != nil {
		return err
	}

	// 启动chrome服务
	if err = global.LoadChromeServer(); err != nil {
		return err
	}

	return nil
}

// Start is called after Init. This method must be non-blocking.
func (p *Program) Start() error {
	// 连接数据库
	var err error
	global.DBReptiledata, err = global.OpenMysql(global.ProjectConfig.MysqlCfg)
	if err != nil {
		return err
	}

	return nil
}

// Stop is called in response to syscall.SIGINT, syscall.SIGTERM, or when a
// Windows Service is stopped.
func (p *Program) Stop() error {
	dbsrv.DBServer().Stop()
	qserv.QueryServer().Stop()
	parseserv.ParseServer().Stop()
	qtasks.ContentRetrievalServer().Stop()
	global.StopChromeServer()
	log.Error("KeywordCollection·服务退出")
	return nil
}

func main() {
	tysys.RegisterOnInterrupt(func(sign os.Signal) {
		dbsrv.DBServer().Stop()
		qserv.QueryServer().Stop()
		parseserv.ParseServer().Stop()
		qtasks.ContentRetrievalServer().Stop()
		global.StopChromeServer()
		log.WithField("sign", fmt.Sprintf("%v", sign)).Error("KeywordCollection·服务退出")
	})

	var err error
	// 加载全局配置
	if err = global.LoadGlobalConfig(); err != nil {
		log.Panic(err)
	}

	// 启动chrome服务
	if err = global.LoadChromeServer(); err != nil {
		log.Panic(err)
	}

	// 连接数据库
	global.DBReptiledata, err = global.OpenMysql(global.ProjectConfig.MysqlCfg)
	if err != nil {
		log.Panic(err)
	}

	// service.TestDao()

	dbsrv.DBServer().Start()
	qserv.QueryServer().Start()
	parseserv.ParseServer().Start()
	qtasks.ContentRetrievalServer().Start()

	app := tygin.NewApplication(global.ProjectConfig.Debug)
	app.SetEnableTLS(global.ProjectConfig.TLS.Enable)
	if app.EnableTLS() {
		app.SetTlsFile(global.ProjectConfig.TLS.CertFile, global.ProjectConfig.TLS.KeyFile)
	}

	err = app.Run(global.ProjectConfig.Port, func(app *tygin.WebApplication) {
		app.Engine().Use(tymiddleware.CorsHandler(nil))
		ws.RegisterService(app.Engine())
		api.RegiserService(app.Engine())
	})
	log.Panic(err)

	// pro := &Program{}
	// if err := svc.Run(pro); err != nil {
	// 	log.Error(err)
	// }
}
