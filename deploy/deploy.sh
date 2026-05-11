#!/bin/bash

set -e

echo "=========================================="
echo "  AFT集团线上服务平台 - 一键部署"
echo "=========================================="

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo "错误: 请先安装 Docker"
    echo "访问: https://docs.docker.com/get-docker/"
    exit 1
fi

if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "错误: 请先安装 Docker Compose"
    exit 1
fi

# 获取服务器 IP
if command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE="docker-compose"
else
    DOCKER_COMPOSE="docker compose"
fi

echo ""
echo "请设置 SFU 服务器公网 IP（留空则使用 127.0.0.1）:"
read -p ">>> " SFU_IP
SFU_IP=${SFU_IP:-127.0.0.1}

echo ""
echo "是否配置微信登录? (y/N)"
read -p ">>> " CONFIG_WECHAT

if [ "$CONFIG_WECHAT" = "y" ] || [ "$CONFIG_WECHAT" = "Y" ]; then
    echo "请输入微信 AppID:"
    read -p ">>> " WECHAT_APP_ID
    echo "请输入微信 AppSecret:"
    read -p ">>> " WECHAT_APP_SECRET
    export WECHAT_APP_ID
    export WECHAT_APP_SECRET
fi

export SFU_IP

echo ""
echo "正在构建和启动服务..."

# 构建镜像
echo "1. 构建 Docker 镜像..."
$DOCKER_COMPOSE build

# 启动服务
echo "2. 启动服务..."
$DOCKER_COMPOSE up -d

# 等待服务启动
echo "3. 等待服务就绪..."
sleep 10

# 检查服务状态
echo "4. 检查服务状态..."
$DOCKER_COMPOSE ps

echo ""
echo "=========================================="
echo "  部署完成!"
echo "=========================================="
echo ""
echo "访问地址: http://你的服务器IP/"
echo ""
echo "服务状态检查:"
echo "  - 前端: http://你的服务器IP/"
echo "  - API:  http://你的服务器IP/api/health"
echo "  - SFU:  http://你的服务器IP:3000/health"
echo ""
echo "查看日志: docker-compose logs -f"
echo "停止服务: docker-compose down"
echo ""
