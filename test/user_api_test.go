package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"group_ten_server/router"
)

type UserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role,omitempty"`
}

func TestUserRegisterAndLogin(t *testing.T) {
	r := router.SetupRouter()

	// 1. 注册用户
	register := UserReq{
		Username: "testuser1",
		Password: "123456",
		Role:     "user",
	}
	body, _ := json.Marshal(register)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK && !bytes.Contains(w.Body.Bytes(), []byte("用户名已存在")) {
		t.Fatalf("注册失败: %s", w.Body.String())
	}

	// 2. 登录用户
	login := UserReq{
		Username: "testuser1",
		Password: "123456",
	}
	body, _ = json.Marshal(login)
	req = httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("登录失败: %s", w.Body.String())
	}
}
