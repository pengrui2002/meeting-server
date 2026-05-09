# 企业视频会议系统 - 实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 搭建一个私有化部署的企业视频会议系统，支持 500 人并发，具备视频通话、屏幕共享、录制等核心功能

**Architecture:**
- **后端**: Go + mediasoup SFU，支持横向扩展
- **客户端**: Flutter 跨平台 (iOS/Android/PC)
- **信令**: WebSocket 长连接
- **媒体**: WebRTC SFU 架构

**Tech Stack:**
- Go 1.21+
- mediasoup 3.x
- Flutter 3.x
- WebRTC (flutter_webrtc)
- PostgreSQL
- Redis

---

## 项目结构

```
meeting-system/
├── server/                    # 后端服务
│   ├── cmd/server/           # 主入口
│   ├── internal/
│   │   ├── config/          # 配置
│   │   ├── handler/         # HTTP handler
│   │   ├── signal/          # 信令服务
│   │   ├── media/           # mediasoup 封装
│   │   ├── model/           # 数据模型
│   │   ├── repository/      # 数据层
│   │   └── service/         # 业务逻辑
│   └── pkg/
│       └── middleware/       # 中间件
│
├── client/                   # Flutter 客户端
│   └── lib/
│       ├── main.dart
│       ├── config/          # 配置
│       ├── models/          # 数据模型
│       ├── pages/           # 页面
│       ├── services/        # 服务 (WebRTC, API)
│       └── widgets/          # 组件
│
└── deploy/                   # 部署配置
    └── docker-compose.yml
```

---

## 任务列表

### 阶段一：后端基础设施

#### Task 1: 项目初始化与配置

**Files:**
- Create: `server/go.mod`
- Create: `server/cmd/server/main.go`
- Create: `server/internal/config/config.go`
- Create: `server/internal/model/meeting.go`

- [ ] **Step 1: 初始化 Go module**

```bash
cd /Users/imac/Desktop/自用会议室/server && go mod init meeting-server
```

- [ ] **Step 2: 创建目录结构**

```bash
mkdir -p server/cmd/server server/internal/{config,handler,signal,media,model,repository,service}
```

- [ ] **Step 3: 编写配置文件**

`server/internal/config/config.go`:
```go
package config

type Config struct {
    Server  ServerConfig
    Media   MediaConfig
    Database DatabaseConfig
    Redis   RedisConfig
}

type ServerConfig struct {
    Host string
    Port int
}

type MediaConfig struct {
    WorkerPort int
    RtcPort   int
}

type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    Database string
}

type RedisConfig struct {
    Host string
    Port int
}
```

- [ ] **Step 4: 创建数据模型**

`server/internal/model/meeting.go`:
```go
package model

import "time"

type Meeting struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Password  string    `json:"-"`
    CreatedAt time.Time `json:"created_at"`
    Status    string    `json:"status"` // pending, active, ended
}

type Participant struct {
    ID       string `json:"id"`
    MeetingID string `json:"meeting_id"`
    UserID   string `json:"user_id"`
    Role     string `json:"role"` // host, participant
}
```

- [ ] **Step 5: 提交**

```bash
cd /Users/imac/Desktop/自用会议室 && git init && git add server/go.mod server/cmd server/internal/config server/internal/model && git commit -m "feat: init project structure"
```

---

#### Task 2: 信令服务 (WebSocket)

**Files:**
- Create: `server/internal/signal/server.go`
- Create: `server/internal/signal/room.go`
- Create: `server/internal/signal/client.go`
- Modify: `server/cmd/server/main.go`

- [ ] **Step 1: 实现 Room 管理**

`server/internal/signal/room.go`:
```go
package signal

import (
    "sync"
    "github.com/google/uuid"
)

type Room struct {
    ID          string
    clients     map[*Client]bool
    register    chan *Client
    unregister  chan *Client
    broadcast   chan []byte
    mu          sync.Mutex
}

func NewRoom() *Room {
    return &Room{
        ID:         uuid.New().String(),
        clients:    make(map[*Client]bool),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        broadcast:  make(chan []byte),
    }
}

func (r *Room) Run() {
    for {
        select {
        case client := <-r.register:
            r.mu.Lock()
            r.clients[client] = true
            r.mu.Unlock()

        case client := <-r.unregister:
            r.mu.Lock()
            if _, ok := r.clients[client]; ok {
                delete(r.clients, client)
                close(client.send)
            }
            r.mu.Unlock()

        case message := <-r.broadcast:
            r.mu.Lock()
            for client := range r.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(r.clients, client)
                }
            }
            r.mu.Unlock()
        }
    }
}
```

- [ ] **Step 2: 实现 Client 连接**

