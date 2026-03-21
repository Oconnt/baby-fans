# Baby-Fans Backend

家长积分管理系统后端服务，基于 Go + Gin 框架。

## 技术栈

- Go 1.25
- Gin Web Framework
- MySQL 8.0 / SQLite (备用)
- JWT Authentication

## 快速开始

### 本地运行

```bash
# 运行开发模式
make dev

# 或先编译再运行
make run
```

### Docker 运行（推荐）

```bash
# 使用 docker-compose (带 MySQL)
make docker-compose-up

# 或一键启动 (自动创建 MySQL 容器)
make docker-run-mysql
```

## 配置

配置文件位于 `etc/config.yaml`：

```yaml
server:
  port: "18081"
  mode: "release"

db:
  type: "mysql"           # mysql 或 sqlite
  host: "baby-fans-db"    # MySQL 容器名
  port: 3306
  username: "root"
  password: "baby_fans_password"
  name: "baby_fans"

wechat:
  app_id: ""
  app_secret: ""

jwt:
  secret: "your_jwt_secret"
  expire: 24
```

### 环境变量

可通过环境变量覆盖配置：

```bash
export DB_TYPE=mysql
export DB_HOST=baby-fans-db
export DB_PASSWORD=your_password
```

## Makefile 命令

| 命令 | 说明 |
|------|------|
| `make build` | 编译应用程序 |
| `make run` | 编译并运行 |
| `make dev` | 开发模式运行 |
| `make test` | 运行测试 |
| `make migrate` | 运行数据库迁移 |
| `make clean` | 清理构建产物 |
| `make docker-build` | 构建 Docker 镜像 |
| `make docker-run` | 运行容器 (SQLite 模式) |
| `make docker-run-mysql` | 运行容器 + MySQL (一键启动) |
| `make docker-stop` | 停止所有容器 |
| `make docker-logs` | 查看日志 |
| `make docker-compose-up` | 使用 docker-compose 启动 |
| `make docker-compose-down` | 使用 docker-compose 停止 |
| `make docker-shell` | 进入后端容器 |
| `make docker-db-shell` | 进入 MySQL 容器 |

## 架构说明

- **docker-compose**: 推荐的生产部署方式，MySQL 作为独立容器
- **docker-run-mysql**: 一键启动，自动创建 MySQL 容器和网络
- **docker-run**: 独立运行，使用内置 SQLite（适合开发调试）

迁移脚本是幂等的，可安全重复执行。

## API 端口

默认端口：`18081`
