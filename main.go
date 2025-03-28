package main

import (
	"fmt"
	"net/http"
	"os"
	"sifu-tool/ddns"
	"sifu-tool/ent"
	"sifu-tool/initial"
	"sifu-tool/middleware"
	"sifu-tool/models"
	"sifu-tool/route"
	"sifu-tool/schedule"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)
var taskLogger *zap.Logger
var webLogger *zap.Logger
var entClient *ent.Client

var setting models.Setting
var environment, config, address *string
var domains *[]string
func init() {
	environment, config, address, domains = initial.Initkingpin() 
	initLogger := initial.GetLogger(*environment, "init")
	defer initLogger.Sync()
	taskLogger = initial.GetLogger(*environment, "task")
	webLogger = initial.GetLogger(*environment, "web")
	entClient = initial.InitEntdb(*environment, initLogger)
	initLogger.Info("初始化sqlite数据库客户端完成")

	viper.SetConfigFile(*config)
	if err := viper.ReadInConfig(); err != nil {
		initLogger.Fatal(fmt.Sprintf("读取配置文件失败: %s", err.Error()))
		os.Exit(1)
	}
	if err := viper.Unmarshal(&setting); err != nil {
		initLogger.Fatal(fmt.Sprintf("解析配置文件失败: %s",err.Error()))
		os.Exit(1)
	}
	initLogger.Info("加载配置文件完成")
	ddnsSchedule := cron.New()
	ddnsSchedule.AddFunc("@every 1m", func() {
		client := http.DefaultClient
		schedule.DDNSJob(entClient, client, map[string][]string{"ipv4": setting.DDNS.V4API, "ipv6": setting.DDNS.V6API}, setting.DDNS.Resolver, taskLogger)})
	ddnsSchedule.Start()
}

func main()  {
	defer func() {
		taskLogger.Sync()
		entClient.Close()	
	}()
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	server.Use(middleware.Logger(webLogger), middleware.Recovery(true, webLogger), cors.New(middleware.Cors(*domains)))
	api := server.Group("/api")
	route.SettingLogin(api, setting.User, webLogger)
	route.SettingAcme(api, setting.User.Secret, setting.User.Key, entClient, webLogger)
	route.SettingDDNS(api, setting.User.Secret, setting.DDNS.Resolver, map[string][]string{"ipv4": setting.DDNS.V4API, "ipv6": setting.DDNS.V6API}, entClient, webLogger)
	a,_ := ddns.IPfromInterface("enp6s18",`^fe.*$`, webLogger)
	fmt.Println(a)
	
	if setting.Server.Tls != nil {
		server.RunTLS(fmt.Sprintf(":%d", setting.Server.Tls.Port), setting.Server.Tls.Cert, setting.Server.Tls.Key)
	}
	server.Run(*address)
}