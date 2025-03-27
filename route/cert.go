package route

import (
	"net/http"
	"sifu-tool/ent"
	"sifu-tool/middleware"
	"sifu-tool/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SettingAcme(api *gin.RouterGroup, secret, key string, entClient *ent.Client, logger *zap.Logger){
	acme := api.Group("/acme")
	acme.Use(middleware.Jwt(secret, logger))
	acme.POST("/submit/:method", func(ctx *gin.Context) {
		method := ctx.Param("method")
		switch method {
			case "add":
				form := models.AcmeForm{}
				if err := ctx.BindJSON(&form); err != nil{
					ctx.JSON(http.StatusBadRequest, gin.H{"message": "Json解析失败"})
					return
				}
				
		}
	})
}