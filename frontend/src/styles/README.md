# EMS 前端样式指南

本文档介绍 EMS 设备管理系统的前端设计系统和样式规范。

## 📁 样式文件结构

```
frontend/src/styles/
├── design-system.css    # 核心设计系统（CSS变量、主题）
├── utilities.css        # 工具类（快速构建界面）
├── h5.css              # 移动端H5样式
├── pc.css              # PC端样式优化
└── README.md           # 本文档
```

## 🎨 设计系统

### 配色方案

#### 主色调
- **Primary**: `#667eea` - 工业蓝（主要操作、强调）
- **Hover**: `#7c8ff5`
- **Dim**: `rgba(102, 126, 234, 0.15)`

#### 状态色
- **Success**: `#00D2A0` / `#00FFA3` - 成功/正常
- **Warning**: `#FFB800` - 警告/进行中
- **Danger**: `#FF6B6B` / `#FF4757` - 危险/错误
- **Info**: `#5C7CFA` - 信息提示

#### 间距系统
- `--space-xs`: 4px
- `--space-sm`: 8px
- `--space-md`: 16px
- `--space-lg`: 24px
- `--space-xl`: 32px
- `--space-2xl`: 48px

#### 圆角
- `--radius-sm`: 4px
- `--radius-md`: 8px
- `--radius-lg`: 12px
- `--radius-xl`: 16px
- `--radius-2xl`: 24px

#### 过渡动画
- `--transition-fast`: 150ms
- `--transition-base`: 250ms
- `--transition-slow`: 350ms

### 主题支持

系统支持深色/浅色主题切换，通过 `data-theme` 属性控制：

```html
<div data-theme="light">...</div>
<div data-theme="dark">...</div>
```

## 🛠 工具类使用

### 间距类
```html
<div class="p-4">内边距 24px</div>
<div class="mb-3">下边距 16px</div>
<div class="px-2 py-3">水平8px 垂直16px</div>
```

### 文字类
```html
<div class="text-lg font-semibold">大号加粗文字</div>
<div class="text-sm text-secondary">小号次要文字</div>
```

### Flexbox 类
```html
<div class="flex items-center justify-between">
  <span>左边</span>
  <span>右边</span>
</div>
```

### 卡片类
```html
<div class="card">标准卡片</div>
<div class="card-elevated">高阴影卡片</div>
```

### 状态标签
```html
<span class="badge badge-success">成功</span>
<span class="badge badge-warning">警告</span>
<span class="badge badge-danger">危险</span>
```

## 📱 移动端 H5 样式

### 页面结构
```html
<div class="h5-page">
  <div class="h5-header">...</div>
  <div class="content">...</div>
</div>
```

### 卡片样式
```html
<div class="h5-card">标准卡片</div>
<div class="h5-card-elevated">高阴影卡片</div>
```

### 按钮样式
```html
<button class="h5-btn h5-btn-primary">主要按钮</button>
<button class="h5-btn h5-btn-success">成功按钮</button>
<button class="h5-btn h5-btn-warning">警告按钮</button>
<button class="h5-btn h5-btn-danger">危险按钮</button>
```

### 列表项
```html
<div class="h5-list-item">
  <div class="h5-list-item-icon">📋</div>
  <div class="h5-list-item-content">
    <div class="h5-list-item-title">标题</div>
    <div class="h5-list-item-desc">描述</div>
  </div>
</div>
```

### 快捷操作网格
```html
<div class="h5-action-grid">
  <div class="h5-action-item">
    <div class="h5-action-icon">📋</div>
    <span class="h5-action-label">点检</span>
  </div>
</div>
```

## 💻 PC 端样式

### 页面容器
```html
<div class="page-container">
  <div class="page-header">
    <h1 class="page-title">
      <span class="page-title-icon">📋</span>
      页面标题
    </h1>
    <div class="page-actions">...</div>
  </div>
</div>
```

### 统计卡片
使用 `el-statistic` 组件，样式已自动优化：
- 渐变顶部边框（hover时显示）
- 悬停抬升效果
- 统一圆角和阴影

### 表格样式
- 表头：大写字母、字间距、底部加粗边框
- 行悬停：半透明背景色
- 分页：圆角按钮、渐变激活色

### 按钮样式
所有按钮已应用渐变效果：
- Primary: `#667eea` → `#764ba2`
- Success: `#00D2A0` → `#00A67A`
- Warning: `#FFB800` → `#FF9500`
- Danger: `#FF6B6B` → `#EE5A5A`

## 🎯 最佳实践

### 1. 使用 CSS 变量
```css
/* ✅ 推荐 */
color: var(--color-text-primary);
background: var(--color-bg-card);

/* ❌ 避免 */
color: #1A202C;
background: #FFFFFF;
```

### 2. 使用工具类
```html
<!-- ✅ 推荐 -->
<div class="flex items-center justify-between gap-3 p-4">

<!-- ❌ 避免 -->
<div style="display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 16px;">
```

### 3. 响应式设计
```css
/* 移动端优先，使用媒体查询 */
@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
```

### 4. 主题适配
```css
/* 使用 prefers-color-scheme 自动适配暗色模式 */
@media (prefers-color-scheme: dark) {
  .component {
    background: var(--color-bg-card-dark);
    color: var(--color-text-primary-dark);
  }
}
```

## 🔄 更新日志

### 2026-01-17
- 新增 `utilities.css` 工具类文件
- 新增 `h5.css` 移动端样式
- 新增 `pc.css` PC端优化样式
- 优化 H5 TasksView 页面
- 优化 H5 IndexView 页面
- 优化 MobileHeader 组件
- 优化 PC端 TaskListView 页面

## 📚 参考资源

- [Element Plus 文档](https://element-plus.org/)
- [Vant 文档](https://vant-ui.github.io/)
- [CSS 变量 MDN](https://developer.mozilla.org/zh-CN/docs/Web/CSS/Using_CSS_custom_properties)
