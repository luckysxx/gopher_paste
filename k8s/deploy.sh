#!/bin/bash
# ==========================================
# Kubernetes 一键部署脚本
# ==========================================
# 用途：按正确顺序部署所有 K8s 资源
# 使用：chmod +x deploy.sh && ./deploy.sh

set -e  # 遇到错误立即退出

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}=========================================="
echo "🚀 GopherPaste K8s 部署脚本"
echo -e "==========================================${NC}\n"

# ==========================================
# 第 1 步：检查 kubectl 是否可用
# ==========================================
echo -e "${YELLOW}📋 检查依赖...${NC}"
if ! command -v kubectl &> /dev/null; then
    echo -e "${RED}❌ kubectl 未安装，请先安装 kubectl${NC}"
    exit 1
fi

# 检查 K8s 集群连接
if ! kubectl cluster-info &> /dev/null; then
    echo -e "${RED}❌ 无法连接到 K8s 集群，请检查配置${NC}"
    exit 1
fi

echo -e "${GREEN}✅ kubectl 可用${NC}"
echo -e "${GREEN}✅ K8s 集群连接正常${NC}\n"

# ==========================================
# 第 2 步：创建命名空间
# ==========================================
echo -e "${YELLOW}🏗️  创建命名空间...${NC}"

# 创建 gopher-paste 命名空间
kubectl apply -f namespace.yaml

echo -e "${GREEN}✅ 命名空间创建完成${NC}\n"

# ==========================================
# 第 3 步：配置镜像拉取凭证
# ==========================================
echo -e "${YELLOW}🔐 配置镜像拉取凭证...${NC}"

# 检查是否已存在 aliyun-secret
if kubectl get secret aliyun-secret -n gopher-paste &> /dev/null; then
    echo -e "${GREEN}  ✅ 镜像拉取凭证已存在${NC}"
else
    echo -e "${YELLOW}  ⚠️  镜像拉取凭证不存在${NC}"
    echo -e "${YELLOW}  📝 请手动创建（需要你的阿里云账号信息）：${NC}"
    echo ""
    echo "  kubectl create secret docker-registry aliyun-secret \\"
    echo "    --docker-server=crpi-owej7l76wqheszs7.cn-hangzhou.personal.cr.aliyuncs.com \\"
    echo "    --docker-username=你的阿里云用户名 \\"
    echo "    --docker-password=你的阿里云密码 \\"
    echo "    --namespace=gopher-paste"
    echo ""
    echo -e "${YELLOW}  或者直接跳过（如果镜像是公开的）${NC}"
    read -p "  是否继续部署？(y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo ""

# ==========================================
# 第 4 步：部署 ConfigMap 和 Secret
# ==========================================
echo -e "${YELLOW}📦 部署配置资源...${NC}"

# 检查 Secret 文件是否存在
if [ ! -f "secret.yaml" ]; then
    echo -e "${RED}❌ secret.yaml 不存在！${NC}"
    echo -e "${YELLOW}请执行：cp secret.yaml.example secret.yaml${NC}"
    echo -e "${YELLOW}然后修改 secret.yaml 中的密码${NC}"
    exit 1
fi

# 部署 ConfigMap
echo "  → 创建 ConfigMap..."
kubectl apply -f configmap.yaml

# 部署 Secret
echo "  → 创建 Secret..."
kubectl apply -f secret.yaml

echo -e "${GREEN}✅ 配置资源部署完成${NC}\n"

# ==========================================
# 第 3 步：部署基础设施（Redis）
# ==========================================
echo -e "${YELLOW}🗄️  部署基础设施...${NC}"

# 部署 Redis
echo "  → 部署 Redis..."
kubectl apply -f redis.yaml

# 等待 Redis 就绪
echo "  → 等待 Redis 启动..."
kubectl wait --for=condition=ready pod -l app=redis -n gopher-paste --timeout=60s || {
    echo -e "${RED}❌ Redis 启动失败，请检查日志：${NC}"
    echo "    kubectl logs -l app=redis -n gopher-paste"
    exit 1
}

echo -e "${GREEN}✅ Redis 部署完成${NC}\n"

# ==========================================
# 第 4 步：部署应用服务
# ==========================================
echo -e "${YELLOW}🚀 部署应用服务...${NC}"

# 部署 Paste Service
echo "  → 部署 Paste Service..."
kubectl apply -f paste-service.yaml

# 等待 Paste Service 就绪
echo "  → 等待 Paste Service 启动..."
kubectl wait --for=condition=ready pod -l app=paste-service -n gopher-paste --timeout=120s || {
    echo -e "${RED}❌ Paste Service 启动失败，请检查日志：${NC}"
    echo "    kubectl logs -l app=paste-service -n gopher-paste"
    exit 1
}

echo -e "${GREEN}✅ Paste Service 部署完成${NC}\n"

# ==========================================
# 第 5 步：显示部署结果
# ==========================================
echo -e "${GREEN}=========================================="
echo "🎉 部署完成！"
echo -e "==========================================${NC}\n"

# 显示所有资源
echo -e "${YELLOW}📊 资源状态：${NC}"
kubectl get pods,svc,configmap,secret -n gopher-paste

# 获取服务访问地址
echo -e "\n${YELLOW}🌐 服务访问地址：${NC}"
PASTE_SERVICE_URL=$(kubectl get svc paste-service -n gopher-paste -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "localhost")
PASTE_SERVICE_PORT=$(kubectl get svc paste-service -n gopher-paste -o jsonpath='{.spec.ports[0].port}')

echo "  Paste Service: http://${PASTE_SERVICE_URL}:${PASTE_SERVICE_PORT}"

# 显示有用的命令
echo -e "\n${YELLOW}📚 常用命令：${NC}"
echo "  查看 Pod 状态：    kubectl get pods -n gopher-paste"
echo "  查看日志：         kubectl logs -f deployment/paste-service-deploy -n gopher-paste"
echo "  进入容器：         kubectl exec -it deployment/paste-service-deploy -n gopher-paste -- sh"
echo "  删除所有资源：     kubectl delete namespace gopher-paste"
echo "  重启服务：         kubectl rollout restart deployment/paste-service-deploy -n gopher-paste"

echo -e "\n${GREEN}✨ 部署成功！${NC}"
