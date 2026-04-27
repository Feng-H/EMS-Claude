<template>
  <div class="h5-knowledge-view">
    <mobile-header title="知识库" />

    <!-- 搜索栏 -->
    <van-search
      v-model="searchKeyword"
      placeholder="搜索故障现象、解决方案..."
      @search="onSearch"
      @clear="onReset"
    />

    <!-- 列表 -->
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <div class="list-container">
        <template v-if="loading && !refreshing">
          <van-skeleton v-for="i in 4" :key="i" title :row="3" style="margin-bottom: 20px" />
        </template>
        <van-list
          v-else
          v-model:loading="loading"
          :finished="finished"
          finished-text="没有更多知识了"
          @load="onLoad"
        >
          <van-empty v-if="articles.length === 0" description="暂无相关知识" />
          
          <div
            v-for="item in articles"
            :key="item.id"
            class="article-card"
            @click="showDetail(item)"
          >
            <div class="card-header">
              <span class="article-title">{{ item.title }}</span>
              <van-tag :type="getSourceTypeColor(item.source_type)" plain>
                {{ getSourceTypeName(item.source_type) }}
              </van-tag>
            </div>
            <div class="card-body">
              <p class="phenomenon" v-if="item.fault_phenomenon">
                <span class="label">现象:</span> {{ item.fault_phenomenon }}
              </p>
              <div class="card-tags" v-if="item.tags?.length">
                <van-tag
                  v-for="tag in item.tags"
                  :key="tag"
                  round
                  type="primary"
                  color="#f2f2f2"
                  text-color="#666"
                  class="tag-item"
                >{{ tag }}</van-tag>
              </div>
            </div>
            <div class="card-footer">
              <span class="equipment-type" v-if="item.equipment_type_name">
                <van-icon name="cluster-o" /> {{ item.equipment_type_name }}
              </span>
              <span class="create-time">{{ formatTime(item.created_at) }}</span>
            </div>
          </div>
        </van-list>
      </div>
    </van-pull-refresh>

    <!-- 详情弹出层 -->
    <van-popup
      v-model:show="detailShow"
      position="right"
      :style="{ width: '100%', height: '100%' }"
    >
      <div class="article-detail-page" v-if="currentArticle">
        <van-nav-bar
          title="知识详情"
          left-arrow
          fixed
          @click-left="detailShow = false"
        />
        <div class="detail-content">
          <h2 class="detail-title">{{ currentArticle.title }}</h2>
          <div class="detail-meta">
            <van-tag :type="getSourceTypeColor(currentArticle.source_type)">
              {{ getSourceTypeName(currentArticle.source_type) }}
            </van-tag>
            <span class="meta-item">{{ currentArticle.creator_name }}</span>
            <span class="meta-item">{{ formatDateTime(currentArticle.created_at) }}</span>
          </div>

          <div class="detail-section">
            <h4 class="section-title">故障现象</h4>
            <p class="section-text">{{ currentArticle.fault_phenomenon || '无' }}</p>
          </div>

          <div class="detail-section">
            <h4 class="section-title">原因分析</h4>
            <p class="section-text">{{ currentArticle.cause_analysis || '无' }}</p>
          </div>

          <div class="detail-section">
            <h4 class="section-title">解决方案</h4>
            <div class="solution-box">
              <pre class="solution-text">{{ currentArticle.solution }}</pre>
            </div>
          </div>

          <div class="detail-tags" v-if="currentArticle.tags?.length">
            <van-tag
              v-for="tag in currentArticle.tags"
              :key="tag"
              round
              type="primary"
              size="medium"
              class="tag-item"
            >{{ tag }}</van-tag>
          </div>
        </div>
      </div>
    </van-popup>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { showToast } from 'vant'
import { getKnowledgeArticles, type KnowledgeArticle } from '@/api/knowledge'
import MobileHeader from '@/components/mobile/MobileHeader.vue'

const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const articles = ref<KnowledgeArticle[]>([])
const searchKeyword = ref('')

