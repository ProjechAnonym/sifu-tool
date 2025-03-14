package ddns

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"sifu-tool/models"

	"go.uber.org/zap"
)

func IPfromAPI(api string, client *http.Client, logger *zap.Logger) (string, error) {
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("创建请求失败: [%s]", err.Error()))
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("请求失败: [%s]", err.Error()))
		return "", err
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("读取响应失败: [%s]", err.Error()))
		return "", err
	}
	
	rev4 := regexp.MustCompile(models.Rev4)
	rev6 := regexp.MustCompile(models.Rev6)
	if ip := rev4.FindString(string(content)); ip != "" {
		return fmt.Sprintf("v4:[%s]",ip), nil
	}
	if ip := rev6.FindString(string(content)); ip != "" {
		return fmt.Sprintf("v6:[%s]",ip), nil
	}
	return "", fmt.Errorf("未查找到IP字段")
}

func IPfromInterface(interfaceName, reStr string, logger *zap.Logger) (string, error) {
	targetInterface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		logger.Error(fmt.Sprintf(`获取网卡"%s"失败: [%s]`, interfaceName, err.Error()))
		return "", fmt.Errorf(`获取网卡"%s"失败`, interfaceName)
	}
	addrs, err := targetInterface.Addrs()
	if err != nil {
		logger.Error(fmt.Sprintf("获取网卡地址失败: [%s]", err.Error()))
		return "", fmt.Errorf("获取网卡地址失败")
	}
	for _, addr := range addrs {
		reAddr := regexp.MustCompile(reStr)
		ip := reAddr.FindString(addr.String())
		if ip == "" {
			continue
		}
		address, _, err := net.ParseCIDR(ip)
		if err != nil {
			logger.Error(fmt.Sprintf(`解析网卡"%s"地址"%s"失败: [%s]`, targetInterface.Name, addr.String(), err.Error()))
			return "", fmt.Errorf(`解析网卡"%s"地址"%s"失败`, targetInterface.Name, addr.String())
		}
		return address.String(), nil
	}
	return "", fmt.Errorf("未查找到IP字段")
}