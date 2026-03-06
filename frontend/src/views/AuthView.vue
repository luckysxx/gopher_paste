<template>
  <div class="auth-page">
    <div class="auth-hero">
      <p class="eyebrow">GopherPaste · 安全空间</p>
      <h1>登录或注册，开始分享代码</h1>
      <p class="lede">
        集中管理你的代码片段，随时随地查看与分享。纯 UI 预览，暂无接入真实后端。
      </p>
    </div>

    <el-card class="auth-card" shadow="always">
      <div class="card-grid">
        <section class="welcome">
          <p class="badge">快速上手</p>
          <h2>保持专注的工作流</h2>
          <ul>
            <li>多语言高亮，代码片段整洁呈现</li>
            <li>可分享链接，协作沟通更高效</li>
            <li>轻量化设计，打开即用</li>
          </ul>
          <div class="cta">立即体验全新的粘贴板体验</div>
        </section>

        <section class="forms">
          <el-tabs v-model="activeTab" stretch>
            <el-tab-pane label="登录" name="login">
              <el-form :model="loginForm" label-position="top" size="large" class="form-panel">
                <el-form-item label="邮箱或用户名">
                  <el-input v-model="loginForm.account" placeholder="you@example.com" />
                </el-form-item>
                <el-form-item label="密码">
                  <el-input v-model="loginForm.password" type="password" placeholder="••••••••" show-password />
                </el-form-item>
                <div class="form-actions">
                  <el-checkbox v-model="loginForm.remember">保持登录</el-checkbox>
                  <el-button link type="primary">忘记密码？</el-button>
                </div>
                <el-button type="primary" class="full" size="large" @click="handleLogin">
                  <el-icon>
                    <Right />
                  </el-icon>
                  <span>进入控制台</span>
                </el-button>
              </el-form>
            </el-tab-pane>

            <el-tab-pane label="注册" name="register">
              <el-form :model="registerForm" label-position="top" size="large" class="form-panel">
                <el-form-item label="邮箱">
                  <el-input v-model="registerForm.email" placeholder="you@example.com" />
                </el-form-item>
                <el-form-item label="用户名">
                  <el-input v-model="registerForm.username" placeholder="起一个有辨识度的名字" />
                </el-form-item>
                <el-form-item label="密码">
                  <el-input v-model="registerForm.password" type="password" placeholder="至少 8 位，包含数字" show-password />
                </el-form-item>
                <el-form-item label="确认密码">
                  <el-input v-model="registerForm.confirm" type="password" placeholder="再次输入密码" show-password />
                </el-form-item>
                <el-button type="success" class="full" size="large" @click="handleRegister">
                  <el-icon>
                    <CircleCheck />
                  </el-icon>
                  <span>创建账号</span>
                </el-button>
              </el-form>
            </el-tab-pane>
          </el-tabs>

          <div class="divider">
            <span>或使用快捷方式</span>
          </div>

          <div class="socials">
            <el-button plain class="social-btn" @click="informPreview">
              <el-icon>
                <Message />
              </el-icon>
              <span>邮件登录</span>
            </el-button>
            <el-button plain class="social-btn" @click="informPreview">
              <el-icon>
                <ChatDotRound />
              </el-icon>
              <span>企业单点</span>
            </el-button>
            <el-button plain class="social-btn" @click="informPreview">
              <el-icon>
                <Key />
              </el-icon>
              <span>临时访问码</span>
            </el-button>
          </div>
        </section>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ChatDotRound, CircleCheck, Key, Message, Right } from '@element-plus/icons-vue'
import { login, register } from '@/api/user'

const router = useRouter()
const activeTab = ref('login')

const loginForm = reactive({
  account: '',
  password: '',
  remember: true,
})

const registerForm = reactive({
  email: '',
  username: '',
  password: '',
  confirm: '',
})

const informPreview = () => {
  ElMessage.info('当前为界面示例，尚未接入实际认证逻辑')
}

const getErrorMessage = (error: unknown, fallback: string) => {
  if (error instanceof Error && error.message) return error.message
  return fallback
}

const handleLogin = async () => {
  if (!loginForm.account || !loginForm.password) {
    ElMessage.warning('请填写完整的登录信息')
    return
  }

  try {
    const res = await login({
      username: loginForm.account,
      password: loginForm.password,
    })
    
    // 保存 token 到 localStorage
    localStorage.setItem('token', res.token)
    localStorage.setItem('user', JSON.stringify({
      id: res.user_id,
      username: res.username,
      email: res.email,
    }))
    
    ElMessage.success('登录成功！')
    router.push('/')
  } catch (error: unknown) {
    ElMessage.error(getErrorMessage(error, '登录失败，请检查用户名和密码'))
  }
}

