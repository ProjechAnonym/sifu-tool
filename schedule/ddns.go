package schedule

import (
	"context"
	"fmt"
	"net/http"
	"sifu-tool/ddns"
	"sifu-tool/ent"
	"sifu-tool/models"
	"sync"

	"go.uber.org/zap"
)

func DDNSJob(entClient *ent.Client, client *http.Client, ipAPI map[string][]string, resolvers map[string]map[string]any, logger *zap.Logger) {
	ddnsJobs, err := entClient.DDNS.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取任务失败: [%s]", err.Error()))
		return
	}
	if len(ddnsJobs) == 0 {
		logger.Debug("没有DDNS任务")
		return
	}
	var jobs sync.WaitGroup

	for _, ddnsJob := range ddnsJobs {
		jobs.Add(1)
		go func(){
			defer jobs.Done()
			setDDNS(ddnsJob, entClient, client, ipAPI, resolvers, logger)
		}()
	}
	jobs.Wait()
}

func setDDNS(ddnsJob *ent.DDNS, entClient *ent.Client, client *http.Client, ipAPI map[string][]string, resolvers map[string]map[string]any, logger *zap.Logger){

	// 获取IPv4和IPv6API接口
	ipv4api, ok := ipAPI["ipv4"]
	if !ok {
		logger.Error("未配置IPV4接口")
		return
	}
	ipv6api, ok := ipAPI["ipv6"]
	if !ok {
		logger.Error("未配置IPV6接口")
		return
	}
	// 按照不同方法获取IPv4和IPv6地址
	ipv4, err := ddns.GetIP(ddnsJob.V4method, ddnsJob.V4interface, ddnsJob.Rev4, models.Rev4, ddnsJob.V4script, client, logger, ipv4api...)
	if err != nil {
		logger.Error(fmt.Sprintf("获取IP地址失败: [%s]", err.Error()))
	}

	ipv6, err := ddns.GetIP(ddnsJob.V6method, ddnsJob.V6interface, ddnsJob.Rev6, models.Rev6, ddnsJob.V6script, client, logger, ipv6api...)
	if err != nil {
		logger.Error(fmt.Sprintf("获取IP地址失败: [%s]", err.Error()))
	}

	domains := make([]models.Domain, len(ddnsJob.Domains))
	if ipv4 == ddnsJob.Ipv4 && ipv6 == ddnsJob.Ipv6 {
		for i, domain :=  range(ddnsJob.Domains){
			domain.Status = models.STATIC
			domain.Result = "本地记录未发生变化"
			domains[i] = domain
		}
		if err := entClient.DDNS.UpdateOneID(ddnsJob.ID).SetDomains(domains).Exec(context.Background()); err != nil {
			logger.Error(fmt.Sprintf("保存任务失败: [%s]", err.Error()))
			return
		}
		return
	}
	

	switch ddnsJob.Config[models.RESOLVER]{
		case models.CF:
			// 根据不同的DNS服务商断言不同的配置类型以获取所需信息
			cloudflare, ok := resolvers[ddnsJob.Config[models.RESOLVER]]["api"].(string)
			if !ok {
				logger.Error(`未能读取"cloudflare"接口信息`)
				return
			}
			
			results, err := ddns.CloudFlare(cloudflare, ddnsJob.Config[models.CFTOKEN], ddnsJob.Domains, client, logger)
			if err != nil {
				// 该函数返回的错误为初始化错误，若出错则整个域名列表都不会更新
				// 因此将每条域名的结果都设置为该错误
				logger.Error(fmt.Sprintf("更新域名失败: [%s]", err.Error()))
				for i, domain := range results {
					domain.Result = fmt.Sprintf("更新域名失败: [%s]", err.Error())
					domain.Status = models.FAILURE
					domains[i] = domain
				}
			}else{
				domains = make([]models.Domain, 0)
				for _, domain := range results {
					if domain.Status != models.DELETE {
						domains = append(domains, domain)
					}
				}
			}
	}
	if err := entClient.DDNS.UpdateOneID(ddnsJob.ID).SetIpv4(ipv4).SetIpv6(ipv6).SetDomains(domains).Exec(context.Background()); err != nil {
		logger.Error(fmt.Sprintf("保存任务失败: [%s]", err.Error()))
		return
	}
}
