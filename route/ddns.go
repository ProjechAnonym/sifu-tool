package route

import (
	"net/http"
	"sifu-tool/controller"
	"sifu-tool/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SettingDDNS(api *gin.RouterGroup, secret string, logger *zap.Logger) {
	api.POST("/interface", middleware.Jwt(secret, logger), func(c *gin.Context) {
		inerfaceIPs, errors := controller.GetInterfaceIPs(logger)
		errs := make([]string, len(errors))
		if errors != nil {
			for i, err := range errors {
				errs[i] = err.Error()
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": inerfaceIPs, "errors": errs})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": inerfaceIPs})
	})
	api.POST("/ddns", middleware.Jwt(secret, logger), func(ctx *gin.Context){
		
	})
}