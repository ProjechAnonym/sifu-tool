package ddns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sifu-tool/models"
	"sync"

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
func getRecordID(zoneID, api, token string, client *http.Client, logger *zap.Logger) (map[string]*models.Domain, map[string][]string, error){
	api += fmt.Sprintf("/%s/dns_records", zoneID)
	req, err := http.NewRequest("GET", api, nil)
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
	records := make(map[string][]string)
	domains := make(map[string]*models.Domain)
	if len(results) != 0 {
		for _, result := range results {
			id, ok := result.(map[string]any)[CFID].(string)
			if !ok {
				logger.Error(`不存在"id"字段`)
				continue
			}

			domain := models.Domain{}

			if recordType, ok := result.(map[string]any)[CFTYPE].(string); !ok {
				logger.Error(`不存在"type"字段`)
				continue
			}else{
				domain.Type = recordType
			}

			if content, ok := result.(map[string]any)[CFVALUE].(string); !ok {
				logger.Error(`不存在"content"字段`)
				continue
			}else{
				domain.Value = content
			}
			
			if name, ok := result.(map[string]any)["name"].(string); !ok {
				logger.Error(`不存在"name"字段`)
				continue
			}else{
				domain.Domain = name
				if records[name] == nil {
					records[name] = []string{id}
				}else{
					records[name] = append(records[name], id)
				}
			}
			domains[id] = &domain
		}
		return domains, records, nil
	}
	return nil, nil, nil
}
func SetCFRecord(api, token string, value map[string]string, editDomains []models.Domain, client *http.Client, logger *zap.Logger) ([]*models.Domain, error) {
	zoneID, err := getZoneID(api, token, client, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取zoneID失败: [%s]", err.Error()))
		return nil, err
	}
	domains, records, err := getRecordID(zoneID, api, token, client, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取recordID失败: [%s]", err.Error()))
		return nil, err
	}
	addList := make([]models.Domain, 0)
	editList := make([]struct {
		domain models.Domain 
		id     string
	}, 0)
	deleteList := make([]struct {
		domain models.Domain 
		id     string
	}, 0)
	for _, domain := range editDomains {
		if records[domain.Domain] == nil {
			addList = append(addList, domain)
			continue
		}
		edit := false
		for _, recordID := range records[domain.Domain] {
			if domains[recordID].Type == domain.Type {
				editList = append(editList, struct{domain models.Domain; id string}{domain: domain, id: recordID})
				edit = true
				break
			}
		}
		if !edit {
			addList = append(addList, domain)
		}
	}
	for id, item := range domains {
		deleteTag := true
		for _, domain := range editDomains {
			if item.Domain == domain.Domain && item.Type == domain.Type {
				deleteTag = false
				break
			}
		}
		if deleteTag {
			deleteList = append(deleteList, struct{domain models.Domain; id string}{domain: *domains[id], id: id})
		}
	}
	var jobs sync.WaitGroup
	itemChannel := make(chan *models.Domain, 10)
	countChannel := make(chan int, 10)
	domainList := make([]*models.Domain, 0)
	jobs.Add(1)
	go func() {
		defer func(){
            jobs.Done()
            var ok bool
            if _, ok = <- countChannel; ok {close(countChannel)}
			if _, ok = <- itemChannel; ok {close(itemChannel)}
        }()
		sum := 0
		for {
			if sum == len(addList) + len(editList) + len(deleteList) {
				return
			}
			select{
				case count, ok := <- countChannel:
					if !ok {return}
					sum += count
				case item, ok := <- itemChannel:
					if !ok {return}
					domainList = append(domainList, item)
			}
		}
	}()
	for _, domain := range addList {
		jobs.Add(1)
		go func() {
			defer func() {
				countChannel <- 1
				jobs.Done()
			}()
			item, err := setCFRecord(zoneID, "", api, token, value[domain.Type], "add", domain, client, logger)
			if err != nil {
				logger.Error(fmt.Sprintf("添加记录失败: [%s]", err.Error()))
				item.Result = err.Error()
			}
			itemChannel <- item
		}()
	}
	for _, editItem := range editList {
		jobs.Add(1)
		go func() {
			defer func() {
				countChannel <- 1
				jobs.Done()
			}()
			item, err := setCFRecord(zoneID, editItem.id, api, token, value[editItem.domain.Type], "edit", editItem.domain, client, logger)
			if err != nil {
				logger.Error(fmt.Sprintf("修改记录失败: [%s]", err.Error()))
				item.Result = err.Error()
			}
			itemChannel <- item
		}()
	}
	for _, deleteItem := range deleteList {
		jobs.Add(1)
		go func() {
			defer func() {
				countChannel <- 1
				jobs.Done()
			}()
			item, err := setCFRecord(zoneID, deleteItem.id, api, token, value[deleteItem.domain.Type], "delete", deleteItem.domain, client, logger)
			if err != nil {
				logger.Error(fmt.Sprintf("删除记录失败: [%s]", err.Error()))
				item.Result = err.Error()
			}
			itemChannel <- item
		}()
	}
	jobs.Wait()
	
	return domainList, nil
}

