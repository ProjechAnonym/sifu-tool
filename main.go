package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"log"
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
	route.SettingDDNS(api, setting.User.Secret, setting.DDNS.Resolver, map[string][]string{"ipv4": setting.DDNS.V4API, "ipv6": setting.DDNS.V6API}, entClient, webLogger)
	a,_ := ddns.IPfromInterface("enp6s18",`^fe.*$`, webLogger)
	fmt.Println(a)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}	
	b, _ := x509.MarshalECPrivateKey(privateKey)
	block, _ := aes.NewCipher([]byte("1234567890123456"))
	gcm, _ := cipher.NewGCM(block)

	nonce := make([]byte, gcm.NonceSize())

	ciphertext := gcm.Seal(nonce, nonce, b, nil)
	k := base64.StdEncoding.EncodeToString(ciphertext)
	ciphertext1, _ := base64.StdEncoding.DecodeString(k)
	block1, _ := aes.NewCipher([]byte("1234567890123456"))


	gcm1, _ := cipher.NewGCM(block1)
	

	nonceSize := gcm1.NonceSize()


	nonce, ciphertextBytes := ciphertext1[:nonceSize], ciphertext1[nonceSize:]
	plaintext, _ := gcm.Open(nil, nonce, ciphertextBytes, nil)
	d,_ := x509.ParseECPrivateKey(plaintext)
	fmt.Println(d.Equal(privateKey))
	if setting.Server.Tls != nil {
		server.RunTLS(fmt.Sprintf(":%d", setting.Server.Tls.Port), setting.Server.Tls.Cert, setting.Server.Tls.Key)
	}
	server.Run(*address)
}
func encrypt(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt 解密数据
func decrypt(ciphertext string, key []byte) ([]byte, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertextBytes) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertextBytes := ciphertextBytes[:nonceSize], ciphertextBytes[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}