`server/internal/signal/client.go`:
```go
package signal

import (
    "encoding/json"
    "log"
    "time"

    "github.com/gorilla/websocket"
)

const (
    writeWait      = 10 * time.Second
    pongWait       = 60 * time.Second
    pingPeriod     = (pongWait * 9) / 10
    maxMessageSize = 512
)

type Message struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}

type Client struct {
    hub    *Hub
    conn   *websocket.Conn
    send   chan []byte
    roomID string
    userID string
}

func (c *Client) ReadPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()

    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }

        var msg Message
        if err := json.Unmarshal(message, &msg); err != nil {
            continue
        }

        c.handleMessage(&msg)
    }
}

func (c *Client) WritePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            w, err := c.conn.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(message)

            if err := w.Close(); err != nil {
                return
            }

        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}

func (c *Client) handleMessage(msg *Message) {
    switch msg.Type {
    case "join":
        c.hub.register <- c
    case "leave":
        c.hub.unregister <- c
    case "offer", "answer", "ice-candidate":
        c.hub.broadcast <- msg.Data
    }
}
```

- [ ] **Step 3: 实现 Hub (连接管理器)**

`server/internal/signal/server.go`:
```go
package signal

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type Hub struct {
    rooms      map[string]*Room
    register   chan *Client
    unregister chan *Client
}

func NewHub() *Hub {
    return &Hub{
        rooms:      make(map[string]*Room),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            room, ok := h.rooms[client.roomID]
            if !ok {
                room = NewRoom()
                h.rooms[client.roomID] = room
                go room.Run()
            }
            room.register <- client

        case client := <-h.unregister:
            if room, ok := h.rooms[client.roomID]; ok {
                room.unregister <- client
            }
        }
    }
}

func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request) {
    roomID := r.URL.Query().Get("room")
    userID := r.URL.Query().Get("user")

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("upgrade error: %v", err)
        return
    }

    client := &Client{
        hub:    h,
        conn:   conn,
        send:   make(chan []byte, 256),
        roomID: roomID,
        userID: userID,
    }

    h.register <- client

    go client.WritePump()
    go client.ReadPump()
}
```

- [ ] **Step 4: 更新 main.go**

`server/cmd/server/main.go`:
```go
package main

import (
    "log"
    "net/http"

    "meeting-server/internal/signal"
)

func main() {
    hub := signal.NewHub()
    go hub.Run()

    http.HandleFunc("/ws", hub.HandleWS)

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

- [ ] **Step 5: 添加依赖并测试编译**

```bash
cd server && go get github.com/gorilla/websocket github.com/google/uuid && go build ./cmd/server
```

- [ ] **Step 6: 提交**

```bash
git add server/internal/signal server/cmd/server/main.go && git commit -m "feat: add WebSocket signal server with room management"
```

---

#### Task 3: API 服务 (会议管理)

**Files:**
- Create: `server/internal/handler/meeting.go`
- Create: `server/internal/repository/meeting.go`
- Create: `server/internal/service/meeting.go`
- Modify: `server/cmd/server/main.go`

- [ ] **Step 1: 创建 Meeting Handler**

`server/internal/handler/meeting.go`:
```go
package handler

import (
    "encoding/json"
    "net/http"

    "meeting-server/internal/model"
    "meeting-server/internal/service"

    "github.com/google/uuid"
)

type MeetingHandler struct {
    service *service.MeetingService
}

func NewMeetingHandler(s *service.MeetingService) *MeetingHandler {
    return &MeetingHandler{service: s}
}

func (h *MeetingHandler) CreateMeeting(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Title    string `json:"title"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    meeting := &model.Meeting{
        ID:     uuid.New().String()[:8],
        Title:  req.Title,
        Status: "pending",
    }

    h.service.CreateMeeting(meeting)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(meeting)
}

func (h *MeetingHandler) GetMeeting(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")

    meeting, err := h.service.GetMeeting(id)
    if err != nil {
        http.Error(w, "meeting not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(meeting)
}

func (h *MeetingHandler) JoinMeeting(w http.ResponseWriter, r *http.Request) {
    var req struct {
        MeetingID string `json:"meeting_id"`
        UserID   string `json:"user_id"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    participant := &model.Participant{
        ID:        uuid.New().String(),
        MeetingID: req.MeetingID,
        UserID:    req.UserID,
        Role:      "participant",
    }

    h.service.AddParticipant(participant)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(participant)
}

func (h *MeetingHandler) RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/api/meeting/create", h.CreateMeeting)
    mux.HandleFunc("/api/meeting/get", h.GetMeeting)
    mux.HandleFunc("/api/meeting/join", h.JoinMeeting)
}
```

- [ ] **Step 2: 创建 Repository**

`server/internal/repository/meeting.go`:
```go
package repository

