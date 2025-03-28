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

// GetInterfaceIPs 获取所有网络接口的IP地址
// 该函数接受一个 *zap.Logger 作为日志记录器, 用于记录错误信息
// 返回值是一个映射, 键是接口名称, 值是该接口的IP地址列表；
// 第二个返回值是一个错误列表, 如果在获取或解析IP地址时遇到任何错误, 将被添加到此列表中
func GetInterfaceIPs(logger *zap.Logger) (map[string][]string, []error){
    // 获取所有网络接口信息
    interfaces, err := net.Interfaces()
    if err != nil {
        logger.Error(fmt.Sprintf("获取网卡信息失败: [%s]", err.Error()))
        return nil, []error{fmt.Errorf("获取网卡信息失败")}
    }
    
    interfaceIPs := make(map[string][]string)
    var errors []error
    
    // 遍历每个网络接口
    for _, netInterface := range interfaces {
        // 获取当前接口的所有地址
        addrs, err := netInterface.Addrs()
        if err != nil {
            logger.Error(fmt.Sprintf("获取网卡地址失败: [%s]", err.Error()))
            errors = append(errors, fmt.Errorf(`获取网卡"%s"IP地址失败`, netInterface.Name))
        }
        
        addresses := make([]string, len(addrs))
        
        // 遍历当前接口的所有地址
        for i, addr := range addrs {
            // 解析CIDR地址, 获取IP
            ip, _, err := net.ParseCIDR(addr.String())
            if err != nil {
                logger.Error(fmt.Sprintf(`解析网卡"%s"地址"%s"失败: [%s]`, netInterface.Name, addr.String(), err.Error()))
                errors = append(errors, fmt.Errorf(`解析网卡"%s"地址"%s"失败`, netInterface.Name, addr.String()))
            }
            addresses[i] = ip.String()
        }
        
        // 将接口名称与其对应的IP地址列表存入映射
        interfaceIPs[netInterface.Name] = addresses
    }
    
    // 返回所有接口的IP地址映射和遇到的错误列表
    return interfaceIPs, errors
}

// AddJobs 添加DDNS任务
// 该函数根据JobForm中的配置, 获取IPv4和IPv6地址, 并更新到指定的域名中
// 参数:
//   form: 包含任务配置的JobForm对象
//   resolver: 域名解析器类型
//   api: 解析器的API密钥
//   ipAPI: 用于获取IP地址的API映射
//   entClient: 数据库客户端
//   logger: 日志记录器
// 返回值:
//   更新后的JobForm对象和可能的错误
func AddJob(form models.JobForm, resolver string, ipAPI map[string][]string, config map[string]map[string]any, entClient *ent.Client, logger *zap.Logger) (models.JobForm, error) {
    // 初始化HTTP客户端
    client := http.DefaultClient

    // 初始化域名数组, 用于存储处理后的域名信息
    domains := make([]models.Domain, len(form.Domains))
	
	// 获取IPv4和IPv6API接口
	ipv4api, ok := ipAPI["ipv4"]
	if !ok {
		logger.Error("未配置IPV4接口")
		return form, fmt.Errorf("未配置IPV4接口")
	}
	ipv6api, ok := ipAPI["ipv6"]
	if !ok {
		logger.Error("未配置IPV6接口")
		return form, fmt.Errorf("未配置IPV6接口")
	}

	// 按照不同方法获取IPv4和IPv6地址
	ipv4, err := ddns.GetIP(form.V4method, form.V4interface, form.Rev4, models.Rev4, form.V4script, client, logger, ipv4api...)
	if err != nil {
		logger.Error(fmt.Sprintf("获取IP地址失败: [%s]", err.Error()))
	}

	ipv6, err := ddns.GetIP(form.V6method, form.V6interface, form.Rev6, models.Rev6, form.V6script, client, logger, ipv6api...)
    if err != nil {
		logger.Error(fmt.Sprintf("获取IP地址失败: [%s]", err.Error()))
	}

    // 设置任务域名的要更新的IP地址
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
        }else{
			domains[i] = domain
		}
    }
    // 根据托管商更新域名
    switch resolver {
		case models.CF:
			// 根据不同的DNS服务商断言不同的配置类型以获取所需信息
			cloudflare, ok := config[resolver]["api"].(string)
			if !ok {
				logger.Error(`未能读取"cloudflare"接口信息`)
				return form, fmt.Errorf(`未能读取"cloudflare"接口信息`)
			}
			// 更新域名
			results, err := ddns.CloudFlare(cloudflare, form.Config[models.CFTOKEN], domains, client, logger)
			// 该函数返回的错误为初始化错误，若出错则整个域名列表都不会更新
			// 因此将每条域名的结果都设置为该错误
			if err != nil {
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
		default:
			logger.Error(fmt.Sprintf("暂不支持%s", resolver))
			return form, fmt.Errorf("暂不支持%s", resolver)
    }

    // 更新任务的域名信息
    form.Domains = domains
    if err := entClient.DDNS.Create().SetConfig(form.Config).SetDomains(domains).SetWebhook(form.Webhook).
							SetV4method(form.V4method).SetV4interface(form.V4interface).
							SetV4script(form.V4script).SetRev4(form.Rev4).
							SetIpv4(ipv4).SetIpv6(ipv6).
							SetV6method(form.V6method).SetV6interface(form.V6interface).
							SetV4script(form.V6script).SetRev6(form.Rev6).
							Exec(context.Background()); err != nil {
        logger.Error(fmt.Sprintf("保存任务失败: [%s]", err.Error()))
        return form, fmt.Errorf("保存任务失败")
    }

    // 返回更新后的任务信息
    return form, nil
}