func setCFRecord(zoneID, recordID, api, token, value, operation string, domain models.Domain, client *http.Client, logger *zap.Logger) (*models.Domain, error) {
	if domain.Value == value && operation == "edit" {
		domain.Status = models.STATIC
		domain.Result = "与托管商记录一致"
		return &domain, nil
	}
	url := api + fmt.Sprintf("/%s/dns_records", zoneID)
	var req *http.Request
	var err error
	var result string
	switch operation {
	case "add":
		result = "添加域名记录"
		content, err := json.Marshal(map[string]string{CFNAME: domain.Domain, CFTYPE: domain.Type, CFVALUE: value})
		if err != nil {
			logger.Error(fmt.Sprintf("解析数据失败: [%s]", err.Error()))
			return &domain, err
		}
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(content))
		if err != nil {
			logger.Error(fmt.Sprintf("创建请求失败: [%s]", err.Error()))
			return &domain, err
		}
	case "delete":
		result = "删除域名记录"
		url += fmt.Sprintf("/%s", recordID)
		req, err = http.NewRequest("DELETE", url, nil)
		if err != nil {
			logger.Error(fmt.Sprintf("创建请求失败: [%s]", err.Error()))
			return &domain, err
		}
	case "edit":
		result = "修改域名记录"
		url += fmt.Sprintf("/%s", recordID)
		content, err := json.Marshal(map[string]string{CFNAME: domain.Domain, CFTYPE: domain.Type, CFVALUE: value})
		if err != nil {
			logger.Error(fmt.Sprintf("解析数据失败: [%s]", err.Error()))
			return &domain, err
		}
		req, err = http.NewRequest("PUT", url, bytes.NewBuffer(content))
		if err != nil {
			logger.Error(fmt.Sprintf("创建请求失败: [%s]", err.Error()))
			return &domain, err
		}
		
	default:
		logger.Error("操作类型错误")
		return &domain, fmt.Errorf("操作类型错误")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("请求失败: [%s]", err.Error()))
		return &domain, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("读取响应失败: [%s]", err.Error()))
		return &domain, err
	}
	content := make(map[string]any)
	if err := json.Unmarshal(respBody, &content); err != nil{
		logger.Error(fmt.Sprintf("解析响应失败: [%s]", err.Error()))
		return &domain, err
	}
	if status, ok := content["success"].(bool); ok {
		if status {
			return &models.Domain{Domain: domain.Domain, Result: fmt.Sprintf("%s成功", result), Status: models.SUCCESS, Type: domain.Type, Value: value}, nil
		}else{
			if errors, ok := content["errors"].([]map[string]any); ok {
				for _, err := range errors {
					if msg, ok := err["message"].(string); ok {
						return &models.Domain{Domain: domain.Domain, Result: msg, Status: models.FAILURE, Type: domain.Type, Value: domain.Value}, fmt.Errorf("%s", msg)
					}else{
						return &models.Domain{Domain: domain.Domain, Result: fmt.Sprintf(`%s失败, 不存在"message"字段`, result), Status: models.FAILURE, Type: domain.Type, Value: domain.Value}, fmt.Errorf(`不存在"message"字段`)
					}
				}
			}else{
				return &models.Domain{Domain: domain.Domain, Result: fmt.Sprintf(`%s失败, 不存在"errors"字段`, result), Status: models.FAILURE, Type: domain.Type, Value: domain.Value}, fmt.Errorf(`不存在"errors"字段`)
			}
		}
	}
	return &models.Domain{Domain: domain.Domain, Result: fmt.Sprintf(`%s失败, 不存在"success"字段`, result), Status: models.FAILURE, Type: domain.Type, Value: domain.Value}, fmt.Errorf(`不存在"success"字段`)
}