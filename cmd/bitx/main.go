package main

import (
	"fmt"
	"github.com/chainxx/bitx/internal/conf"
	"github.com/chainxx/bitx/internal/server"
	"github.com/chainxx/bitx/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/startopsz/rule/pkg/os/env"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"os"
	
	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	name                     string
	version                  string
	configPath               string
	region                   string
	environment              string
	etcdEndpoints            []string
	openTelemetryTraceEnable bool
	openTelemetryEndpoint    string
	id, _                    = os.Hostname()
)

func init() {
	
	name = "bitx"
	version = "1.0.0"
	etcdEndpoints = env.GetEtcdEndpoints()
	configPath = env.GetConfigPath()
	region = env.GetRegion()
	environment = env.GetEnvironment()
	openTelemetryTraceEnable = env.GetOpenTelemetryTraceEnable()
	openTelemetryEndpoint = env.GetOpenTelemetryEndpoint()
	
	fmt.Printf("name: %s.\nversion: %s.\nconfigPath: %s.\nregion: %s.\nenvironment: %s.\netcdEndpoints: %s.\nopenTelemetryTraceEnable: %t.\nopenTelemetryEndpoint: %s.\n",
		name, version, configPath, region, environment, etcdEndpoints, openTelemetryTraceEnable, openTelemetryEndpoint)
}

type app struct {
	bitxService *service.BitxService
}

func newApp(bitxService *service.BitxService) *app {
	return &app{
		bitxService: bitxService,
	}
}

//	@title			StartOps
//	@version		1.0
//	@description	StartOps运维平台
//	@contact.name	qx.liu
//	@BasePath		/
func main() {
	
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", "id",
		"service.name", "Name",
		"service.version", "Version",
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	
	var configs []config.Source
	// 本地配置文件配置
	if configPath != "" {
		configs = append(configs, file.NewSource(configPath))
		
	}
	
	if len(configs) == 0 {
		panic("缺少配置信息")
	}
	c := config.New(
		config.WithSource(
			configs...,
		),
	)
	
	if err := c.Load(); err != nil {
		panic(err)
	}
	
	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	
	iapp, clean, err := wireApp(bc.Data, bc.Binance, logger)
	if err != nil {
		panic(err)
	}
	defer clean()
	pid := int64(os.Getpid())
	
	if openTelemetryTraceEnable {
		server.InitProvider(openTelemetryEndpoint, name, region, environment, pid)
	}
	
	route := gin.New()
	
	gin.SetMode(gin.DebugMode)
	
	// openTelemetry
	route.Use(otelgin.Middleware(name))
	//route.Use(gin.Logger())
	route.Use(gin.Recovery())
	
	route.GET("/symbol", iapp.bitxService.ListSymbol)
	route.GET("/asset/daily", iapp.bitxService.GetDailyAsset)
	route.GET("/price", iapp.bitxService.ListSymbolPrice)
	
	route.GET("/lotto/doubleColorBall/count", iapp.bitxService.GetDoubleColorBallCount)
	
	err = route.Run(":8010")
	
	if err != nil {
		errMessage := fmt.Sprint("Start Gateway Server Error,", err)
		panic(errMessage)
	}
}
