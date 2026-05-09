package main

import (
	"log"
	"net/http"

	"meeting-server/internal/handler"
	"meeting-server/internal/repository"
	"meeting-server/internal/service"
	"meeting-server/internal/signal"
)

func main() {
	// 初始化组件
	meetingRepo := repository.NewMeetingRepository()
	meetingService := service.NewMeetingService(meetingRepo)
	meetingHandler := handler.NewMeetingHandler(meetingService)

	// 信令服务
	hub := signal.NewHub()
	go hub.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", hub.HandleWS)
	meetingHandler.RegisterRoutes(mux)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
