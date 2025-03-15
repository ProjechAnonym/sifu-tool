package ddns

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sifu-tool/models"

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
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("读取响应失败: [%s]", err.Error()))
		return "", err
	}
	content := make(map[string]any)
	if err := json.Unmarshal(respBody, &content); err != nil{
		logger.Error(fmt.Sprintf("解析响应失败: [%s]", err.Error()))
		return "", err
	}
	
	result, ok := content["result"].([]any)
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
func getRecordID(zoneID, api, token string, client *http.Client, logger *zap.Logger) (map[string]map[string]string, map[string][]string, error){
	req, err := http.NewRequest("GET", fmt.Sprintf(api, zoneID), nil)
	if err != nil {
		logger.Error(fmt.Sprintf("创建请求失败: [%s]", err.Error()))
		return nil, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("请求失败: [%s]", err.Error()))
		return nil, nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("读取响应失败: [%s]", err.Error()))
		return nil, nil, err
	}
	content := make(map[string]any)
	if err := json.Unmarshal(respBody, &content); err != nil{
		logger.Error(fmt.Sprintf("解析响应失败: [%s]", err.Error()))
		return nil, nil, err
	}
	
	results, ok := content["result"].([]any)
	if !ok {
		logger.Error(`不存在"result"字段`)
		return nil, nil, fmt.Errorf(`不存在"result"字段`)
	}
	records := make(map[string]map[string]string)
	domains := make(map[string][]string)
	if len(results) != 0 {
		for _, result := range results {
			id, ok := result.(map[string]any)["id"].(string)
			if !ok {
				logger.Error(`不存在"id"字段`)
				continue
			}

			msg := map[string]string{}

			if recordType, ok := result.(map[string]any)["type"].(string); !ok {
				logger.Error(`不存在"type"字段`)
				continue
			}else{
				msg["type"] = recordType
			}

			if ip, ok := result.(map[string]any)["content"].(string); !ok {
				logger.Error(`不存在"id"字段`)
				continue
			}else{
				msg["ip"] = ip
			}
			
			if name, ok := result.(map[string]any)["name"].(string); !ok {
				logger.Error(`不存在"name"字段`)
				continue
			}else{
				msg["name"] = name
				if domains[name] == nil {
					domains[name] = []string{id}
				}else{
					domains[name] = append(domains[name], id)
				}
			}
			records[id] = msg
		}
		return records, domains, nil
	}
	return nil, nil, fmt.Errorf("获取recordID失败")
}
func SetCFRecord(zoneAPI, recordAPI, token, ipv4, ipv6 string, domains []models.Domain, client *http.Client, logger *zap.Logger) error {
	zoneID, err := getZoneID(zoneAPI, token, client, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取zoneID失败: [%s]", err.Error()))
		return err
	}
	records, recordsMap, err := getRecordID(zoneID, recordAPI, token, client, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取recordID失败: [%s]", err.Error()))
		return err
	}
	for _, domain := range domains {
		if recordsMap[domain.Domain] == nil {
			continue
		}
		for _, recordID := range recordsMap[domain.Domain] {
			if records[recordID]["type"] == domain.Type {
				
			}
		}

	}

	return nil
}