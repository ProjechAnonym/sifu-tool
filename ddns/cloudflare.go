package ddns

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func getZoneID(api, token string, client *http.Client, logger *zap.Logger) (string, error) {
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("创建请求失败: [%s]", err.Error()))
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
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
	results := make(map[string]any)
	if err := json.Unmarshal(content, &results); err != nil{
		logger.Error(fmt.Sprintf("解析响应失败: [%s]", err.Error()))
		return "", err
	}
	
	result, ok := results["result"].([]any)
	if !ok {
		logger.Error(`不存在"result"字段`)
		return "", fmt.Errorf(`不存在"result"字段`)
	}
	if len(result) != 0 {
		if id, ok := result[0].(map[string]any)["id"].(string); !ok {
			logger.Error(`不存在"id"字段`)
			return "", fmt.Errorf(`不存在"id"字段`)
		}else{
			return id, nil
		}
	}
	return "", fmt.Errorf("获取zoneID失败")
}
func getRecordID(zoneID, api, token string, client *http.Client, logger *zap.Logger) {
	req, err := http.NewRequest("GET", fmt.Sprintf(api, zoneID), nil)
	if err != nil {
		logger.Error(fmt.Sprintf("创建请求失败: [%s]", err.Error()))
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("请求失败: [%s]", err.Error()))

	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("读取响应失败: [%s]", err.Error()))
	}
	results := make(map[string]any)
	if err := json.Unmarshal(content, &results); err != nil{
		logger.Error(fmt.Sprintf("解析响应失败: [%s]", err.Error()))
		return 
	}
	
	result, ok := results["result"].([]any)
	if !ok {
		logger.Error(`不存在"result"字段`)
		return 
	}
	domainsMsg := make(map[string]map[string]string)
	if len(result) != 0 {
		for _, domainMsg := range result {
			name, ok := domainMsg.(map[string]any)["name"].(string)
			if !ok {
				logger.Error(`不存在"name"字段`)
				return
			}

			msg :=  make(map[string]string)
			if id, ok := domainMsg.(map[string]any)["id"].(string); !ok {
				logger.Error(`不存在"id"字段`)
				return
			}else{
				msg["id"] = id
			}

			if domainType, ok := domainMsg.(map[string]any)["type"].(string); !ok {
				logger.Error(`不存在"type"字段`)
				return
			}else{
				msg["type"] = domainType
			}

			if ip, ok := domainMsg.(map[string]any)["ip"].(string); !ok {
				logger.Error(`不存在"id"字段`)
				return
			}else{
				msg["ip"] = ip
			}
			domainsMsg[name] = msg
		}
	}
}
func SetCFRecord(zoneAPI, recordAPI, token string, client *http.Client, logger *zap.Logger)  {
	a,e:=getZoneID(zoneAPI, token, client, logger)
	fmt.Println(a,e)
}