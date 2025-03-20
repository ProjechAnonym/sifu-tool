package controller

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sifu-tool/ddns"
	"sifu-tool/ent"
	entddns "sifu-tool/ent/ddns"
	"sifu-tool/models"

	"go.uber.org/zap"
)

func GetInterfaceIPs(logger *zap.Logger) (map[string][]string, []error){
	interfaces, err := net.Interfaces()
	if err != nil {
		logger.Error(fmt.Sprintf("获取网卡信息失败: [%s]", err.Error()))
		return nil, []error{fmt.Errorf("获取网卡信息失败")}
	}
	interfaceIPs := make(map[string][]string)
	var errors []error
	for _, netInterface := range interfaces {
		addrs, err := netInterface.Addrs()
		if err != nil {
			logger.Error(fmt.Sprintf("获取网卡地址失败: [%s]", err.Error()))
			errors = append(errors, fmt.Errorf(`获取网卡"%s"IP地址失败`, netInterface.Name))
		}
		addresses := make([]string, len(addrs))
		for i, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				logger.Error(fmt.Sprintf(`解析网卡"%s"地址"%s"失败: [%s]`, netInterface.Name, addr.String(), err.Error()))
				errors = append(errors, fmt.Errorf(`解析网卡"%s"地址"%s"失败`, netInterface.Name, addr.String()))
			}
			addresses[i] = ip.String()
		}
		interfaceIPs[netInterface.Name] = addresses
	}
	return interfaceIPs, errors
}

