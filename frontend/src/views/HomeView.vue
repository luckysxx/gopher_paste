<template>
  <div class="home-page">
    <section class="hero">
      <div>
        <p class="eyebrow">My Snippets</p>
        <h1>管理你发布的代码片段</h1>
        <p class="lede">像 GitHub Gist 一样浏览、搜索、编辑你自己的代码片段。</p>
      </div>
      <div class="hero-actions">
        <el-input
          v-model="keyword"
          placeholder="按标题搜索..."
          clearable
          class="search-input"
          @input="handleSearch"
        />
        <el-button type="primary" @click="router.push('/snippets/new')">新建片段</el-button>
      </div>
    </section>

    <el-card class="list-card" shadow="never">
      <template v-if="loading">
        <el-skeleton :rows="8" animated />
      </template>

      <el-result v-else-if="error" icon="warning" title="加载失败" :sub-title="error">
        <template #extra>
          <el-button type="primary" @click="fetchSnippets">重试</el-button>
        </template>
      </el-result>

      <el-empty v-else-if="filteredSnippets.length === 0" description="暂无代码片段">
        <el-button type="primary" @click="router.push('/snippets/new')">创建第一个片段</el-button>
      </el-empty>

      <el-table v-else :data="filteredSnippets" stripe>
        <el-table-column prop="title" label="标题" min-width="260">
          <template #default="{ row }">
            <el-link type="primary" @click="openSnippet(row.id)">{{ row.title }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="language" label="语言" width="140" />
        <el-table-column label="更新时间" width="220">
          <template #default="{ row }">{{ formatDate(row.updated_at) }}</template>
        </el-table-column>
        <el-table-column label="可见性" width="120">
          <template #default="{ row }">
            <el-tag size="small" :type="row.visibility === 'public' ? 'success' : 'info'">
              {{ row.visibility === 'public' ? '公开' : '私有' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-space>
              <el-button text type="primary" @click="openSnippet(row.id)">查看</el-button>
              <el-button text @click="editSnippet(row.id)">编辑</el-button>
            </el-space>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { listMySnippets, type Snippet } from '@/api/paste'

const router = useRouter()
const loading = ref(false)
const error = ref('')
const keyword = ref('')
const snippets = ref<Snippet[]>([])

const filteredSnippets = computed(() => {
  const currentKeyword = keyword.value.trim().toLowerCase()
  if (!currentKeyword) {
    return snippets.value
  }
  return snippets.value.filter((item) => item.title.toLowerCase().includes(currentKeyword))
})

const fetchSnippets = async () => {
  loading.value = true
  error.value = ''
  try {
    snippets.value = await listMySnippets()
  } catch (err) {
    snippets.value = []
    error.value = err instanceof Error ? err.message : '无法加载你的代码片段'
  } finally {
    loading.value = false
  }
}

const openSnippet = (id: string | number) => {
  router.push(`/snippets/${id}`)
}

const editSnippet = (id: string | number) => {
  router.push(`/snippets/${id}/edit`)
}

const handleSearch = () => {
  return
}

const formatDate = (value: string) => {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  return date.toLocaleString()
}

onMounted(() => {
  fetchSnippets()
})
</script>

<style scoped lang="scss">
.home-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.hero {
  background: #fff;
  border: 1px solid #e4e7ed;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;

  h1 {
    margin: 8px 0;
    font-size: 28px;
    color: #1f2d3d;
  }

  .lede {
    margin: 0;
    color: #606266;
  }
}

.eyebrow {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: 999px;
  background: rgba(64, 158, 255, 0.12);
  color: #2468c0;
  font-weight: 700;
  font-size: 12px;
}

.hero-actions {
  display: flex;
  gap: 10px;
  align-items: center;

  .search-input {
    width: 240px;
  }
}

.list-card {
  border-radius: 12px;

  :deep(.el-card__body) {
    padding: 16px;
  }
}

@media (max-width: 860px) {
  .hero {
    flex-direction: column;
  }

  .hero-actions {
    width: 100%;

    .search-input {
      width: 100%;
    }
  }
}
</style>