// EditJobs 编辑作业以更新DNS记录
// 该函数根据JobForm中的配置获取IPv4和IPv6地址, 并更新域名解析记录
// 参数:
// - form: 包含作业配置和域名信息的JobForm对象
// - resolver: 解析器类型, 用于确定使用哪个DNS服务提供商的API来更新域名记录
// - api: API密钥, 用于访问DNS服务提供商的API
// - id: 作业的数据库ID
// - ipAPI: 包含IPv4和IPv6API地址的映射, 用于获取IP地址
// - entClient: 数据库客户端, 用于与数据库交互
// - logger: 日志记录器, 用于记录日志信息
// 返回值:
// - models.JobForm: 更新后的JobForm对象
// - error: 如果操作过程中发生错误, 则返回错误
func EditJob(form models.JobForm, resolver string, id int, ipAPI map[string][]string, config map[string]map[string]any, entClient *ent.Client, logger *zap.Logger) (models.JobForm, error) {
	var (
		ipv4 string
		ipv6 string
		err error
	)
    // 从数据库中获取当前任务的IPv4和IPv6地址
    record, err := entClient.DDNS.Query().Select(entddns.FieldIpv4,entddns.FieldIpv6).Where(entddns.IDEQ(id)).Only(context.Background())
    if err != nil {
        logger.Error(fmt.Sprintf("数据库获取任务失败: [%s]", err.Error()))
        return form, fmt.Errorf("数据库获取任务失败")
    }

    client := http.DefaultClient
    domains := make([]models.Domain, len(form.Domains))

	// 获取IPv4和IPv6API接口
	ipv4api, ok := ipAPI["ipv4"]
	if !ok {
		logger.Error("未配置IPV4接口")
		return form, fmt.Errorf("未配置IPV4接口")
	}
	ipv6api, ok := ipAPI["ipv6"]
	if !ok {
		logger.Error("未配置IPV6接口")
		return form, fmt.Errorf("未配置IPV6接口")
	}

	// 获取IPv4和IPv6地址
	// 若记录中已经存在IP地址则直接使用记录中的地址
	// 否则按照不同方法获取IP地址
	if record.Ipv4 != ""{
		ipv4 = record.Ipv4
	}else{
		ipv4, err = ddns.GetIP(form.V4method, form.V4interface, form.Rev4, models.Rev4, form.V4script, client, logger, ipv4api...)
		if err != nil {
			logger.Error(fmt.Sprintf("获取IP地址失败: [%s]", err.Error()))
		}
	}

	if record.Ipv6 != ""{
		ipv6 = record.Ipv6
	}else{
		ipv6, err = ddns.GetIP(form.V6method, form.V6interface, form.Rev6, models.Rev6, form.V6script, client, logger, ipv6api...)
		if err != nil {
			logger.Error(fmt.Sprintf("获取IP地址失败: [%s]", err.Error()))
		}
	}

    // 设置任务域名的要更新的IP地址
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
        }else{
			domains[i] = domain
		}
    }

    // 根据托管商类型更新域名记录
    switch resolver {
		case models.CF:
			// 根据不同的DNS服务商断言不同的配置类型以获取所需信息
			cloudflare, ok := config[resolver]["api"].(string)
			if !ok {
				logger.Error(`未能读取"cloudflare"接口信息`)
				return form, fmt.Errorf(`未能读取"cloudflare"接口信息`)
			}

			// 更新域名
			results, err := ddns.CloudFlare(cloudflare, form.Config[models.CFTOKEN], domains, client, logger)
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
		default:
			logger.Error(fmt.Sprintf("暂不支持%s", resolver))
			return form, fmt.Errorf("暂不支持%s", resolver)
	}

    // 更新任务的域名列表
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