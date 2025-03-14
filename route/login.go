package route

import (
	"net/http"
	"sifu-tool/controller"
	"sifu-tool/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SettingLogin(api *gin.RouterGroup, user models.User, logger *zap.Logger) {
	api.POST("/login", func(ctx *gin.Context) {
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		if username == user.Username && password == user.Password {
			token, err := controller.Login(user.Secret, logger)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"message": token})
			return
		}
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "用户名或密码错误"})
	})
	api.GET("/verify", func(ctx *gin.Context) {
		authorization := ctx.GetHeader("Authorization")
		token, err := controller.Verify(authorization, user.Secret, logger)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": token})
	})
}