package route

import (
	"fmt"
	"net/http"
	"sifu-tool/controller"
	"sifu-tool/ent"
	"sifu-tool/middleware"
	"sifu-tool/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SettingDDNS(api *gin.RouterGroup, secret string, resolvers map[string]map[string]any, ipAPI map[string][]string, entClient *ent.Client, logger *zap.Logger) {
	api.Use(middleware.Jwt(secret, logger))
	ddns := api.Group("/ddns")
	ddns.POST("/interface", func(c *gin.Context) {
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
	ddns.POST("/submit/:method", func(ctx *gin.Context){
		method := ctx.Param("method")
		var form = models.JobForm{}
		if err := ctx.BindJSON(&form); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Json解析失败"})
			return
		}
		if (form.V4method == models.INTERFACE && (form.V4interface == "" || form.Rev4 == "")) || (form.V6method == models.INTERFACE && (form.V6interface == "" || form.Rev6 == "")) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "通过接口获取IP地址时必须指定接口名称以及匹配IP的正则表达式"})
			return
		}
		if (form.V4method == models.SCRIPT && form.V4script == "") || (form.V6method == models.SCRIPT && form.V6script == "") {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "通过脚本获取IP地址时必须指定脚本"})
			return
		}
		if form.Domains == nil || form.Config == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "域名列表和托管商配置不能为空"})
			return
		}
		switch method {
			case "edit":
				id, err := strconv.Atoi(ctx.Query("id"))
				if err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{"message": "字符参数转换Int失败"})
					return
				}
				if id == 0{
					ctx.JSON(http.StatusBadRequest, gin.H{"message": "编辑任务时必须指定任务ID"})
					return
				}
				result, err := controller.EditJobs(form, form.Config[models.RESOLVER], id, ipAPI, resolvers, entClient, logger)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("编辑任务失败: [%s]", err.Error()), "result": result})
					return
				} else {
					ctx.JSON(http.StatusOK, gin.H{"message": "编辑任务成功", "result": result})
					return
				}
			case "add":
				result, err := controller.AddJobs(form, form.Config[models.RESOLVER], ipAPI, resolvers, entClient, logger)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("添加任务失败: [%s]", err.Error()), "result": result})
					return
				} else {
					ctx.JSON(http.StatusOK, gin.H{"message": "添加任务成功", "result": result})
					return
				}
			default:
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "未知操作"})
				return
		}
	})
	ddns.GET("fetch",func(ctx *gin.Context) {
		jobs, err := controller.GetJobs(entClient, logger)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{"message": "获取任务成功", "result": jobs})
			return
		}
	})
	ddns.DELETE("delete",func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Query("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "字符参数转换Int失败"})
			return
		}
		if id == 0{
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "删除任务时必须指定任务ID"})
			return
		}
		if err := controller.DeleteJobs(id, entClient, logger); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{"message": "删除任务成功"})
			return
		}
	})
	ddns.POST("/script", func(ctx *gin.Context){
		script := ctx.PostForm("script")
		result, err := controller.TestScript(script, logger)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("测试脚本失败: [%s]", err.Error()), "result": result})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "测试脚本成功", "result": result})
	})
}