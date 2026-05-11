package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type WechatHandler struct {
	appID     string
	appSecret string
	jwtSecret []byte
}

type WechatUser struct {
	OpenID    string `json:"openid"`
	UnionID   string `json:"unionid,omitempty"`
	Nickname  string `json:"nickname"`
	HeadImgURL string `json:"headimgurl"`
}

type WechatAccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID      string `json:"openid"`
	Scope       string `json:"scope"`
}

type WechatCode2SessionResp struct {
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	OpenID     string `json:"openid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func NewWechatHandler(appID, appSecret, jwtSecret string) *WechatHandler {
	return &WechatHandler{
		appID:     appID,
		appSecret: appSecret,
		jwtSecret: []byte(jwtSecret),
	}
}

// GetAuthorizeURL 获取微信授权URL
func (h *WechatHandler) GetAuthorizeURL(w http.ResponseWriter, r *http.Request) {
	redirectURI := r.URL.Query().Get("redirect_uri")
	state := r.URL.Query().Get("state")

	params := url.Values{}
	params.Set("appid", h.appID)
	params.Set("redirect_uri", redirectURI)
	params.Set("response_type", "code")
	params.Set("scope", "snsapi_login")
	params.Set("state", state)

	authURL := fmt.Sprintf("https://open.weixin.qq.com/connect/qrconnect?%s", params.Encode())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": authURL})
}

// Code2Session 微信code换session
func (h *WechatHandler) Code2Session(code string) (*WechatCode2SessionResp, error) {
	params := url.Values{}
	params.Set("appid", h.appID)
	params.Set("secret", h.appSecret)
	params.Set("js_code", code)
	params.Set("grant_type", "authorization_code")

	resp, err := http.Get("https://api.weixin.qq.com/sns/jscode2session?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result WechatCode2SessionResp
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("wechat error: %s", result.ErrMsg)
	}

	return &result, nil
}

// GenerateToken 生成JWT token
func (h *WechatHandler) GenerateToken(openID, nickname string) (string, error) {
	claims := jwt.MapClaims{
		"openid":   openID,
		"nickname": nickname,
		"exp":      time.Now().Add(24 * 7 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.jwtSecret)
}

// AuthCallback 授权回调处理
func (h *WechatHandler) AuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	// 用code换取session
	session, err := h.Code2Session(code)
	if err != nil {
		http.Error(w, "failed to get session", http.StatusInternalServerError)
		return
	}

	// 生成应用token
	token, err := h.GenerateToken(session.OpenID, state)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":   token,
		"openid":  session.OpenID,
		"unionid": session.UnionID,
	})
}

// MockWechatLogin 模拟微信登录（用于开发/演示）
func (h *WechatHandler) MockWechatLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OpenID  string `json:"openid"`
		Name    string `json:"name"`
		Avatar  string `json:"avatar"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.OpenID == "" {
		req.OpenID = "mock_" + uuid.New().String()
	}
	if req.Name == "" {
		req.Name = "微信用户"
	}

	// 生成token
	token, err := h.GenerateToken(req.OpenID, req.Name)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       req.OpenID,
			"name":     req.Name,
			"avatar":   req.Avatar,
			"platform": "wechat",
		},
	})
}

// DemoLogin 演示模式登录（无需微信）
func (h *WechatHandler) DemoLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		req.Name = "演示用户"
	}

	openID := "demo_" + uuid.New().String()
	token, err := h.GenerateToken(openID, req.Name)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       openID,
			"name":     req.Name,
			"avatar":   "",
			"platform": "demo",
			"demo":     true,
		},
	})
}

// RegisterWithPassword 注册账号
func (h *WechatHandler) RegisterWithPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// 验证输入
	req.Username = strings.TrimSpace(req.Username)
	req.Name = strings.TrimSpace(req.Name)

	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password required", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		req.Name = req.Username
	}

	// 密码哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	// 生成用户ID和token
	userID := uuid.New().String()
	token, err := h.GenerateToken(userID, req.Name)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       userID,
			"name":     req.Name,
			"username": req.Username,
			"avatar":   "",
			"platform": "password",
		},
		"password_hash": string(hash),
	})
}

// PasswordLogin 账号密码登录
func (h *WechatHandler) PasswordLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	req.Username = strings.TrimSpace(req.Username)

	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password required", http.StatusBadRequest)
		return
	}

	// 演示模式：任意账号密码都能登录
	// 生产环境需要查询数据库验证
	userID := "user_" + req.Username
	token, err := h.GenerateToken(userID, req.Username)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       userID,
			"name":     req.Username,
			"username": req.Username,
			"avatar":   "",
			"platform": "password",
		},
	})
}
