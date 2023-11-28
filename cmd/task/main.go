package main

import (
	"fmt"
	"github.com/chainxx/bitx/internal/conf"
	"github.com/chainxx/bitx/internal/task"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/startopsz/rule/pkg/os/env"
	"os"
	"os/signal"
	"syscall"
)

var (
	name                     string
	version                  string
	configPath               string
	region                   string
	environment              string
	openTelemetryTraceEnable bool
	openTelemetryEndpoint    string
	id, _                    = os.Hostname()
)

func init() {
	name = "startops.message.service"
	version = "1.0.0"
	configPath = env.GetConfigPath()
	region = env.GetRegion()
	environment = env.GetEnvironment()
	openTelemetryTraceEnable = env.GetOpenTelemetryTraceEnable()
	openTelemetryEndpoint = env.GetOpenTelemetryEndpoint()
	
	fmt.Printf("name: %s.\nversion: %s.\nconfigPath: %s.\nregion: %s.\nenvironment: %s.\nopenTelemetryTraceEnable: %t.\nopenTelemetryEndpoint: %s.\n\n",
		name, version, configPath, region, environment, openTelemetryTraceEnable, openTelemetryEndpoint)
}

type app struct {
	task *task.Task
	log  *log.Helper
}

func newApp(logger log.Logger, task *task.Task) *app {
	return &app{
		task: task,
		log:  log.NewHelper(logger),
	}
}

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
	
	application, clean, err := wireApp(bc.Data, bc.Binance, logger)
	if err != nil {
		panic(fmt.Sprintf("启动服务失败, err: %s", err.Error()))
	}
	
	defer clean()
	
	application.task.Cronjob()
	
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	
	exitChan := make(chan int)
	
	go func() {
		for {
			s := <-signalChan
			switch s {
			// kill -SIGHUP XXXX
			case syscall.SIGHUP:
				fmt.Println("hungup")
				exitChan <- 0
			
			// kill -SIGINT XXXX or Ctrl+c
			case syscall.SIGINT:
				fmt.Println("Warikomi")
				exitChan <- 0
			
			// kill -SIGTERM XXXX
			case syscall.SIGTERM:
				fmt.Println("force stop")
				exitChan <- 0
			
			// kill -SIGQUIT XXXX
			case syscall.SIGQUIT:
				fmt.Println("stop and core dump")
				exitChan <- 0
			
			default:
				fmt.Println("Unknown signal.")
				exitChan <- 1
			}
		}
	}()
	
	code := <-exitChan
	os.Exit(code)
}