func AddJobs(form models.JobForm, resolver, api string, ipAPI map[string][]string, entClient *ent.Client, logger *zap.Logger) (models.JobForm, error) {
	var ipv4 string
	var ipv6 string
	var err error
	client := http.DefaultClient
	domains := make([]models.Domain, len(form.Domains))
	switch form.V4method {
		case models.INTERFACE:
			ipv4, err = ddns.IPfromInterface(form.V4interface, form.Rev4, logger)
			if err != nil {
				logger.Error(fmt.Sprintf("获取IPV4地址失败: [%s]", err.Error()))
			}
		case models.IPAPI:
			ipv4api, ok := ipAPI["ipv4"]
			if !ok {
				logger.Error(fmt.Sprintf("获取IPV4地址失败: [%s]", "未配置IPV4接口"))
				err = fmt.Errorf("获取IPV4地址失败")
			}else{
				tag := false
				for _, addr := range ipv4api {
					ipv4, err = ddns.IPfromAPI(addr, models.Rev4, client, logger)
					if err != nil || ipv4 == "" {
						logger.Error(fmt.Sprintf(`"%s"获取IPV4地址失败: [%s]`, addr, err.Error()))
						continue
					}
					tag = true
				}
				if !tag {
					err = fmt.Errorf("API获取IPV4地址失败")
				}
			}
		case models.SCRIPT:
			ipv4, err = ddns.IPfromScript(form.V4script, logger)
			if err != nil {
				logger.Error(fmt.Sprintf("获取IPV4地址失败: [%s]", err.Error()))
			}
	}
	if err != nil {
		for i, domain := range form.Domains {
			domain.Result = err.Error()
			domains[i] = domain
		}
		form.Domains = domains
		if err := entClient.DDNS.Create().SetConfig(form.Config).SetDomains(domains).SetWebhook(form.Webhook).
									SetV4method(form.V4method).SetV4interface(form.V4interface).
									SetV4script(form.V4script).SetRev4(form.Rev4).
									SetIpv4(ipv4).SetIpv6(ipv6).
									SetV6method(form.V6method).SetV6interface(form.V6interface).
									SetV4script(form.V6script).SetRev6(form.Rev6).
									Exec(context.Background()); err != nil{
			logger.Error(fmt.Sprintf("保存任务失败: [%s]", err.Error()))
			return form, fmt.Errorf("保存任务失败")
		}
		return form, fmt.Errorf("获取IPV4地址失败")
	}
	switch form.V6method {
		case models.INTERFACE:
			ipv6, err = ddns.IPfromInterface(form.V6interface, form.Rev6, logger)
			if err != nil {
				logger.Error(fmt.Sprintf("获取IPV6地址失败: [%s]", err.Error()))
			}
		case models.IPAPI:
			ipv6api, ok := ipAPI["ipv6"]
			if !ok {
				logger.Error(fmt.Sprintf("获取IPV6地址失败: [%s]", "未配置IPV6接口"))
				err = fmt.Errorf("获取IPV6地址失败")
			}else{
				tag := false
				for _, addr := range ipv6api {
					ipv6, err = ddns.IPfromAPI(addr, models.Rev6, client, logger)
					if err != nil || ipv6 == "" {
						logger.Error(fmt.Sprintf(`"%s"获取IPV6地址失败: [%s]`, addr, err.Error()))
						continue
					}
					tag = true
				}
				if !tag {
					err = fmt.Errorf("API获取IPV6地址失败")
				}
			}
		case models.SCRIPT:
			ipv6, err = ddns.IPfromScript(form.V6script, logger)
			if err != nil {
				logger.Error(fmt.Sprintf("获取IPV6地址失败: [%s]", err.Error()))
			}
	}
	if err != nil {
		for i, domain := range form.Domains {
			domain.Result = err.Error()
			domains[i] = domain
		}
		form.Domains = domains
		if err := entClient.DDNS.Create().SetConfig(form.Config).SetDomains(domains).SetWebhook(form.Webhook).
									SetV4method(form.V4method).SetV4interface(form.V4interface).
									SetV4script(form.V4script).SetRev4(form.Rev4).
									SetIpv4(ipv4).SetIpv6(ipv6).
									SetV6method(form.V6method).SetV6interface(form.V6interface).
									SetV4script(form.V6script).SetRev6(form.Rev6).
									Exec(context.Background()); err != nil{
			logger.Error(fmt.Sprintf("保存任务失败: [%s]", err.Error()))
			return form, fmt.Errorf("保存任务失败")
		}
		return form, fmt.Errorf("获取IPV6地址失败")
	}
	for i, domain := range form.Domains {
		if domain.Value == "" {
			switch domain.Type {
				case models.IPV4:
					domain.Value = ipv4
					domains[i] = domain
				case models.IPV6:
					domain.Value = ipv6
					domains[i] = domain
			}
		}
	}
	switch resolver {
		case models.CF:
			results, err := ddns.CloudFlare(api, form.Config[models.CFTOKEN], domains, client, logger)
			if err != nil {
				logger.Error(fmt.Sprintf("更新域名失败: [%s]", err.Error()))
				for i, domain := range results {
					domain.Result = fmt.Sprintf("更新域名失败: [%s]", err.Error())
					domains[i] = domain
				}
			}
			domains = make([]models.Domain, 0)
			for _, domain := range results {
				if domain.Status != models.DELETE {
					domains = append(domains, domain)
				}
			}
		default:
			logger.Error(fmt.Sprintf("暂不支持%s", resolver))
			return form, fmt.Errorf("暂不支持%s", resolver)
	}
	form.Domains = domains
	if err := entClient.DDNS.Create().SetConfig(form.Config).SetDomains(domains).SetWebhook(form.Webhook).
							SetV4method(form.V4method).SetV4interface(form.V4interface).
							SetV4script(form.V4script).SetRev4(form.Rev4).
							SetIpv4(ipv4).SetIpv6(ipv6).
							SetV6method(form.V6method).SetV6interface(form.V6interface).
							SetV4script(form.V6script).SetRev6(form.Rev6).
							Exec(context.Background()); err != nil{
			logger.Error(fmt.Sprintf("保存任务失败: [%s]", err.Error()))
			return form, fmt.Errorf("保存任务失败")
	}
	return form, nil
}

