package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"group_ten_server/router"
)

type EnvReq struct {
	Name     string  `json:"name"`
	Temp     float64 `json:"temp"`
	Humidity float64 `json:"humidity"`
	Status   string  `json:"status"`
	Info     string  `json:"info"`
	Light    float64 `json:"light"`
	Flame    bool    `json:"flame"`
}

func TestEnvironmentAPI(t *testing.T) {
	r := router.SetupRouter()

	// 1. 创建环境
	create := EnvReq{
		Name:     "testenv",
		Temp:     25.5,
		Humidity: 60,
		Status:   "normal",
		Info:     "测试环境",
		Light:    300,
		Flame:    false,
	}
	body, _ := json.Marshal(create)
	req := httptest.NewRequest("POST", "/environments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("创建环境失败: %s", w.Body.String())
	}

	// 2. 查询所有环境
	req = httptest.NewRequest("GET", "/environments", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("查询所有环境失败: %s", w.Body.String())
	}

	// 3. 查询单个环境
	req = httptest.NewRequest("GET", "/environments/testenv", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("查询单个环境失败: %s", w.Body.String())
	}

	// 4. 更新环境
	update := EnvReq{
		Temp:     26,
		Humidity: 65,
		Status:   "updated",
		Info:     "更新环境",
		Light:    350,
		Flame:    true,
	}
	body, _ = json.Marshal(update)
	req = httptest.NewRequest("PUT", "/environments/testenv", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("更新环境失败: %s", w.Body.String())
	}

	// 5. 删除环境
	req = httptest.NewRequest("DELETE", "/environments/testenv", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("删除环境失败: %s", w.Body.String())
	}
}