const handleRegister = async () => {
  if (!registerForm.email || !registerForm.username || !registerForm.password) {
    ElMessage.warning('请填写完整的注册信息')
    return
  }
  
  if (registerForm.password !== registerForm.confirm) {
    ElMessage.warning('两次输入的密码不一致')
    return
  }

  try {
    const res = await register({
      username: registerForm.username,
      email: registerForm.email,
      password: registerForm.password,
    })
    
    ElMessage.success(`注册成功！欢迎你，${res.username}`)
    // 注册成功后切换到登录页
    activeTab.value = 'login'
    loginForm.account = registerForm.username
  } catch (error: unknown) {
    ElMessage.error(getErrorMessage(error, '注册失败，请稍后重试'))
  }
}
</script>

<style scoped lang="scss">
.auth-page {
  min-height: calc(100vh - 160px);
  background: radial-gradient(circle at 10% 20%, rgba(64, 158, 255, 0.08), transparent 35%),
    radial-gradient(circle at 90% 10%, rgba(103, 194, 58, 0.08), transparent 30%),
    linear-gradient(135deg, #f6f9ff 0%, #fdfdfd 60%, #f6fbff 100%);
  border-radius: 16px;
  padding: 32px 20px 48px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.auth-hero {
  max-width: 860px;
  margin: 0 auto;
  text-align: center;
  display: flex;
  flex-direction: column;
  gap: 10px;

  .eyebrow {
    display: inline-flex;
    align-self: center;
    align-items: center;
    gap: 8px;
    padding: 6px 14px;
    background: rgba(64, 158, 255, 0.1);
    color: #409eff;
    border-radius: 999px;
    font-weight: 700;
    letter-spacing: 0.2px;
  }

  h1 {
    margin: 0;
    font-size: 2.6rem;
    font-weight: 800;
    color: #1f2d3d;
  }

  .lede {
    margin: 0;
    color: #5f6b7a;
    font-size: 1.05rem;
  }
}

.auth-card {
  border: none;
  border-radius: 18px;
  box-shadow: 0 20px 60px rgba(31, 45, 61, 0.08);
  padding: 6px;

  :deep(.el-card__body) {
    padding: 0;
  }
}

.card-grid {
  display: grid;
  grid-template-columns: 1fr 1.2fr;
  gap: 0;
  overflow: hidden;
  border-radius: 14px;
  background: #fff;
}

.welcome {
  background: linear-gradient(180deg, #1d2b64 0%, #1e3c72 60%, #15244e 100%);
  color: #e8f1ff;
  padding: 32px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  position: relative;

  &::after {
    content: '';
    position: absolute;
    inset: 0;
    background: radial-gradient(circle at 30% 20%, rgba(255, 255, 255, 0.08), transparent 40%);
    pointer-events: none;
  }

  .badge {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    background: rgba(255, 255, 255, 0.1);
    padding: 6px 12px;
    border-radius: 999px;
    width: fit-content;
    font-weight: 700;
  }

  h2 {
    margin: 0;
    font-size: 1.8rem;
    line-height: 1.2;
  }

  ul {
    padding-left: 20px;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 10px;
    color: #d6e4ff;
    line-height: 1.5;
  }

  .cta {
    margin-top: auto;
    padding: 12px 14px;
    border-radius: 12px;
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.14);
    font-weight: 700;
    text-align: center;
  }
}

.forms {
  padding: 28px 32px;
  display: flex;
  flex-direction: column;
  gap: 18px;

  .form-panel {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .form-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin: 6px 0 16px;
  }

  .full {
    width: 100%;
    border-radius: 10px;
    font-weight: 700;
  }
}

.divider {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #a0a8b3;
  font-size: 0.95rem;
  margin-top: -6px;

  &::before,
  &::after {
    content: '';
    flex: 1;
    height: 1px;
    background: linear-gradient(90deg, transparent, #e4e7ed 50%, transparent);
  }
}

.socials {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 10px;

  .social-btn {
    justify-content: flex-start;
    width: 100%;
    border-radius: 10px;
    font-weight: 600;
    color: #3a4250;

    .el-icon {
      margin-right: 8px;
    }
  }
}

@media (max-width: 960px) {
  .card-grid {
    grid-template-columns: 1fr;
  }

  .welcome {
    border-bottom: 1px solid rgba(255, 255, 255, 0.12);
  }
}
</style>