func EditJobs(form models.JobForm, resolver, api string, id int, ipAPI map[string][]string, entClient *ent.Client, logger *zap.Logger) (models.JobForm, error) {
	var ipv4 string
	var ipv6 string
	var err error
	record, err := entClient.DDNS.Query().Select(entddns.FieldIpv4,entddns.FieldIpv6).Where(entddns.IDEQ(id)).Only(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("数据库获取任务失败: [%s]", err.Error()))
		return form, fmt.Errorf("数据库获取任务失败")
	}
	client := http.DefaultClient
	domains := make([]models.Domain, len(form.Domains))
	switch form.V4method {
		case models.INTERFACE:
			if record.Ipv4 != ""{
				ipv4 = record.Ipv4
			}else{
				ipv4, err = ddns.IPfromInterface(form.V4interface, form.Rev4, logger)
				if err != nil {
					logger.Error(fmt.Sprintf("获取IPV4地址失败: [%s]", err.Error()))
				}
			}	
		case models.IPAPI:
			if record.Ipv4 != ""{
				ipv4 = record.Ipv4
			}else{
				ipv4api, ok := ipAPI["ipv4"]
				if !ok {
					logger.Error(fmt.Sprintf("获取IPV4地址失败: [%s]", "未配置IPV4接口"))
					err = fmt.Errorf("获取IPV4地址失败")
				}else{
					tag := false
					for _, addr := range ipv4api {
						ipv4, err = ddns.IPfromAPI(addr, models.Rev4, client, logger)
						if err != nil || ipv4 == "" {
							logger.Error(fmt.Sprintf(`"%s"获取IPV4地址失败: [%s]`, addr, err.Error()))
							continue
						}
						tag = true
					}
					if !tag {
						err = fmt.Errorf("API获取IPV4地址失败")
					}
				}
			}
		case models.SCRIPT:
			if record.Ipv4 != ""{
				ipv4 = record.Ipv4
			}else{
				ipv4, err = ddns.IPfromScript(form.V4script, logger)
				if err != nil {
					logger.Error(fmt.Sprintf("获取IPV4地址失败: [%s]", err.Error()))
				}
			}	
	}
	if err != nil {
		for i, domain := range form.Domains {
			domain.Result = err.Error()
			domains[i] = domain
		}
		form.Domains = domains
		if err := entClient.DDNS.UpdateOneID(id).SetConfig(form.Config).SetDomains(domains).SetWebhook(form.Webhook).
									SetV4method(form.V4method).SetV4interface(form.V4interface).
									SetV4script(form.V4script).SetRev4(form.Rev4).
									SetIpv4(ipv4).SetIpv6(ipv6).
									SetV6method(form.V6method).SetV6interface(form.V6interface).
									SetV4script(form.V6script).SetRev6(form.Rev6).
									Exec(context.Background()); err != nil{
			logger.Error(fmt.Sprintf("保存任务失败: [%s]", err.Error()))
			return form, fmt.Errorf("保存任务失败")
		}
		return form, fmt.Errorf("获取IPV4地址失败")
	}
	switch form.V6method {
		case models.INTERFACE:
			if record.Ipv6 != ""{
				ipv6 = record.Ipv6
			}else{
				ipv6, err = ddns.IPfromInterface(form.V6interface, form.Rev6, logger)
				if err != nil {
					logger.Error(fmt.Sprintf("获取IPV6地址失败: [%s]", err.Error()))
				}
			}	
		case models.IPAPI:
			if record.Ipv6 != ""{
				ipv6 = record.Ipv6
			}else{
				ipv6api, ok := ipAPI["ipv6"]
				if !ok {
					logger.Error(fmt.Sprintf("获取IPV6地址失败: [%s]", "未配置IPV6接口"))
					err = fmt.Errorf("获取IPV6地址失败")
				}else{
					tag := false
					for _, addr := range ipv6api {
						ipv6, err = ddns.IPfromAPI(addr, models.Rev6, client, logger)
						if err != nil || ipv6 == "" {
							logger.Error(fmt.Sprintf(`"%s"获取IPV6地址失败: [%s]`, addr, err.Error()))
							continue
						}
						tag = true
					}
					if !tag {
						err = fmt.Errorf("API获取IPV6地址失败")
					}
				}
			}
		case models.SCRIPT:
			if record.Ipv6 != ""{
				ipv6 = record.Ipv6
			}else{
				ipv6, err = ddns.IPfromScript(form.V6script, logger)
				if err != nil {
					logger.Error(fmt.Sprintf("获取IPV6地址失败: [%s]", err.Error()))
				}
			}	
	}
	if err != nil {
		for i, domain := range form.Domains {
			domain.Result = err.Error()
			domains[i] = domain
		}
		form.Domains = domains
		if err := entClient.DDNS.UpdateOneID(id).SetConfig(form.Config).SetDomains(domains).SetWebhook(form.Webhook).
									SetV4method(form.V4method).SetV4interface(form.V4interface).
									SetV4script(form.V4script).SetRev4(form.Rev4).
									SetIpv4(ipv4).SetIpv6(ipv6).
									SetV6method(form.V6method).SetV6interface(form.V6interface).
									SetV4script(form.V6script).SetRev6(form.Rev6).
									Exec(context.Background()); err != nil{
			logger.Error(fmt.Sprintf("保存任务失败: [%s]", err.Error()))
			return form, fmt.Errorf("保存任务失败")
		}
		return form, fmt.Errorf("获取IPV6地址失败")
	}
	for i, domain := range form.Domains {
		if domain.Value == "" {
			switch domain.Type {
				case models.IPV4:
					domain.Value = ipv4
					domains[i] = domain
				case models.IPV6:
					domain.Value = ipv6
					domains[i] = domain
			}
		}
	}
	switch resolver {
		case models.CF:
			results, err := ddns.CloudFlare(api, form.Config[models.CFTOKEN], domains, client, logger)
			if err != nil {
				logger.Error(fmt.Sprintf("更新域名失败: [%s]", err.Error()))
				for i, domain := range results {
					domain.Result = fmt.Sprintf("更新域名失败: [%s]", err.Error())
					domains[i] = domain
				}
			}
			domains = make([]models.Domain, 0)
			for _, domain := range results {
				if domain.Status != models.DELETE {
					domains = append(domains, domain)
				}
			}
		default:
			logger.Error(fmt.Sprintf("暂不支持%s", resolver))
			return form, fmt.Errorf("暂不支持%s", resolver)
	}
	form.Domains = domains
	if err := entClient.DDNS.UpdateOneID(id).SetConfig(form.Config).SetDomains(domains).SetWebhook(form.Webhook).
							SetV4method(form.V4method).SetV4interface(form.V4interface).
							SetV4script(form.V4script).SetRev4(form.Rev4).
							SetIpv4(ipv4).SetIpv6(ipv6).
							SetV6method(form.V6method).SetV6interface(form.V6interface).
							SetV4script(form.V6script).SetRev6(form.Rev6).
							Exec(context.Background()); err != nil{
			logger.Error(fmt.Sprintf("保存任务失败: [%s]", err.Error()))
			return form, fmt.Errorf("保存任务失败")
	}
	return form, nil
}

