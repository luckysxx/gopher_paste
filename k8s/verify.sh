#!/bin/bash
# ==========================================
# 配置验证脚本
# ==========================================
# 用途：验证 ConfigMap 和 Secret 是否正确配置

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}=========================================="
echo "🔍 验证 ConfigMap 和 Secret 配置"
echo -e "==========================================${NC}\n"

# 检查 ConfigMap
echo -e "${YELLOW}1️⃣  检查 ConfigMap...${NC}"
if kubectl get configmap gopher-config &> /dev/null; then
    echo -e "${GREEN}   ✅ ConfigMap 存在${NC}"
else
    echo -e "${RED}   ❌ ConfigMap 不存在${NC}"
fi

# 检查 Secret
echo -e "\n${YELLOW}2️⃣  检查 Secret...${NC}"
if kubectl get secret gopher-secrets &> /dev/null; then
    echo -e "${GREEN}   ✅ Secret 存在${NC}"
else
    echo -e "${RED}   ❌ Secret 不存在${NC}"
fi

echo -e "\n${GREEN}✅ 验证完成${NC}"
