package route

import (
	"fmt"
	"net/http"
	"sifu-tool/controller"
	"sifu-tool/ent"
	"sifu-tool/middleware"
	"sifu-tool/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SettingDDNS(api *gin.RouterGroup, secret string, resolver map[string]map[string]string, entClient *ent.Client, logger *zap.Logger) {
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
		var form = models.JobForm{}
		if err := ctx.BindJSON(&form); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Json解析失败"})
			return
		}
		if (form.V4method == models.INTERFACE && (form.V4interface == "" || form.IPV4 == "" || form.Rev4 == "")) || (form.V6method == models.INTERFACE && (form.V6interface == "" || form.IPV6 == "" || form.Rev6 == "")) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "通过接口获取IP地址时必须指定接口名称以及匹配IP的正则表达式"})
			return
		}
		if (form.V4method == models.SCRIPT && (form.V4script == "" || form.IPV4 == "")) || (form.V6method == models.SCRIPT && (form.V6script == "" || form.IPV6 == "")) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "通过脚本获取IP地址时必须指定脚本以及IP"})
			return
		}
		if form.Domains == nil || form.Config == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "域名列表和托管商配置不能为空"})
			return
		}

		result, err := controller.AddJobs(form, form.Config[models.RESOLVER], resolver[form.Config[models.RESOLVER]][models.CFAPI], entClient, logger)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("添加任务失败: [%s]", err.Error())})
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{"message": "添加任务成功", "result": result})
			return
		}
	})
}