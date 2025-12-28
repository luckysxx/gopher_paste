<template>
  <div class="about-container">
    <el-card class="about-card">
      <h1>关于 GopherPaste</h1>
      <p>
        GopherPaste 是一个基于 Go (Gin) 和 Vue 3 (Element Plus) 构建的现代化代码粘贴板工具。
      </p>
      <p>
        旨在提供快速、简洁的代码分享体验。
      </p>

      <el-divider />

      <h3>技术栈</h3>
      <ul>
        <li><strong>Frontend:</strong> Vue 3, TypeScript, Vite, Element Plus, Pinia</li>
        <li><strong>Backend:</strong> Go, Gin, GORM</li>
      </ul>
    </el-card>
  </div>
</template>

<style scoped>
.about-container {
  max-width: 800px;
  margin: 0 auto;
}

.about-card {
  padding: 20px;
  border-radius: 12px;
}

h1 {
  margin-top: 0;
  color: #409eff;
}

ul {
  padding-left: 20px;
  line-height: 1.8;
}
</style>
