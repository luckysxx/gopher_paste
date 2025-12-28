<template>
  <div class="common-layout">
    <el-container class="main-container">
      <el-header class="app-header">
        <div class="header-inner">
          <div class="logo" @click="$router.push('/')">
            <el-icon class="logo-icon" :size="24">
              <EditPen />
            </el-icon>
            <span class="logo-text">GopherPaste</span>
          </div>
          <div class="nav-actions">
            <el-button text @click="$router.push('/')">首页</el-button>
            <el-button text @click="$router.push('/auth')">登录 / 注册</el-button>
            <el-button text @click="$router.push('/about')">关于</el-button>
            <a href="https://github.com" target="_blank" class="github-link">
              <el-button text>GitHub</el-button>
            </a>
          </div>
        </div>
      </el-header>

      <el-main class="app-main">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>

      <el-footer class="app-footer">
        <div class="footer-content">
          <p>© {{ new Date().getFullYear() }} GopherPaste. Built with Go & Vue 3.</p>
        </div>
      </el-footer>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { EditPen } from '@element-plus/icons-vue'
</script>

<style lang="scss">
/* Global Styles */
body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background-color: #f5f7fa;
  color: #2c3e50;
}

.common-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-container {
  min-height: 100vh;
}

.app-header {
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  position: sticky;
  top: 0;
  z-index: 100;
  padding: 0;
  height: 64px;

  .header-inner {
    max-width: 1200px;
    margin: 0 auto;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 20px;
  }

  .logo {
    display: flex;
    align-items: center;
    gap: 10px;
    cursor: pointer;
    font-weight: 700;
    font-size: 20px;
    color: #409eff;
    transition: opacity 0.2s;

    &:hover {
      opacity: 0.8;
    }
  }

  .nav-actions {
    display: flex;
    gap: 10px;

    .github-link {
      text-decoration: none;
    }
  }
}

.app-main {
  padding: 40px 20px;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
}

.app-footer {
  background-color: #fff;
  border-top: 1px solid #ebeef5;
  padding: 20px;
  text-align: center;
  color: #909399;
  font-size: 14px;
  margin-top: auto;
}

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
