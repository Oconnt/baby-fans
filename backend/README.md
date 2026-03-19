# Baby-Fans Backend

家长积分管理系统后端服务，基于 Go + Gin 框架。

## 技术栈

- Go 1.25
- Gin Web Framework
- SQLite Database
- JWT Authentication

## 快速开始

### 本地运行

```bash
# 运行开发模式
make dev

# 或先编译再运行
make run
```

### Docker 运行

```bash
# 构建 Docker 镜像
make docker-build

# 运行容器
make docker-run

# 使用 docker-compose
make docker-compose-up
```

## 配置

配置文件位于 `config/config.yaml`，支持以下配置项：

```yaml
server:
  port: "18081"
  mode: "release"

db:
  dsn: "baby-fans.db"

wechat:
  app_id: "your_app_id"
  app_secret: "your_app_secret"

jwt:
  secret: "your_jwt_secret"
  expire: 24
```

### 环境变量

可以通过环境变量覆盖配置：

```bash
export SERVER_PORT=18081
export JWT_SECRET=your_secret
```

## Makefile 命令

| 命令 | 说明 |
|------|------|
| `make build` | 编译应用程序 |
| `make run` | 编译并运行 |
| `make dev` | 开发模式运行 |
| `make test` | 运行测试 |
| `make clean` | 清理构建产物 |
| `make docker-build` | 构建 Docker 镜像 |
| `make docker-run` | 运行 Docker 容器 |
| `make docker-stop` | 停止 Docker 容器 |
| `make docker-logs` | 查看容器日志 |
| `make docker-compose-up` | 使用 docker-compose 启动 |
| `make docker-compose-down` | 使用 docker-compose 停止 |
| `make docker-shell` | 进入容器 shell |

## API 端口

默认端口：`18081`

## 目录结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go          # 入口文件
├── config/
│   ├── config.go            # 配置加载
│   └── config.yaml         # 配置文件
├── internal/
│   ├── api/
│   │   ├── handler/        # HTTP 处理器
│   │   ├── middleware/     # 中间件
│   │   └── router.go       # 路由配置
│   ├── model/              # 数据模型
│   ├── repository/         # 数据库操作
│   └── service/           # 业务逻辑
├── storage/                # 文件存储
├── Dockerfile
├── Makefile
├── docker-compose.yml
└── go.mod
```
