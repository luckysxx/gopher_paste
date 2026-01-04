# 运行方式说明

本项目支持三种运行方式，每种方式的配置和使用场景不同。

## 方式 1：本地开发（直接运行 Go 代码）

### 适用场景
- 开发调试
- 快速测试

### 前置条件
确保本地已安装并启动：
- PostgreSQL（端口 5432）
- Redis（端口 6379）

或使用 docker-compose-infra.yaml 启动基础设施：
```bash
docker-compose -f docker-compose-infra.yaml up -d postgres redis
```

### 运行步骤

1. **设置环境变量**（可选）
```bash
cd backend
cp .env.example .env
# 编辑 .env 文件，修改配置
```

2. **运行服务**
```bash
# Paste 服务
cd backend/services/paste
go run main.go

# User 服务（另开终端）
cd backend/services/user
go run main.go
```

3. **使用默认配置**
如果不设置环境变量，会使用 config.go 中的默认值：
- DB_SOURCE: `postgres://postgres:123456@postgres:5432/gopher_paste?sslmode=disable`
- REDIS_ADDR: `gopher_redis:6379`
- 注意：主机名可能需要改为 `localhost`

---

## 方式 2：Docker Compose

### 适用场景
- 本地完整环境测试
- 模拟生产环境

### 运行步骤

1. **启动基础设施**（数据库、Redis、监控）
```bash
docker-compose -f docker-compose-infra.yaml up -d
```

2. **启动应用服务**
```bash
docker-compose up -d
```

3. **访问服务**
- Paste 服务: http://localhost:8080
- User 服务: http://localhost:8081
- 前端: http://localhost:80

4. **查看日志**
```bash
docker-compose logs -f paste-service
docker-compose logs -f user-service
```

### 配置说明
环境变量在 `docker-compose.yaml` 中定义：
```yaml
environment:
  - ENV=production
  - DB_SOURCE=postgres://luckys:123456@gopher_db:5432/gopher_paste?sslmode=disable
  - REDIS_ADDR=gopher_redis:6379
  - REDIS_PASSWORD=123456
  - JWT_SECRET=gopherpaste_secret_key
```

---

## 方式 3：Kubernetes

### 适用场景
- 生产环境部署
- 云原生环境

### 前置条件
- Kubernetes 集群（本地可用 OrbStack、Minikube 等）
- kubectl 已配置

### 运行步骤

1. **创建命名空间**
```bash
kubectl apply -f k8s/namespace.yaml
```

2. **配置 Secret（⚠️ 先修改密码）**
```bash
# 编辑 k8s/secret.yaml，修改密码的 Base64 值
kubectl apply -f k8s/secret.yaml
```

3. **应用配置**
```bash
kubectl apply -f k8s/configmap.yaml
```

4. **部署服务**
```bash
kubectl apply -f k8s/paste-service.yaml
kubectl apply -f k8s/redis.yaml
# 部署其他服务...
```

5. **查看状态**
```bash
kubectl get pods -n gopher-paste
kubectl get svc -n gopher-paste
```

6. **访问服务**
```bash
# 查看 LoadBalancer 分配的地址
kubectl get svc -n gopher-paste

# 或使用端口转发
kubectl port-forward svc/paste-service 8080:8080 -n gopher-paste
```

### 配置说明
配置分为两部分：
- **ConfigMap**（非敏感配置）: `k8s/configmap.yaml`
  - 服务器端口
  - Redis 地址
  - 日志级别等
  
- **Secret**（敏感信息）: `k8s/secret.yaml`
  - 数据库连接字符串（包含密码）
  - Redis 密码
  - JWT 密钥

---

## 配置对比

| 配置项 | 本地开发 | Docker Compose | Kubernetes |
|-------|---------|----------------|------------|
| **数据库主机** | localhost | gopher_db | gopher-postgres |
| **数据库账号** | postgres | luckys | postgres |
| **数据库密码** | 123456 | 123456 | mysupersecretpassword |
| **Redis 主机** | localhost | gopher_redis | gopher-redis |
| **Redis 密码** | 123456 | 123456 | 123456 |
| **日志输出** | stdout + 文件 | stdout（容器） | stdout（容器） |
| **配置来源** | 环境变量/默认值 | docker-compose.yaml | ConfigMap + Secret |

---

## 常见问题

### Q: 本地开发时连不上数据库
A: 检查主机名，改为 `localhost`：
```bash
export DB_SOURCE="postgres://postgres:123456@localhost:5432/gopher_paste?sslmode=disable"
```

### Q: Docker Compose 启动失败
A: 确保先启动 docker-compose-infra.yaml：
```bash
docker-compose -f docker-compose-infra.yaml up -d
```

### Q: K8s 中 Pod 启动失败
A: 检查 Secret 配置：
```bash
kubectl get secret gopher-secrets -n gopher-paste -o yaml
# 检查 db-source 字段是否存在
```

### Q: 如何修改数据库密码
**本地/Docker**: 修改环境变量中的 DB_SOURCE
**K8s**: 
1. 生成新的 Base64: `echo -n "postgres://user:newpass@host:5432/db" | base64`
2. 更新 secret.yaml 中的 db-source
3. 重新应用: `kubectl apply -f k8s/secret.yaml`
4. 重启 Pod: `kubectl rollout restart deployment paste-service-deploy -n gopher-paste`
