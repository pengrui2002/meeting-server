# AFT集团线上服务平台 - 部署指南

## 快速部署

### 方式一：Docker 一键部署（推荐）

```bash
# 1. 进入部署目录
cd deploy

# 2. 一键部署
./deploy.sh
```

### 方式二：手动部署

```bash
# 1. 进入部署目录
cd deploy

# 2. 启动服务
docker-compose up -d

# 3. 查看状态
docker-compose ps
```

## 服务访问

| 服务 | 地址 |
|:---|:---|
| 前端页面 | http://服务器IP/standalone.html |
| API 健康检查 | http://服务器IP/health |

## 配置说明

### 环境变量

| 变量 | 说明 | 默认值 |
|:---|:---|:---|
| SFU_IP | SFU 服务器公网 IP | 127.0.0.1 |
| WECHAT_APP_ID | 微信开放平台 AppID | 空 |
| WECHAT_APP_SECRET | 微信开放平台 AppSecret | 空 |
| JWT_SECRET | JWT 密钥 | your-secret-key-change-in-production |

### 微信登录配置

如需微信扫码登录，需：
1. 申请微信开放平台账号
2. 设置环境变量后重启服务

```bash
export WECHAT_APP_ID=wx1234567890abcdef
export WECHAT_APP_SECRET=your_app_secret
docker-compose up -d
```

## 端口说明

| 端口 | 服务 | 说明 |
|:---|:---|:---|
| 80 | Nginx | HTTP |
| 443 | Nginx | HTTPS（需配置证书）|
| 3000 | SFU | mediasoup HTTP API |
| 40000-40100 | SFU | WebRTC 媒体传输 |
| 5432 | PostgreSQL | 数据库 |
| 8080 | Go Server | API 服务 |

## 常用命令

```bash
# 查看日志
docker-compose logs -f

# 重启服务
docker-compose restart

# 停止服务
docker-compose down

# 重新构建
docker-compose build --no-cache
docker-compose up -d
```

## 目录结构

```
deploy/
├── docker-compose.yml    # 服务编排
├── Dockerfile.server     # Go 服务镜像
├── Dockerfile.sfu        # SFU 服务镜像
├── nginx.conf            # Nginx 配置
├── init.sql              # 数据库初始化
├── deploy.sh             # 一键部署脚本
└── README.md             # 本文件
```
