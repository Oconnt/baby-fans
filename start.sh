#!/bin/bash

# Baby-Fans 一键启动脚本

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 正在启动 Baby-Fans 全栈服务...${NC}"
echo -e "${BLUE}🖥️  后端服务将运行在: http://localhost:18081${NC}"

# 1. 启动后端
echo -e "${GREEN}📂 正在后台启动 Go 后端服务...${NC}"
cd backend
go run cmd/server/main.go > backend.log 2>&1 &
BACKEND_PID=$!
cd ..

# 2. 启动前端 (默认启动 H5 预览模式)
echo -e "${GREEN}🌐 正在启动 UniApp 前端服务 (H5 模式)...${NC}"
echo -e "${BLUE}💡 如果需要调试微信小程序，请手动在 frontend 目录下运行: npm run dev:mp-weixin${NC}"
cd frontend
npm run dev:h5 &
FRONTEND_PID=$!
cd ..

# 捕获退出信号
trap "echo -e '${BLUE}🛑 正在停止所有服务...${NC}'; kill $BACKEND_PID $FRONTEND_PID; exit" SIGINT SIGTERM

echo -e "${GREEN}✅ 服务已就绪！${NC}"
echo -e "🖥️  后端 API: http://localhost:18081"
echo -e "📱 前端 H5:  请查看上方 Vite 输出的地址"
echo -e "⌨️  按 Ctrl+C 可同时停止前后端服务"

# 保持脚本运行以等待信号
wait
