# GopherPaste

一个基于 Go + Vue 的 Pastebin 风格代码分享平台，采用前后端分离与双后端服务（`paste-service`、`user-service`）架构，支持本地开发、Docker Compose 和 Kubernetes 部署。

## 项目特性

- 用户注册、登录与 JWT 鉴权
- 代码片段创建与查询（Paste 服务）
- 前端基于 Vue 3 + Vite + Element Plus
- 基础设施包含 PostgreSQL、Redis
- 可选监控与日志组件：Prometheus、Grafana、Loki、Promtail

## 技术栈

- 后端：Go、Gin、PostgreSQL、Redis、sqlc、JWT
- 前端：Vue 3、TypeScript、Vite、Pinia、Vue Router、Element Plus
- 部署：Docker Compose、Kubernetes

## 目录说明

```text
backend/                 Go 后端
  common/                公共组件（配置、鉴权、数据库、日志等）
  services/paste/        Paste 服务（默认 8080）
  services/user/         User 服务（默认 8081）

frontend/                Vue 前端（本地 dev 默认 5173）
k8s/                     Kubernetes 部署清单
docker-compose*.yaml     Compose 运行文件
RUNNING.md               详细运行说明
Makefile                 常用开发命令
```

## 环境要求

- Go `1.25+`
- Node.js `20.19+`（或 `22.12+`）
- pnpm
- Docker / Docker Compose（用于基础设施或容器化运行）

## 快速开始（推荐）

在项目根目录执行：

```bash
make dev
```

该命令会：

1. 启动本地基础设施（PostgreSQL、Redis）
2. 启动 `paste-service`（8080）
3. 启动 `user-service`（8081）
4. 启动前端开发服务器（Vite）

首次启动前建议先初始化数据库：

```bash
make db-init
```

## 运行方式

### 1) 本地开发（Go 直接运行）

```bash
# 启动基础设施
make dev-infra

# 初始化数据库（首次）
make db-init

# 分别启动服务（多终端）
cd backend/services/paste && go run main.go
cd backend/services/user && SERVER_PORT=8081 go run main.go
cd frontend && pnpm install && pnpm dev
```

### 2) Docker Compose

```bash
# 先启动基础设施
docker compose -f docker-compose-infra.yaml up -d

# 首次初始化数据库
make db-init

# 启动应用服务
docker compose -f docker-compose.local.yaml up -d
```

访问地址：

- 前端：http://localhost:80
- Paste 服务：http://localhost:8080
- User 服务：http://localhost:8081

### 3) Kubernetes

参考 `k8s/` 下清单按顺序部署：

```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/redis.yaml
kubectl apply -f k8s/paste-service.yaml
```

## 常用命令

```bash
make help             # 查看全部命令
make db-init          # 初始化 users/pastes 表
make db-migrate-paste # 执行 paste 表迁移
make swagger-paste    # 生成 paste 服务 Swagger
make lint             # 后端测试 + 前端 lint
make stop-infra       # 停止 postgres/redis
```

## 服务与端口

- `paste-service`: `8080`
- `user-service`: `8081`
- `frontend-service`: `80`（容器模式）
- `postgres`: `5432`
- `redis`: `6379`
- `prometheus`: `9090`
- `grafana`: `3000`
- `loki`: `3100`

## 额外说明

- 更完整的运行与排障文档见 `RUNNING.md`
- `pg_data/`、`redis_data/` 是本地持久化目录，建议仅用于开发环境
- 若前端请求失败，优先检查前端代理配置（`frontend/vite.config.ts` 与 `frontend/nginx.conf`）