func DeleteJobs(id int, entClient *ent.Client, logger *zap.Logger) error {
	if err := entClient.DDNS.DeleteOneID(id).Exec(context.Background()); err != nil {
		logger.Error(fmt.Sprintf("删除任务失败: [%s]", err.Error()))
		return fmt.Errorf("删除任务失败")
	}
	return nil
}

func GetJobs(entClient *ent.Client, logger *zap.Logger) ([]struct{
																ID	int	`json:"id"`
																Job	models.JobForm	`json:"job"`
														}, error) {
	ddnsJobs, err := entClient.DDNS.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取任务失败: [%s]", err.Error()))
		return nil, fmt.Errorf("获取任务失败")
	}
	var jobs = make([]struct{
		ID	int	`json:"id"`
		Job	models.JobForm	`json:"job"`
	}, len(ddnsJobs))
	for i, ddnsJob := range ddnsJobs {
		jobs[i] = struct{ID int `json:"id"`; Job models.JobForm `json:"job"`}{
			ID: ddnsJob.ID,
			Job: models.JobForm{
				Config: ddnsJob.Config,
				Domains: ddnsJob.Domains,
				Webhook: ddnsJob.Webhook,
				IPV4: ddnsJob.Ipv4,
				IPV6: ddnsJob.Ipv6,
				Rev4: ddnsJob.Rev4,
				Rev6: ddnsJob.Rev6,
				V4interface: ddnsJob.V4interface,
				V6interface: ddnsJob.V6interface,
				V4method: ddnsJob.V4method,
				V6method: ddnsJob.V6method,
				V4script: ddnsJob.V4script,
				V6script: ddnsJob.V6script,
			},
		}
	}
	return jobs, nil
}

func TestScript(script string, logger *zap.Logger) (string, error){
	result, err := ddns.IPfromScript(script, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("测试脚本失败: [%s]", err.Error()))
		return result, fmt.Errorf("测试脚本出错")
	}
	return result, nil
}