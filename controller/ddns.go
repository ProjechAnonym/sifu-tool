package controller

import (
	"fmt"
	"net"
	"net/http"
	"sifu-tool/ddns"
	"sifu-tool/ent"
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

func AddJobs(form models.JobForm, resolver, api string, entClient *ent.Client, logger *zap.Logger) (models.JobForm, error) {
	var ipv4 string
	var ipv6 string
	var err error
	client := http.DefaultClient
	switch form.V4method {
		case models.INTERFACE:
			ipv4 = form.IPV4
		case models.IPAPI:
			ipv4, err = ddns.IPfromAPI(form.V4interface, models.Rev4, client, logger)
			if err != nil {
				logger.Error(fmt.Sprintf("API获取IPV4地址失败: [%s]", err.Error()))
				entClient.DDNS.Create().SetDomains(form.Domains).SetConfig(form.Config).SetV4method(form.V4method)
				return form, fmt.Errorf("API获取IPV4地址失败")
			}
		case models.SCRIPT:
			ipv4 = form.IPV4
	}
	switch form.V6method {
		case models.INTERFACE:
			ipv6 = form.IPV6
		case models.IPAPI:
			ipv6, err = ddns.IPfromAPI(form.V6interface, models.Rev6, client, logger)
			if err != nil {
				logger.Error(fmt.Sprintf("获取网卡IPV6地址失败: [%s]", err.Error()))
				return form, fmt.Errorf("API获取IPV6地址失败")
			}
		case models.SCRIPT:
			ipv6 = form.IPV6
	}
	domains := make([]models.Domain, len(form.Domains))
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
				return form, fmt.Errorf("更新域名失败")
			}
			form.Domains = results
		
		default:
			logger.Error(fmt.Sprintf("暂不支持%s", resolver))
			return form, fmt.Errorf("暂不支持%s", resolver)
	}
	
	return form, nil
}