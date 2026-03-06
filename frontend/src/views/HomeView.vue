<template>
  <div class="home-page">
    <section class="hero">
      <div class="hero-copy">
        <p class="eyebrow">GopherPaste · 代码快贴</p>
        <h1>创建、分享、通过短链快速打开</h1>
        <p class="lede">无需登录，生成短链即可分享。正式项目样式，接口对接后即可上线。</p>

        <div class="quick-card">
          <div class="quick-title">已有短链？直接打开</div>
          <div class="quick-row">
            <el-input v-model="shortLink" placeholder="输入短链或 ID，例如 gp-8f3c1" clearable @keyup.enter="goToPaste" />
            <el-button type="primary" :disabled="!shortLink.trim()" @click="goToPaste">
              <el-icon>
                <View />
              </el-icon>
              <span>打开</span>
            </el-button>
          </div>
        </div>
      </div>

      <el-card class="side-card" shadow="hover">
        <div class="side-header">
          <el-tag type="success" effect="dark" size="small">即时</el-tag>
          <h3>快速发布一个代码贴</h3>
          <p>选择语言，粘贴代码，生成短链即可分享。</p>
        </div>

        <el-form :model="form" label-position="top" size="large">
          <el-form-item label="语言">
            <el-select v-model="form.language" placeholder="选择语言" filterable>
              <el-option label="Plain Text" value="text" />
              <el-option label="Go" value="go" />
              <el-option label="JavaScript" value="javascript" />
              <el-option label="TypeScript" value="typescript" />
              <el-option label="Python" value="python" />
              <el-option label="Java" value="java" />
              <el-option label="C++" value="cpp" />
              <el-option label="Rust" value="rust" />
              <el-option label="SQL" value="sql" />
              <el-option label="JSON" value="json" />
              <el-option label="YAML" value="yaml" />
              <el-option label="Markdown" value="markdown" />
              <el-option label="HTML" value="html" />
              <el-option label="CSS" value="css" />
              <el-option label="Shell" value="shell" />
            </el-select>
          </el-form-item>

          <el-form-item label="代码内容">
            <el-input v-model="form.content" type="textarea" :rows="12" placeholder="在此粘贴或输入代码..." resize="none"
              spellcheck="false" class="code-input" />
          </el-form-item>

          <div class="form-actions">
            <el-button type="primary" size="large" :loading="loading" @click="onSubmit" class="submit-btn">
              <el-icon>
                <Promotion />
              </el-icon>
              <span>生成短链</span>
            </el-button>
          </div>
        </el-form>
      </el-card>
    </section>

    <section class="info-grid">
      <el-card class="info-card" shadow="never">
        <h4>正式项目样式</h4>
        <p>保留真实逻辑位点，接口待接入后即可运行，无示例占位内容。</p>
      </el-card>
      <el-card class="info-card" shadow="never">
        <h4>短链直达</h4>
        <p>顶部即输即开，访问路径为 /paste/:id，便于部署后直接使用。</p>
      </el-card>
      <el-card class="info-card" shadow="never">
        <h4>安全可控</h4>
        <p>后续可扩展访问权限、过期时间、登录态等能力。</p>
      </el-card>
    </section>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { createPaste } from '@/api/paste'
import { ElMessage } from 'element-plus'
import { Promotion, View } from '@element-plus/icons-vue'

const router = useRouter()
const loading = ref(false)
const shortLink = ref('')

const form = reactive({
  content: '',
  language: 'text',
})

const goToPaste = () => {
  const id = shortLink.value.trim()
  if (!id) return
  router.push(`/paste/${id}`)
}

const onSubmit = async () => {
  if (!form.content.trim()) {
    ElMessage.warning('请输入代码内容')
    return
  }

  loading.value = true
  try {
    const res = await createPaste({
      content: form.content,
      language: form.language,
    })

    ElMessage.success('创建成功！')

    const pasteId = res.short_link

    if (pasteId) {
      router.push(`/paste/${pasteId}`)
    } else {
      console.error('No short_link returned', res)
    }
  } catch (error) {
    console.error(error)
    // request.ts 已统一处理错误提示，这里可省略
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
.home-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.hero {
  display: grid;
  grid-template-columns: 1.2fr 1fr;
  gap: 20px;
  align-items: stretch;
}

.hero-copy {
  background: linear-gradient(135deg, #e8f1ff 0%, #f9fbff 50%, #e7f7ff 100%);
  border: 1px solid #e4e7ed;
  border-radius: 16px;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 14px;

  h1 {
    margin: 0;
    font-size: 2rem;
    color: #1f2d3d;
  }

  .lede {
    margin: 0;
    color: #596273;
    line-height: 1.6;
  }
}

.eyebrow {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  background: rgba(64, 158, 255, 0.12);
  color: #2468c0;
  border-radius: 999px;
  font-weight: 700;
  width: fit-content;
  letter-spacing: 0.2px;
}

.quick-card {
  background: #fff;
  border: 1px solid #dcdfe6;
  border-radius: 12px;
  padding: 14px;
  display: flex;
  flex-direction: column;
  gap: 10px;

  .quick-title {
    font-weight: 600;
    color: #1f2d3d;
  }

  .quick-row {
    display: flex;
    gap: 10px;
    flex-wrap: wrap;

    .el-input {
      flex: 1;
      min-width: 260px;
    }
  }
}

.side-card {
  border-radius: 16px;
  border: 1px solid #e4e7ed;

  :deep(.el-card__body) {
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
}

.side-header {
  display: flex;
  flex-direction: column;
  gap: 8px;

  h3 {
    margin: 0;
  }

  p {
    margin: 0;
    color: #606266;
  }
}

.code-input {
  :deep(.el-textarea__inner) {
    font-family: 'Fira Code', 'Consolas', monospace;
    font-size: 14px;
    line-height: 1.6;
    background-color: #f8f9fa;
    border-color: #e4e7ed;
    border-radius: 8px;
    padding: 14px;

    &:focus {
      background-color: #fff;
      border-color: #409eff;
      box-shadow: 0 0 0 1px #409eff;
    }
  }
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 6px;

  .submit-btn {
    padding-left: 26px;
    padding-right: 26px;
    font-weight: 600;
    border-radius: 8px;

    .el-icon {
      margin-right: 6px;
    }
  }
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 14px;
}

.info-card {
  border-radius: 12px;
  border: 1px solid #ebeef5;

  h4 {
    margin: 0 0 6px;
  }

  p {
    margin: 0;
    color: #606266;
    line-height: 1.5;
  }
}

@media (max-width: 960px) {
  .hero {
    grid-template-columns: 1fr;
  }

  .hero-copy {
    order: 2;
  }

  .side-card {
    order: 1;
  }
}
</style>