import (
    "sync"

    "meeting-server/internal/model"
)

type MeetingRepository struct {
    meetings      map[string]*model.Meeting
    participants  map[string][]*model.Participant
    mu            sync.RWMutex
}

func NewMeetingRepository() *MeetingRepository {
    return &MeetingRepository{
        meetings:     make(map[string]*model.Meeting),
        participants: make(map[string][]*model.Participant),
    }
}

func (r *MeetingRepository) Create(meeting *model.Meeting) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.meetings[meeting.ID] = meeting
}

func (r *MeetingRepository) GetByID(id string) (*model.Meeting, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    if meeting, ok := r.meetings[id]; ok {
        return meeting, nil
    }
    return nil, nil
}

func (r *MeetingRepository) AddParticipant(p *model.Participant) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.participants[p.MeetingID] = append(r.participants[p.MeetingID], p)
}

func (r *MeetingRepository) GetParticipants(meetingID string) []*model.Participant {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return r.participants[meetingID]
}
```

- [ ] **Step 3: 创建 Service**

`server/internal/service/meeting.go`:
```go
package service

import (
    "meeting-server/internal/model"
    "meeting-server/internal/repository"
)

type MeetingService struct {
    repo *repository.MeetingRepository
}

func NewMeetingService(r *repository.MeetingRepository) *MeetingService {
    return &MeetingService{repo: r}
}

func (s *MeetingService) CreateMeeting(meeting *model.Meeting) {
    s.repo.Create(meeting)
}

func (s *MeetingService) GetMeeting(id string) (*model.Meeting, error) {
    return s.repo.GetByID(id)
}

func (s *MeetingService) AddParticipant(p *model.Participant) {
    s.repo.AddParticipant(p)
}
```

- [ ] **Step 4: 更新 main.go**

```go
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
```

- [ ] **Step 5: 提交**

```bash
git add server/internal/handler server/internal/repository server/internal/service server/cmd/server/main.go && git commit -m "feat: add REST API for meeting management"
```

---

#### Task 4: mediasoup SFU 集成

**Files:**
- Create: `server/internal/media/sfu.go`
- Modify: `server/cmd/server/main.go`

- [ ] **Step 1: 创建 SFU 封装**

`server/internal/media/sfu.go`:
```go
package media

import (
    "log"
    "net"
    "github.com/versori/mediasoup"
)

type SFUServer struct {
    worker *mediasoup.Worker
    router *mediasoup.Router
}

func NewSFUServer() (*SFUServer, error) {
    // 创建 Worker
    worker, err := mediasoup.NewWorker([]string{
        "mediasoup:listenIp=127.0.0.1",
        "mediasoup:numWorkers=4",
    })
    if err != nil {
        return nil, err
    }

    // 创建 Router
    router, err := worker.CreateRouter()
    if err != nil {
        return nil, err
    }

    return &SFUServer{
        worker: worker,
        router: router,
    }, nil
}

func (s *SFUServer) CreateTransport() (*mediasoup.Transport, error) {
    transport, err := s.router.CreateWebRtcTransport(mediasoup.WebRtcTransportOptions{
        ListenIps: []mediasoup.TransportListenIp{
            {Ip: "127.0.0.1"},
        },
        InitialAvailableOutgoingBitrate: 1000000,
    })
    return transport, err
}

func (s *SFUServer) Close() {
    if s.worker != nil {
        s.worker.Close()
    }
}
```

- [ ] **Step 2: 创建 Room 路由**

`server/internal/media/room.go`:
```go
package media

import (
    "sync"
    "github.com/versori/mediasoup"
)

type Room struct {
    ID        string
    router    *mediasoup.Router
    peers     map[string]*Peer
    transports map[string]*mediasoup.Transport
    mu        sync.RWMutex
}

type Peer struct {
    ID         string
    DisplayName string
    transports map[string]*mediasoup.Transport
    producers  map[string]*mediasoup.Producer
    consumers  map[string]*mediasoup.Consumer
}

func NewRoom(router *mediasoup.Router, id string) *Room {
    return &Room{
        ID:         id,
        router:     router,
        peers:      make(map[string]*Peer),
        transports: make(map[string]*mediasoup.Transport),
    }
}

func (r *Room) AddPeer(peer *Peer) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.peers[peer.ID] = peer
}

