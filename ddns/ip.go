package ddns

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"sifu-tool/models"
	"strings"

	"github.com/go-cmd/cmd"
	"go.uber.org/zap"
)

func IPfromAPI(api, reStr string, client *http.Client, logger *zap.Logger) (string, error) {
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
	re := regexp.MustCompile(reStr)
	if ip := re.FindString(string(content)); ip != "" {
		return ip, nil
	}
	return "", fmt.Errorf("未查找到IP字段")
}

func IPfromInterface(interfaceName, reStr string, logger *zap.Logger) (string, error) {
	targetInterface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		logger.Error(fmt.Sprintf(`获取网卡"%s"失败: [%s]`, interfaceName, err.Error()))
		return "", fmt.Errorf(`获取网卡"%s"失败: [%s]`, interfaceName, err.Error())
	}
	addrs, err := targetInterface.Addrs()
	if err != nil {
		logger.Error(fmt.Sprintf("获取网卡地址失败: [%s]", err.Error()))
		return "", fmt.Errorf("获取网卡地址失败: [%s]", err.Error())
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
			return "", fmt.Errorf(`解析网卡"%s"地址"%s"失败: [%s]`, targetInterface.Name, addr.String(), err.Error())
		}
		return address.String(), nil
	}
	return "", fmt.Errorf("未查找到IP字段")
}

func IPfromScript(script string, logger *zap.Logger) (string, error) {
	cmd := cmd.NewCmd("sh", "-c", script)
	status := <- cmd.Start()
	// 执行命令
	if status.Error != nil {
		logger.Error(fmt.Sprintf("执行脚本失败: [%s]", status.Error.Error()))
		return "", status.Error
	}
	if len(status.Stderr) != 0 {
		logger.Error(fmt.Sprintf("执行脚本失败: [%s]", strings.Join(status.Stderr, "")))
		return "", errors.New(strings.Join(status.Stderr, ""))
	}
	
	rev4 := regexp.MustCompile(models.Rev4)
	rev6 := regexp.MustCompile(models.Rev6)
	if ip := rev4.FindString(strings.Join(status.Stdout, "")); ip != "" {
		return ip, nil
	}
	if ip := rev6.FindString(strings.Join(status.Stdout, "")); ip != "" {
		return ip, nil
	}
	return "", fmt.Errorf("未查找到IP字段")
}

func GetIP(method int, netInterface, interfaceRe, apiRe, script string, client *http.Client, logger *zap.Logger, api ...string) (string, error) {
	// 根据不同方法获取IP地址
	switch method {
		// 通过网络接口获取IP地址
		case models.INTERFACE:
			ip, err := IPfromInterface(netInterface, interfaceRe, logger)
			// 如果获取IP地址失败, 则记录错误
			if err != nil {
				logger.Error(fmt.Sprintf("获取IP地址失败: [%s]", err.Error()))
				return "", err
			}
			return ip, nil

		// 通过IP接口API获取IP地址
		case models.IPAPI:

			// 如果未正确接收API列表, 则记录错误
			if len(api) == 0 {
				logger.Error(fmt.Sprintf("获取IP地址失败: [%s]", "未配置IP接口"))
				return "", fmt.Errorf("未配置IP接口")
			}else{
				// 通过IP接口API获取IP地址, 如果错误则使用下一个接口
				// 成功后记录获取的地址并终止循环, 并将错误置空
				// 若所有接口都不能正确获取地址, 则记录错误
				for _, addr := range api {
					ip, err := IPfromAPI(addr, apiRe, client, logger)
					if err != nil || ip == "" {
						logger.Error(fmt.Sprintf(`"%s"获取IP地址失败: [%s]`, addr, err.Error()))
						continue
					}
					return ip, nil
				}
				return "", fmt.Errorf("不能正确获取IP地址")
			}
		// 通过脚本获取IP地址
		case models.SCRIPT:
			ip, err := IPfromScript(script, logger)
			// 如果获取IP地址失败, 则记录错误
			if err != nil {
				logger.Error(fmt.Sprintf("获取IP地址失败: [%s]", err.Error()))
				return "", err
			}
			return ip, nil
	}
	return "", nil
}