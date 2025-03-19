package main

import (
	"fmt"
	"os"
	"sifu-tool/ddns"
	"sifu-tool/ent"
	"sifu-tool/initial"
	"sifu-tool/middleware"
	"sifu-tool/models"
	"sifu-tool/route"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)
var taskLogger *zap.Logger
var webLogger *zap.Logger
var entClient *ent.Client
var buntClient *buntdb.DB
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
	buntClient = initial.InitBuntdb(initLogger)
	initLogger.Info("初始化buntdb内存数据库客户端完成")
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
}

func main()  {
	defer func() {
		taskLogger.Sync()
		entClient.Close()	
		buntClient.Close()
	}()
	// domais := []models.Domain{
	// 	models.Domain{Domain: "1.lzhlovelcl.top", Type: "A", TTL: 120, Value: "192.168.50.2",},
	// 	models.Domain{Domain: "3.lzhlovelcl.top", Type: "AAAA", TTL: 120, Value: "240e:379:17f:7400:be24:11ff:fe4d:994c"},
	// 	models.Domain{Domain: "4.lzhlovelcl.top", Type: "AAAA",TTL: 120,Value: "240e:379:17f:7400:be24:11ff:fe4d:994c"},
	// 	models.Domain{Domain: "3.lzhlovelcl.top", Type: "A", Value: "3.3.3.3",TTL: 120},
	// 	models.Domain{Domain: "lzhlovelcl.top", Type: "AAAA", Value: "240e:379:17f:7400:be24:11ff:fe4d:994c",TTL: 120},
	// 	models.Domain{Domain: "*.lzhlovelcl.top", Type: "AAAA", Value: "240e:379:17f:7400:be24:11ff:fe4d:994c",TTL: 120},
	// }
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	server.Use(middleware.Logger(webLogger), middleware.Recovery(true, webLogger), cors.New(middleware.Cors(*domains)))
	api := server.Group("/api")
	route.SettingLogin(api, setting.User, webLogger)
	route.SettingDDNS(api, setting.User.Secret, setting.DDNS.Resolver, entClient, webLogger)
	
	a,_ := ddns.IPfromInterface("enp6s18",`^fe.*$`, webLogger)
	fmt.Println(a)
	// fmt.Println(setting.DDNS.Resolver[models.CF]["api"])
	// client := http.DefaultClient
	// b,e:=ddns.CloudFlare(setting.DDNS.Resolver[models.CF]["api"], "_YnGvSYspTxw1zEBMptWVvBqLWRxPMVR1_M1dsqm",domais, client, webLogger)
	// for _, i := range(b){
	// 	fmt.Println(i)
	// }
	// fmt.Println(e)
	
	if setting.Server.Tls != nil {
		server.RunTLS(fmt.Sprintf(":%d", setting.Server.Tls.Port), setting.Server.Tls.Cert, setting.Server.Tls.Key)
	}
	server.Run(*address)
}