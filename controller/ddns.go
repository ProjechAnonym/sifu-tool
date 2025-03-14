package controller

import (
	"fmt"
	"net"

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