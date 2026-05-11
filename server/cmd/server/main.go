package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"meeting-server/internal/handler"
	"meeting-server/internal/repository"
	"meeting-server/internal/service"
	"meeting-server/internal/signal"
	"meeting-server/internal/sfu"
)

func main() {
	// 数据库连接
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/meeting?sslmode=disable"
	}
	db, err := repository.NewDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 初始化 Repository
	meetingRepo := repository.NewMeetingRepository(db)
	userRepo := repository.NewUserRepository(db)

	// 初始化 Service
	meetingService := service.NewMeetingService(meetingRepo)
	userService := service.NewUserService(userRepo)

	// 初始化 Handler
	meetingHandler := handler.NewMeetingHandler(meetingService)
	authHandler := handler.NewAuthHandler(userService)

	// 微信 OAuth Handler
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}
	wechatHandler := handler.NewWechatHandler(
		os.Getenv("WECHAT_APP_ID"),
		os.Getenv("WECHAT_APP_SECRET"),
		jwtSecret,
	)

	// SFU 客户端
	sfuURL := os.Getenv("SFU_URL")
	if sfuURL == "" {
		sfuURL = "http://localhost:3000"
	}
	sfuClient := sfu.NewClient(sfuURL)
	roomHandler := handler.NewRoomHandler(sfuClient)

	// 信令服务
	hub := signal.NewHub(sfuClient)
	go hub.Run()

	// 路由
	mux := http.NewServeMux()

	// 静态文件服务
	mux.Handle("/", http.FileServer(http.Dir("/Users/imac/Desktop/自用会议室/web")))

	// WebSocket 信令
	mux.HandleFunc("/ws", hub.HandleWS)

	// 公开 API
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/me", authHandler.Me)

	// 微信 OAuth API
	mux.HandleFunc("/api/wechat/authorize", wechatHandler.GetAuthorizeURL)
	mux.HandleFunc("/api/wechat/callback", wechatHandler.AuthCallback)
	mux.HandleFunc("/api/wechat/mock-login", wechatHandler.MockWechatLogin)
	mux.HandleFunc("/api/demo/login", wechatHandler.DemoLogin)
	mux.HandleFunc("/api/auth/register/password", wechatHandler.RegisterWithPassword)
	mux.HandleFunc("/api/auth/login/password", wechatHandler.PasswordLogin)

	// 会议 API
	mux.HandleFunc("/api/meeting/create", meetingHandler.CreateMeeting)
	mux.HandleFunc("/api/meeting/get", meetingHandler.GetMeeting)
	mux.HandleFunc("/api/meeting/join", meetingHandler.JoinMeeting)

	// 房间 API
	mux.HandleFunc("/api/room/create", roomHandler.CreateRoom)
	mux.HandleFunc("/api/room/get", roomHandler.GetRoom)

	// 健康检查
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}