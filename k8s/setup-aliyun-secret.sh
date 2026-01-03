#!/bin/bash
# ==========================================
# 配置阿里云镜像拉取凭证
# ==========================================
# 用途：交互式创建 Kubernetes imagePullSecrets

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${GREEN}=========================================="
echo "🔐 配置阿里云镜像拉取凭证"
echo -e "==========================================${NC}\n"

# 固定的镜像仓库地址
DOCKER_SERVER="crpi-owej7l76wqheszs7.cn-hangzhou.personal.cr.aliyuncs.com"
NAMESPACE="gopher-paste"

echo -e "${BLUE}📝 镜像仓库地址：${NC}$DOCKER_SERVER"
echo -e "${BLUE}📦 命名空间：${NC}$NAMESPACE"
echo ""

# ==========================================
# 检查命名空间是否存在
# ==========================================
echo -e "${YELLOW}📋 检查命名空间...${NC}"
if ! kubectl get namespace $NAMESPACE &> /dev/null; then
    echo -e "${YELLOW}  ⚠️  命名空间不存在，正在创建...${NC}"
    kubectl apply -f namespace.yaml
    echo -e "${GREEN}  ✅ 命名空间创建成功${NC}"
else
    echo -e "${GREEN}  ✅ 命名空间已存在${NC}"
fi

echo ""

# ==========================================
# 检查 Secret 是否已存在
# ==========================================
echo -e "${YELLOW}🔍 检查现有凭证...${NC}"
if kubectl get secret aliyun-secret -n $NAMESPACE &> /dev/null; then
    echo -e "${YELLOW}  ⚠️  凭证已存在！${NC}"
    read -p "  是否删除并重新创建？(y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        kubectl delete secret aliyun-secret -n $NAMESPACE
        echo -e "${GREEN}  ✅ 已删除旧凭证${NC}"
    else
        echo -e "${BLUE}  ℹ️  保留现有凭证，退出${NC}"
        exit 0
    fi
fi

echo ""

# ==========================================
# 提示用户如何获取凭证
# ==========================================
echo -e "${BLUE}=========================================="
echo "📖 如何获取阿里云凭证？"
echo -e "==========================================${NC}"
echo "1️⃣  访问：https://cr.console.aliyun.com/"
echo "2️⃣  点击左侧「访问凭证」"
echo "3️⃣  查看或重置密码"
echo ""
echo -e "${YELLOW}💡 提示：用户名和密码与 GitHub Secrets 中的相同${NC}"
echo ""

# ==========================================
# 输入凭证
# ==========================================
echo -e "${GREEN}=========================================="
echo "🔑 请输入凭证"
echo -e "==========================================${NC}"

# 输入用户名
read -p "阿里云用户名: " DOCKER_USERNAME
if [ -z "$DOCKER_USERNAME" ]; then
    echo -e "${RED}❌ 用户名不能为空${NC}"
    exit 1
fi

# 输入密码（不显示）
read -s -p "阿里云密码: " DOCKER_PASSWORD
echo
if [ -z "$DOCKER_PASSWORD" ]; then
    echo -e "${RED}❌ 密码不能为空${NC}"
    exit 1
fi

echo ""

# ==========================================
# 创建 Secret
# ==========================================
echo -e "${YELLOW}🚀 创建镜像拉取凭证...${NC}"

kubectl create secret docker-registry aliyun-secret \
  --docker-server=$DOCKER_SERVER \
  --docker-username=$DOCKER_USERNAME \
  --docker-password=$DOCKER_PASSWORD \
  --namespace=$NAMESPACE

if [ $? -eq 0 ]; then
    echo -e "${GREEN}=========================================="
    echo "✅ 凭证创建成功！"
    echo -e "==========================================${NC}\n"
    
    # 验证
    echo -e "${YELLOW}🔍 验证凭证...${NC}"
    kubectl get secret aliyun-secret -n $NAMESPACE
    
    echo ""
    echo -e "${GREEN}🎉 完成！现在可以运行 ./deploy.sh 部署应用了${NC}"
else
    echo -e "${RED}=========================================="
    echo "❌ 创建失败！"
    echo -e "==========================================${NC}"
    exit 1
fi