func (r *Room) RemovePeer(peerID string) {
    r.mu.Lock()
    defer r.mu.Unlock()
    delete(r.peers, peerID)
}
```

- [ ] **Step 3: 提交**

```bash
git add server/internal/media server/cmd/server/main.go && git commit -m "feat: add mediasoup SFU integration"
```

---

### 阶段二：Flutter 客户端

#### Task 5: Flutter 项目初始化

**Files:**
- Create: `client/pubspec.yaml`
- Create: `client/lib/main.dart`
- Create: `client/lib/config/app_config.dart`
- Create: `client/lib/models/meeting.dart`

- [ ] **Step 1: 创建 Flutter 项目**

```bash
cd /Users/imac/Desktop/自用会议室 && flutter create --org com.example meeting_client
```

- [ ] **Step 2: 配置 pubspec.yaml**

`client/pubspec.yaml`:
```yaml
name: meeting_client
description: 企业视频会议系统
publish_to: 'none'
version: 1.0.0+1

environment:
  sdk: '>=3.0.0 <4.0.0'

dependencies:
  flutter:
    sdk: flutter
  flutter_webrtc: ^0.11.0
  web_socket_channel: ^2.4.0
  provider: ^6.1.0
  http: ^1.1.0
  uuid: ^4.2.1

dev_dependencies:
  flutter_test:
    sdk: flutter
  flutter_lints: ^3.0.0

flutter:
  uses-material-design: true
```

- [ ] **Step 3: 创建应用配置**

`client/lib/config/app_config.dart`:
```dart
class AppConfig {
  static const String serverUrl = 'http://localhost:8080';
  static const String wsUrl = 'ws://localhost:8080/ws';
  static const String apiBase = '$serverUrl/api';
}
```

- [ ] **Step 4: 创建数据模型**

`client/lib/models/meeting.dart`:
```dart
class Meeting {
  final String id;
  final String title;
  final String status;

  Meeting({
    required this.id,
    required this.title,
    required this.status,
  });

  factory Meeting.fromJson(Map<String, dynamic> json) {
    return Meeting(
      id: json['id'] ?? '',
      title: json['title'] ?? '',
      status: json['status'] ?? 'pending',
    );
  }
}

class Participant {
  final String id;
  final String meetingId;
  final String userId;
  final String role;

  Participant({
    required this.id,
    required this.meetingId,
    required this.userId,
    required this.role,
  });

  factory Participant.fromJson(Map<String, dynamic> json) {
    return Participant(
      id: json['id'] ?? '',
      meetingId: json['meeting_id'] ?? '',
      userId: json['user_id'] ?? '',
      role: json['role'] ?? 'participant',
    );
  }
}
```

- [ ] **Step 5: 创建入口文件**

`client/lib/main.dart`:
```dart
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'pages/home_page.dart';
import 'services/meeting_service.dart';
import 'services/webrtc_service.dart';

void main() {
  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => MeetingService()),
        ChangeNotifierProvider(create: (_) => WebRTCService()),
      ],
      child: const MeetingApp(),
    ),
  );
}

class MeetingApp extends StatelessWidget {
  const MeetingApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: '视频会议',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        useMaterial3: true,
      ),
      home: const HomePage(),
    );
  }
}
```

- [ ] **Step 6: 提交**

```bash
git add client/pubspec.yaml client/lib/ && git commit -m "feat: init Flutter project structure"
```

---

#### Task 6: 会议服务 (API + WebSocket)

**Files:**
- Create: `client/lib/services/meeting_service.dart`
- Create: `client/lib/services/webrtc_service.dart`

- [ ] **Step 1: 创建 MeetingService**

`client/lib/services/meeting_service.dart`:
```dart
import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;
import '../config/app_config.dart';
import '../models/meeting.dart';

class MeetingService extends ChangeNotifier {
  Meeting? _currentMeeting;
  final List<Participant> _participants = [];
  bool _isLoading = false;

  Meeting? get currentMeeting => _currentMeeting;
  List<Participant> get participants => _participants;
  bool get isLoading => _isLoading;

