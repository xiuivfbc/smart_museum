package my_init

import (
	"log"
	"net/http"
)

func GetMaxKBPower() {
	url := "http://154.39.79.242:8080/chat/api/application/profile"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("创建maxkb授权请求失败:", err)
		return
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("Authorization", "Bearer application-8a25368429d8c9414eed9f269622b276")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("maxkb授权请求失败:", err)
	}
	resp.Body.Close()
}
