# EMS 前端样式指南

本文档介绍 EMS 设备管理系统的前端设计系统和样式规范。

## 设计理念

本系统采用 **Claude 风格**设计语言，以暖陶土色调为主，营造知性、温暖的编辑布局体验。

- **温暖知性**：使用羊皮纸色调背景，营造舒适阅读体验
- **暖陶土色**：#c96442 作为品牌主色，温暖而不刺眼
- **圆润友好**：12-32px 圆角，柔和的视觉语言
- **Ring阴影**：使用 0px 0px 0px 1px 模式创造微妙层次

## 📁 样式文件结构

```
frontend/src/styles/
├── design-system.css    # 核心设计系统（CSS变量、主题）
├── utilities.css        # 工具类（快速构建界面）
├── h5.css              # 移动端H5样式
├── pc.css              # PC端样式优化
└── README.md           # 本文档
```

## 🎨 Claude 设计系统

### 配色方案

#### 主色调（暖陶土色系）
- **Terracotta Brand**: `#c96442` - 陶土品牌色（主要操作、强调）
- **Coral Accent**: `#d97757` - 珊瑚强调色（悬停状态）
- **Primary Dim**: `rgba(201, 100, 66, 0.15)` - 半透明背景

#### 背景色（浅色主题）
- **Parchment**: `#f5f4ed` - 羊皮纸主背景
- **Ivory**: `#faf9f5` - 象牙白卡片背景
- **Warm Sand**: `#e8e6dc` - 暖沙色按钮背景

#### 文本色（暖灰色系）
- **Near Black**: `#141413` - 主要文字
- **Olive Gray**: `#5e5d59` - 次要文字
- **Stone Gray**: `#87867f` - 第三级文字
- **Warm Silver**: `#b0aea5` - 深色主题文字

#### 状态色（暖色调）
- **Success**: `#6b8e23` - 橄榄绿
- **Warning**: `#d97757` - 暖橙
- **Danger**: `#b53333` - 深红
- **Info**: `#3898ec` - 蓝色（唯一冷色，用于焦点）

#### 边框色
- **Border Cream**: `#f0eee6` - 浅色边框
- **Border Warm**: `#e8e6dc` - 暖色边框
- **Border Dark**: `#30302e` - 深色边框

### 字体系统

```css
--font-serif: 'Georgia', 'Noto Serif SC', serif;    /* 标题 */
--font-sans: 'Noto Sans SC', -apple-system, sans-serif;  /* 正文/UI */
--font-mono: 'JetBrains Mono', 'SF Mono', monospace;     /* 代码 */
```

### 间距系统

```css
--space-xs: 4px
--space-sm: 8px
--space-md: 12px
--space-lg: 16px
--space-xl: 24px
--space-2xl: 32px
--space-3xl: 48px
```

### 圆角系统

```css
--radius-sharp: 4px       /* 尖锐：微小元素 */
--radius-subtle: 6px      /* 微圆：小按钮 */
--radius-comfortable: 8px /* 舒适：标准按钮 */
--radius-generous: 12px   /* 宽敞：主要按钮、输入框 */
--radius-very: 16px       /* 很圆：卡片 */
--radius-high: 24px       /* 高圆：标签 */
--radius-max: 32px        /* 最圆：英雄容器 */
```

### 过渡动画

```css
--transition-fast: 150ms cubic-bezier(0.4, 0, 0.2, 1)
--transition-base: 250ms cubic-bezier(0.4, 0, 0.2, 1)
--transition-slow: 350ms cubic-bezier(0.4, 0, 0.2, 1)
```

### Ring 阴影系统

Claude 风格的特色阴影模式：

```css
/* Ring 阴影 - 创造边框般的深度效果 */
--ring-warm: 0px 0px 0px 1px #d1cfc5
--ring-subtle: 0px 0px 0px 1px #e0dfd6
--ring-deep: 0px 0px 0px 1px #c2c0b6
```

### 主题支持

系统支持深色/浅色主题切换，通过 `data-theme` 属性控制：

```html
<div data-theme="light">...</div>
<div data-theme="dark">...</div>
```

#### 深色主题配色

- **背景**: `#141413` (Deep Dark)
- **卡片**: `#1a1a19`
- **边框**: `#30302e`
- **文字**: `#faf9f5` (Ivory)

## 🛠 工具类使用

### 间距类
```html
<div class="p-4">内边距 16px</div>
<div class="mb-3">下边距 12px</div>
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
- 悬停时边框变为陶土色
- 悬停抬升效果
- 统一圆角和阴影

### 表格样式
- 表头：大写字母、字间距、底部加粗边框
- 行悬停：暖灰色背景
- 分页：圆角按钮、陶土色激活色

### 按钮样式
所有按钮已应用 Claude 颜色：
- Primary: `#c96442` (Terracotta)
- Success: `#6b8e23` (Olive Green)
- Warning: `#d97757` (Warm Orange)
- Danger: `#b53333` (Deep Red)

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
    background: var(--color-bg-card);
    color: var(--color-text-primary);
  }
}
```

### 5. 字体层级
```css
/* 标题使用 Serif */
.page-title {
  font-family: var(--font-serif);
  font-weight: 500;  /* Serif 只用 500，不用 bold */
}

/* UI 使用 Sans */
.button {
  font-family: var(--font-sans);
}
```

## 🔄 更新日志

### 2026-01-18
- 应用 Claude 设计系统风格
- 更新主色调为暖陶土色 #c96442
- 更新背景色为羊皮纸色 #f5f4ed
- 更新文字色为暖灰色系
- 添加 Ring 阴影系统
- 扩展圆角系统到 32px

### 2026-01-17
- 新增 `utilities.css` 工具类文件
- 新增 `h5.css` 移动端样式
- 新增 `pc.css` PC端优化样式
- 优化 H5 TasksView 页面
- 优化 H5 IndexView 页面
- 优化 MobileHeader 组件
- 优化 PC端 TaskListView 页面

## 📚 参考资源

- [Claude 设计系统](https://getdesign.md/claude/design-md)
- [Element Plus 文档](https://element-plus.org/)
- [Vant 文档](https://vant-ui.github.io/)
- [CSS 变量 MDN](https://developer.mozilla.org/zh-CN/docs/Web/CSS/Using_CSS_custom_properties)