  Future<Meeting> createMeeting(String title, {String? password}) async {
    _isLoading = true;
    notifyListeners();

    final response = await http.post(
      Uri.parse('${AppConfig.apiBase}/meeting/create'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({
        'title': title,
        'password': password ?? '',
      }),
    );

    if (response.statusCode == 200) {
      _currentMeeting = Meeting.fromJson(jsonDecode(response.body));
    }

    _isLoading = false;
    notifyListeners();
    return _currentMeeting!;
  }

  Future<Meeting?> getMeeting(String id) async {
    final response = await http.get(
      Uri.parse('${AppConfig.apiBase}/meeting/get?id=$id'),
    );

    if (response.statusCode == 200) {
      _currentMeeting = Meeting.fromJson(jsonDecode(response.body));
      notifyListeners();
      return _currentMeeting;
    }
    return null;
  }

  Future<Participant?> joinMeeting(String meetingId, String userId, {String? password}) async {
    final response = await http.post(
      Uri.parse('${AppConfig.apiBase}/meeting/join'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({
        'meeting_id': meetingId,
        'user_id': userId,
        'password': password ?? '',
      }),
    );

    if (response.statusCode == 200) {
      final participant = Participant.fromJson(jsonDecode(response.body));
      _participants.add(participant);
      notifyListeners();
      return participant;
    }
    return null;
  }
}
```

- [ ] **Step 2: 创建 WebRTC 服务**

`client/lib/services/webrtc_service.dart`:
```dart
import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:flutter_webrtc/flutter_webrtc.dart';
import 'package:web_socket_channel/web_socket_channel.dart';
import '../config/app_config.dart';

class WebRTCService extends ChangeNotifier {
  RTCPeerConnection? _peerConnection;
  MediaStream? _localStream;
  final Map<String, MediaStream> _remoteStreams = {};
  WebSocketChannel? _channel;
  String? _roomId;
  String? _userId;
  bool _isConnected = false;
  bool _isMuted = false;
  bool _isVideoEnabled = true;

  RTCPeerConnection? get peerConnection => _peerConnection;
  MediaStream? get localStream => _localStream;
  Map<String, MediaStream> get remoteStreams => _remoteStreams;
  bool get isConnected => _isConnected;
  bool get isMuted => _isMuted;
  bool get isVideoEnabled => _isVideoEnabled;

  Future<void> joinRoom(String roomId, String userId) async {
    _roomId = roomId;
    _userId = userId;

    // 初始化本地媒体流
    final Map<String, dynamic> mediaConstraints = {
      'audio': true,
      'video': {
        'facingMode': 'user',
        'width': 640,
        'height': 480,
      },
    };

    _localStream = await navigator.mediaDevices.getUserMedia(mediaConstraints);

    // 创建 WebSocket 连接
    _channel = WebSocketChannel.connect(
      Uri.parse('${AppConfig.wsUrl}?room=$roomId&user=$userId'),
    );

    // 监听消息
    _channel!.stream.listen((message) {
      _handleSignalingMessage(message);
    });

    // 初始化 RTCPeerConnection
    await _initPeerConnection();

    _isConnected = true;
    notifyListeners();
  }

  Future<void> _initPeerConnection() async {
    final Map<String, dynamic> config = {
      'iceServers': [
        {'urls': 'stun:stun.l.google.com:19302'},
      ]
    };

    _peerConnection = await createPeerConnection(config);

    // 添加本地流轨道
    _localStream!.getTracks().forEach((track) {
      _peerConnection!.addTrack(track, _localStream!);
    });

    // 监听远程流
    _peerConnection!.onTrack = (RTCTrackEvent event) {
      if (event.streams.isNotEmpty) {
        _remoteStreams[event.streams[0].id] = event.streams[0];
        notifyListeners();
      }
    };
  }

  void _handleSignalingMessage(dynamic message) {
    final data = jsonDecode(message);

    switch (data['type']) {
      case 'offer':
        _handleOffer(data['sdp']);
        break;
      case 'answer':
        _handleAnswer(data['sdp']);
        break;
      case 'ice-candidate':
        _handleIceCandidate(data);
        break;
    }
  }

  Future<void> _handleOffer(String sdp) async {
    await _peerConnection!.setRemoteDescription(
      RTCSessionDescription(sdp, 'offer'),
    );
    final answer = await _peerConnection!.createAnswer();
    await _peerConnection!.setLocalDescription(answer);
    _sendMessage({'type': 'answer', 'sdp': answer.sdp});
  }

  Future<void> _handleAnswer(String sdp) async {
    await _peerConnection!.setRemoteDescription(
      RTCSessionDescription(sdp, 'answer'),
    );
  }

  Future<void> _handleIceCandidate(Map<String, dynamic> data) async {
    final candidate = RTCIceCandidate(
      data['candidate'],
      data['sdpMid'],
      data['sdpMLineIndex'],
    );
    await _peerConnection!.addCandidate(candidate);
  }

  void _sendMessage(Map<String, dynamic> message) {
    _channel?.sink.add(jsonEncode(message));
  }

  Future<void> toggleMute() async {
    if (_localStream != null) {
      bool enabled = _localStream!.getAudioTracks().first.enabled;
      _localStream!.getAudioTracks().first.enabled = !enabled;
      _isMuted = !enabled;
      notifyListeners();
    }
  }

  Future<void> toggleVideo() async {
    if (_localStream != null) {
      bool enabled = _localStream!.getVideoTracks().first.enabled;
      _localStream!.getVideoTracks().first.enabled = !enabled;
      _isVideoEnabled = !enabled;
      notifyListeners();
    }
  }

  Future<void> switchCamera() async {
    if (_localStream != null) {
      await Helper.switchCamera(_localStream!.getVideoTracks().first);
    }
  }

  Future<void> leaveRoom() async {
    await _localStream?.dispose();
    await _peerConnection?.close();
    await _channel?.sink.close();
    _localStream = null;
    _peerConnection = null;
    _channel = null;
    _remoteStreams.clear();
    _isConnected = false;
    notifyListeners();
  }

  @override
  void dispose() {
    leaveRoom();
    super.dispose();
  }
}
```

- [ ] **Step 3: 提交**

```bash
git add client/lib/services/ && git commit -m "feat: add meeting service and WebRTC service"
```

---

#### Task 7: 页面开发

**Files:**
- Create: `client/lib/pages/home_page.dart`
- Create: `client/lib/pages/meeting_page.dart`
- Create: `client/lib/pages/create_meeting_page.dart`
- Create: `client/lib/pages/join_meeting_page.dart`

- [ ] **Step 1: 创建首页**

`client/lib/pages/home_page.dart`:
```dart
import 'package:flutter/material.dart';
import 'create_meeting_page.dart';
import 'join_meeting_page.dart';

class HomePage extends StatelessWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('视频会议'),
        centerTitle: true,
      ),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(24.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                Icons.video_call,
                size: 120,
                color: Theme.of(context).primaryColor,
              ),
              const SizedBox(height: 48),
              const Text(
                '企业视频会议',
                style: TextStyle(
                  fontSize: 28,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 8),
              Text(
                '安全、私有化部署',
                style: TextStyle(
                  fontSize: 16,
                  color: Colors.grey[600],
                ),
              ),
              const SizedBox(height: 64),
              SizedBox(
                width: double.infinity,
                height: 56,
                child: ElevatedButton.icon(
                  onPressed: () {
                    Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder: (_) => const CreateMeetingPage(),
                      ),
                    );
                  },
                  icon: const Icon(Icons.add),
                  label: const Text('创建会议'),
                  style: ElevatedButton.styleFrom(
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                ),
              ),
              const SizedBox(height: 16),
              SizedBox(
                width: double.infinity,
                height: 56,
                child: OutlinedButton.icon(
                  onPressed: () {
                    Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder: (_) => const JoinMeetingPage(),
                      ),
                    );
                  },
                  icon: const Icon(Icons.login),
                  label: const Text('加入会议'),
                  style: OutlinedButton.styleFrom(
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
```

- [ ] **Step 2: 创建会议页面**

`client/lib/pages/meeting_page.dart`:
```dart
import 'package:flutter/material.dart';
import 'package:flutter_webrtc/flutter_webrtc.dart';
import 'package:provider/provider.dart';
import '../services/webrtc_service.dart';

class MeetingPage extends StatefulWidget {
  final String meetingId;

  const MeetingPage({super.key, required this.meetingId});

  @override
  State<MeetingPage> createState() => _MeetingPageState();
}

class _MeetingPageState extends State<MeetingPage> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      final webrtc = context.read<WebRTCService>();
      webrtc.joinRoom(widget.meetingId, 'user_${DateTime.now().millisecondsSinceEpoch}');
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.black,
      appBar: AppBar(
        backgroundColor: Colors.black,
        title: Text('会议号: ${widget.meetingId}'),
        actions: [
          IconButton(
            icon: const Icon(Icons.copy),
            onPressed: () {
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('会议号已复制')),
              );
            },
          ),
        ],
      ),
      body: Consumer<WebRTCService>(
        builder: (context, webrtc, child) {
          return Column(
            children: [
              // 视频区域
              Expanded(
                child: _buildVideoGrid(webrtc),
              ),
              // 控制栏
              _buildControlBar(webrtc),
            ],
          );
        },
      ),
    );
  }

  Widget _buildVideoGrid(WebRTCService webrtc) {
    final streams = webrtc.remoteStreams.values.toList();

    if (streams.isEmpty) {
      return const Center(
        child: Text(
          '等待其他人加入...',
          style: TextStyle(color: Colors.white70),
        ),
      );
    }

    return GridView.builder(
      padding: const EdgeInsets.all(8),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        childAspectRatio: 16 / 9,
        crossAxisSpacing: 8,
        mainAxisSpacing: 8,
      ),
      itemCount: streams.length + 1, // +1 for local video
      itemBuilder: (context, index) {
        if (index == 0) {
          return _buildVideoView(webrtc.localStream, '我', true);
        }
        return _buildVideoView(streams[index - 1], '用户$index', false);
      },
    );
  }

  Widget _buildVideoView(MediaStream? stream, String label, bool isLocal) {
    if (stream == null) {
      return Container(
        color: Colors.grey[900],
        child: const Center(
          child: Icon(Icons.videocam_off, color: Colors.white54, size: 48),
        ),
      );
    }

    return ClipRRect(
      borderRadius: BorderRadius.circular(8),
      child: Stack(
        fit: StackFit.expand,
        children: [
          RTCVideoView(
            stream.getVideoTracks().first,
            objectFit: RTCVideoViewObjectFit.RTCVideoViewObjectFitCover,
          ),
          Positioned(
            bottom: 8,
            left: 8,
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
              decoration: BoxDecoration(
                color: Colors.black54,
                borderRadius: BorderRadius.circular(4),
              ),
              child: Text(
                label,
                style: const TextStyle(color: Colors.white, fontSize: 12),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildControlBar(WebRTCService webrtc) {
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 16, horizontal: 24),
      color: Colors.grey[900],
      child: SafeArea(
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            _buildControlButton(
              icon: webrtc.isMuted ? Icons.mic_off : Icons.mic,
              label: webrtc.isMuted ? '取消静音' : '静音',
              onPressed: () => webrtc.toggleMute(),
              isActive: !webrtc.isMuted,
            ),
            _buildControlButton(
              icon: webrtc.isVideoEnabled ? Icons.videocam : Icons.videocam_off,
              label: webrtc.isVideoEnabled ? '关闭视频' : '开启视频',
              onPressed: () => webrtc.toggleVideo(),
              isActive: webrtc.isVideoEnabled,
            ),
            _buildControlButton(
              icon: Icons.cameraswitch,
              label: '切换摄像头',
              onPressed: () => webrtc.switchCamera(),
              isActive: true,
            ),
            _buildControlButton(
              icon: Icons.call_end,
              label: '挂断',
              onPressed: () {
                webrtc.leaveRoom();
                Navigator.pop(context);
              },
              isActive: false,
              isEndCall: true,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildControlButton({
    required IconData icon,
    required String label,
    required VoidCallback onPressed,
    required bool isActive,
    bool isEndCall = false,
  }) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        IconButton(
          onPressed: onPressed,
          icon: Icon(icon),
          style: IconButton.styleFrom(
            backgroundColor: isEndCall
                ? Colors.red
                : (isActive ? Colors.white24 : Colors.white10),
            padding: const EdgeInsets.all(16),
          ),
          iconSize: 28,
          color: Colors.white,
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: const TextStyle(color: Colors.white70, fontSize: 12),
        ),
      ],
    );
  }
}
```

- [ ] **Step 3: 创建会议页面**

`client/lib/pages/create_meeting_page.dart`:
```dart
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../services/meeting_service.dart';
import 'meeting_page.dart';

class CreateMeetingPage extends StatefulWidget {
  const CreateMeetingPage({super.key});

  @override
  State<CreateMeetingPage> createState() => _CreateMeetingPageState();
}

class _CreateMeetingPageState extends State<CreateMeetingPage> {
  final _titleController = TextEditingController();
  final _passwordController = TextEditingController();
  bool _isLoading = false;

  @override
  void dispose() {
    _titleController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  Future<void> _createMeeting() async {
    if (_titleController.text.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请输入会议主题')),
      );
      return;
    }

    setState(() => _isLoading = true);

    final meetingService = context.read<MeetingService>();
    final meeting = await meetingService.createMeeting(
      _titleController.text,
      password: _passwordController.text.isNotEmpty
          ? _passwordController.text
          : null,
    );

    setState(() => _isLoading = false);

    if (mounted) {
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(
          builder: (_) => MeetingPage(meetingId: meeting.id),
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('创建会议'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(24.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            TextField(
              controller: _titleController,
              decoration: const InputDecoration(
                labelText: '会议主题',
                hintText: '请输入会议主题',
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: _passwordController,
              decoration: const InputDecoration(
                labelText: '会议密码（可选）',
                hintText: '请输入会议密码',
                border: OutlineInputBorder(),
              ),
              obscureText: true,
            ),
            const SizedBox(height: 32),
            SizedBox(
              width: double.infinity,
              height: 48,
              child: ElevatedButton(
                onPressed: _isLoading ? null : _createMeeting,
                child: _isLoading
                    ? const CircularProgressIndicator(color: Colors.white)
                    : const Text('创建会议'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
```

- [ ] **Step 4: 创建加入会议页面**

`client/lib/pages/join_meeting_page.dart`:
```dart
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../services/meeting_service.dart';
import 'meeting_page.dart';

class JoinMeetingPage extends StatefulWidget {
  const JoinMeetingPage({super.key});

  @override
  State<JoinMeetingPage> createState() => _JoinMeetingPageState();
}

class _JoinMeetingPageState extends State<JoinMeetingPage> {
  final _meetingIdController = TextEditingController();
  final _passwordController = TextEditingController();
  bool _isLoading = false;

  @override
  void dispose() {
    _meetingIdController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  Future<void> _joinMeeting() async {
    if (_meetingIdController.text.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请输入会议号')),
      );
      return;
    }

    setState(() => _isLoading = true);

    final meetingService = context.read<MeetingService>();
    final participant = await meetingService.joinMeeting(
      _meetingIdController.text,
      'user_${DateTime.now().millisecondsSinceEpoch}',
      password: _passwordController.text.isNotEmpty
          ? _passwordController.text
          : null,
    );

    setState(() => _isLoading = false);

    if (participant != null && mounted) {
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(
          builder: (_) => MeetingPage(meetingId: _meetingIdController.text),
        ),
      );
    } else if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('加入会议失败，请检查会议号和密码')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('加入会议'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(24.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            TextField(
              controller: _meetingIdController,
              decoration: const InputDecoration(
                labelText: '会议号',
                hintText: '请输入会议号',
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: _passwordController,
              decoration: const InputDecoration(
                labelText: '会议密码（可选）',
                hintText: '请输入会议密码',
                border: OutlineInputBorder(),
              ),
              obscureText: true,
            ),
            const SizedBox(height: 32),
            SizedBox(
              width: double.infinity,
              height: 48,
              child: ElevatedButton(
                onPressed: _isLoading ? null : _joinMeeting,
                child: _isLoading
                    ? const CircularProgressIndicator(color: Colors.white)
                    : const Text('加入会议'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
```

- [ ] **Step 5: 提交**

```bash
git add client/lib/pages/ && git commit -m "feat: add meeting pages (home, create, join, meeting)"
```

---

### 阶段三：部署与测试

#### Task 8: Docker 部署配置

**Files:**
- Create: `deploy/docker-compose.yml`
- Create: `deploy/nginx.conf`
- Create: `server/Dockerfile`

- [ ] **Step 1: 创建 docker-compose.yml**

`deploy/docker-compose.yml`:
```yaml
version: '3.8'

services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: meeting
      POSTGRES_USER: meeting
      POSTGRES_PASSWORD: meeting123
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  signal:
    build: ../server
    command: ./server
    ports:
      - "8080:8080"
    depends_on:
      - redis
      - postgres
    environment:
      - DATABASE_URL=postgres://meeting:meeting123@postgres:5432/meeting
      - REDIS_URL=redis://redis:6379

  mediasoup:
    image: versori/mediasoup:latest
    ports:
      - "40000-40100:40000-40100/udp"
    environment:
      - MEDIASOUP_LISTEN_IP=0.0.0.0
      - MEDIASOUP_ANNOUNCED_IP=127.0.0.1

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - signal

volumes:
  postgres_data:
```

- [ ] **Step 2: 创建 Nginx 配置**

`deploy/nginx.conf`:
```nginx
events {
    worker_connections 1024;
}

http {
    upstream signal {
        server signal:8080;
    }

    server {
        listen 80;
        server_name localhost;

        location / {
            proxy_pass http://signal;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_set_header Host $host;
        }

        location /api {
            proxy_pass http://signal;
            proxy_set_header Host $host;
        }
    }
}
```

- [ ] **Step 3: 创建 Dockerfile**

`server/Dockerfile`:
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080

CMD ["./server"]
```

- [ ] **Step 4: 提交**

```bash
git add deploy/ && git commit -m "feat: add Docker deployment configuration"
```

---

## 实施检查清单

| 阶段 | 任务 | 状态 |
|:---|:---|:---|
| 后端基础设施 | Task 1: 项目初始化 | ⬜ |
| 后端基础设施 | Task 2: 信令服务 | ⬜ |
| 后端基础设施 | Task 3: API 服务 | ⬜ |
| 后端基础设施 | Task 4: mediasoup SFU | ⬜ |
| Flutter 客户端 | Task 5: 项目初始化 | ⬜ |
| Flutter 客户端 | Task 6: 服务层 | ⬜ |
| Flutter 客户端 | Task 7: 页面开发 | ⬜ |
| 部署测试 | Task 8: Docker 部署 | ⬜ |

---

## 验证步骤

```bash
# 1. 后端编译测试
cd server && go build ./cmd/server

# 2. Flutter 编译测试
cd client && flutter pub get && flutter build apk --debug

# 3. 启动服务
cd deploy && docker-compose up -d

# 4. 查看日志
docker-compose logs -f
```
