package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

// 获取 open 接口的 data 字段内容
func GetMaxKBOpenDataService(s *string) {
	url := "http://154.39.79.242:8080/chat/api/open"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("创建maxkb open请求失败:", err)
		return
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("Authorization", "Bearer application-8a25368429d8c9414eed9f269622b276")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("maxkb open请求失败:", err)
		return
	}
	defer resp.Body.Close()

	var result struct {
		Data string `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("解析maxkb open响应失败:", err)
		return
	}
	*s = result.Data
}

// AI控制业务逻辑
func AiControlService(c *gin.Context, req interface{}) {
	// 这里 req 需断言为原结构体类型
	reqStruct, ok := req.(struct {
		Question string `json:"question"`
	})
	if !ok {
		c.JSON(400, gin.H{"error": "参数类型错误"})
		return
	}

	var chat_id string
	GetMaxKBOpenDataService(&chat_id)

	postBody := struct {
		ChatID  string `json:"chat_id"`
		Message string `json:"message"`
		Stream  bool   `json:"stream"`
		ReChat  bool   `json:"re_chat"`
	}{
		ChatID:  chat_id,
		Message: reqStruct.Question,
		Stream:  false,
		ReChat:  true,
	}
	bodyBytes, _ := json.Marshal(postBody)
	url := "http://154.39.79.242:8080/chat/api/chat_message/" + chat_id
	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		c.JSON(500, gin.H{"error": "创建请求失败"})
		return
	}
	httpReq.Header.Set("accept", "*/*")
	httpReq.Header.Set("Authorization", "Bearer application-8a25368429d8c9414eed9f269622b276")
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(500, gin.H{"error": "请求失败"})
		return
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Content string `json:"content"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(500, gin.H{"error": "解析响应失败"})
		return
	}
	content := result.Data.Content
	re := regexp.MustCompile("```json\\s*(\\{.*?\\})\\s*```")
	matches := re.FindStringSubmatch(content)
	var data map[string]interface{}
	if len(matches) >= 2 {
		if err := json.Unmarshal([]byte(matches[1]), &data); err == nil {
			// 去除json部分
			content = re.ReplaceAllString(content, "")
			c.JSON(200, gin.H{"content": content, "data": data})
			return
		}
	}
	// 没有json则只去除可能的json块
	content = re.ReplaceAllString(content, "")
	c.JSON(200, gin.H{"content": content})
}