const pagination = reactive({
  page: 1,
  pageSize: 10
})

const detailShow = ref(false)
const currentArticle = ref<KnowledgeArticle | null>(null)

const getSourceTypeName = (type: string) => {
  const names: Record<string, string> = {
    repair: '维修单',
    manual: '手动录入',
    other: '其他'
  }
  return names[type] || '未知'
}

const getSourceTypeColor = (type: string) => {
  const colors: Record<string, any> = {
    repair: 'success',
    manual: 'primary',
    other: 'warning'
  }
  return colors[type] || 'default'
}

const formatTime = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`
}

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const loadData = async () => {
  try {
    const res = await getKnowledgeArticles({
      keyword: searchKeyword.value || undefined,
      page: pagination.page,
      page_size: pagination.pageSize
    })
    
    if (refreshing.value) {
      articles.value = res.data.items
    } else {
      articles.value.push(...res.data.items)
    }

    if (articles.value.length >= res.data.total) {
      finished.value = true
    } else {
      pagination.page++
    }
  } catch (err) {
    showToast('加载失败')
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

const onRefresh = () => {
  finished.value = false
  pagination.page = 1
  loadData()
}

const onLoad = () => {
  loadData()
}

const onSearch = () => {
  articles.value = []
  onRefresh()
}

const onReset = () => {
  searchKeyword.value = ''
  onSearch()
}

const showDetail = (item: KnowledgeArticle) => {
  currentArticle.value = item
  detailShow.value = true
}

onMounted(() => {
  // initial load handled by van-list onLoad
})
</script>

<style scoped>
.h5-knowledge-view {
  min-height: 100vh;
  background: var(--color-bg-primary);
  padding-top: 46px;
  padding-bottom: 60px;
}

.list-container {
  padding: 12px;
}

.article-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-very);
  padding: 16px;
  margin-bottom: 12px;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--color-border);
  transition: all 0.2s;
}

.article-card:active {
  background: var(--color-bg-tertiary);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 10px;
  gap: 8px;
}

.article-title {
  font-weight: 600;
  font-size: 16px;
  color: var(--color-text-primary);
  line-height: 1.4;
  flex: 1;
}

.card-body {
  margin-bottom: 12px;
}

.phenomenon {
  font-size: 13px;
  color: var(--color-text-secondary);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.5;
}

.label {
  color: var(--color-text-tertiary);
  font-weight: 500;
}

.card-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 8px;
}

.tag-item {
  font-size: 10px;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: var(--color-text-tertiary);
  border-top: 1px solid var(--color-border-dim);
  padding-top: 10px;
}

.equipment-type {
  display: flex;
  align-items: center;
  gap: 4px;
}

/* 详情页面样式 */
.article-detail-page {
  height: 100%;
  background: var(--color-bg-primary);
  display: flex;
  flex-direction: column;
}

.detail-content {
  padding: 66px 20px 40px;
  overflow-y: auto;
}

.detail-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--color-text-primary);
  margin-bottom: 12px;
  line-height: 1.3;
}

.detail-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
  font-size: 12px;
  color: var(--color-text-tertiary);
  flex-wrap: wrap;
}

.detail-section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 8px;
  display: flex;
  align-items: center;
}

.section-title::before {
  content: '';
  width: 4px;
  height: 14px;
  background: var(--color-terracotta);
  margin-right: 8px;
  border-radius: 2px;
}

.section-text {
  font-size: 14px;
  color: var(--color-text-secondary);
  line-height: 1.6;
}

.solution-box {
  background: var(--color-bg-tertiary);
  padding: 12px;
  border-radius: 8px;
  border: 1px solid var(--color-border);
}

.solution-text {
  font-size: 14px;
  color: var(--color-text-primary);
  white-space: pre-wrap;
  font-family: inherit;
  margin: 0;
}

.detail-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 10px;
}

/* 暗色模式适配 */
@media (prefers-color-scheme: dark) {
  .article-card {
    background: var(--color-bg-card);
  }
  .solution-box {
    background: #2a2a22;
  }
}
</style>
