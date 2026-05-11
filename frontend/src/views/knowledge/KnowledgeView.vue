<template>
  <div class="knowledge-view">
    <div class="header">
      <h2>知识库</h2>
      <div class="header-actions">
        <el-button type="primary" @click="showUploadDialog = true" v-if="canManage" plain>
          <el-icon><Upload /></el-icon>
          PDF 上传
        </el-button>
        <el-button type="primary" @click="showCreateDialog = true" v-if="canManage">
          <el-icon><Plus /></el-icon>
          新增知识
        </el-button>
      </div>
    </div>

    <!-- Search -->
    <el-card class="search-card">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索故障现象、解决方案..."
        clearable
        @clear="handleSearch"
        @keyup.enter="handleSearch"
      >
        <template #append>
          <el-button :icon="Search" @click="handleSearch" />
        </template>
      </el-input>
    </el-card>

    <!-- Articles List -->
    <el-card class="list-card">
      <el-table :data="articles" v-loading="loading" stripe>
        <el-table-column prop="title" label="标题" min-width="200">
          <template #default="{ row }">
            <div class="article-title">{{ row.title }}</div>
            <div class="article-tags" v-if="row.tags?.length">
              <el-tag
                v-for="tag in row.tags"
                :key="tag"
                size="small"
                type="info"
              >{{ tag }}</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="fault_phenomenon" label="故障现象" min-width="150">
          <template #default="{ row }">{{ row.fault_phenomenon || '-' }}</template>
        </el-table-column>
        <el-table-column prop="equipment_type_name" label="设备类型" width="120">
          <template #default="{ row }">{{ row.equipment_type_name || '-' }}</template>
        </el-table-column>
        <el-table-column prop="source_type" label="来源" width="100">
          <template #default="{ row }">
            <el-tag :type="getSourceTypeColor(row.source_type)" size="small">
              {{ getSourceTypeName(row.source_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="creator_name" label="创建人" width="100" />
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="viewArticle(row)">查看</el-button>
            <el-button link type="primary" size="small" @click="handleEdit(row)" v-if="canManage">编辑</el-button>
            <el-popconfirm title="确定删除吗？" @confirm="handleDelete(row.id)" v-if="canManage">
              <template #reference>
                <el-button link type="danger" size="small">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="loadArticles"
          @current-change="loadArticles"
        />
      </div>
    </el-card>

    <!-- Create/Edit Dialog -->
    <el-dialog
      v-model="showCreateDialog"
      :title="editingArticle ? '编辑知识' : '新增知识'"
      width="600px"
    >
      <el-form :model="articleForm" :rules="articleRules" ref="articleFormRef" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="articleForm.title" placeholder="请输入标题" />
        </el-form-item>
        <el-form-item label="设备类型">
          <el-select v-model="articleForm.equipment_type_id" placeholder="请选择设备类型" clearable style="width: 100%">
            <el-option
              v-for="type in equipmentTypes"
              :key="type.id"
              :label="type.name"
              :value="type.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="故障现象">
          <el-input v-model="articleForm.fault_phenomenon" type="textarea" :rows="2" placeholder="请输入故障现象" />
        </el-form-item>
        <el-form-item label="原因分析">
          <el-input v-model="articleForm.cause_analysis" type="textarea" :rows="3" placeholder="请输入原因分析" />
        </el-form-item>
        <el-form-item label="解决方案" prop="solution">
          <el-input v-model="articleForm.solution" type="textarea" :rows="5" placeholder="请输入解决方案" />
        </el-form-item>
        <el-form-item label="标签">
          <el-select
            v-model="articleForm.tags"
            multiple
            filterable
            allow-create
            placeholder="请选择或输入标签"
            style="width: 100%"
          >
            <el-option label="电气故障" value="电气故障" />
            <el-option label="机械故障" value="机械故障" />
            <el-option label="液压故障" value="液压故障" />
            <el-option label="软件故障" value="软件故障" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">确定</el-button>
      </template>
    </el-dialog>

    <!-- View Detail Dialog -->
    <el-dialog v-model="showDetailDialog" title="知识详情" width="700px">
      <div v-if="currentArticle" class="article-detail">
        <div class="detail-header">
          <h3>{{ currentArticle.title }}</h3>
          <div class="detail-meta">
            <el-tag v-if="currentArticle.equipment_type_name" type="info" size="small">
              {{ currentArticle.equipment_type_name }}
            </el-tag>
            <el-tag :type="getSourceTypeColor(currentArticle.source_type)" size="small">
              {{ getSourceTypeName(currentArticle.source_type) }}
            </el-tag>
            <span class="detail-creator">{{ currentArticle.creator_name }}</span>
            <span class="detail-time">{{ formatDateTime(currentArticle.created_at) }}</span>
          </div>
        </div>

        <div class="detail-tags" v-if="currentArticle.tags?.length">
          <el-tag
            v-for="tag in currentArticle.tags"
            :key="tag"
            size="small"
            type="info"
          >{{ tag }}</el-tag>
        </div>

        <el-descriptions :column="1" border class="detail-desc">
          <el-descriptions-item label="故障现象" v-if="currentArticle.fault_phenomenon">
            {{ currentArticle.fault_phenomenon }}
          </el-descriptions-item>
          <el-descriptions-item label="原因分析" v-if="currentArticle.cause_analysis">
            {{ currentArticle.cause_analysis }}
          </el-descriptions-item>
          <el-descriptions-item label="解决方案">
            <pre class="solution-text">{{ currentArticle.solution }}</pre>
          </el-descriptions-item>
        </el-descriptions>

        <div class="detail-actions" v-if="canManage && currentArticle.source_type === 'repair'">
          <el-button type="primary" size="small" @click="convertFromRepair(currentArticle.source_id!)">
            从维修单更新
          </el-button>
        </div>
      </div>
    </el-dialog>

    <!-- PDF Upload Dialog -->
    <el-dialog v-model="showUploadDialog" title="上传技术手册 (PDF)" width="450px">
      <el-form :model="uploadForm" label-width="100px">
        <el-form-item label="手册标题">
          <el-input v-model="uploadForm.title" placeholder="可选，默认使用文件名" />
        </el-form-item>
        <el-form-item label="设备类型">
          <el-select v-model="uploadForm.equipment_type_id" placeholder="请选择设备类型" clearable style="width: 100%">
            <el-option
              v-for="type in equipmentTypes"
              :key="type.id"
              :label="type.name"
              :value="type.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="PDF 文件">
          <el-upload
            class="upload-demo"
            drag
            action="#"
            :auto-upload="false"
            :limit="1"
            accept=".pdf"
            :on-change="handleFileChange"
            :on-remove="handleFileRemove"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              将文件拖到此处，或<em>点击上传</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                只能上传 PDF 文件
              </div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showUploadDialog = false">取消</el-button>
        <el-button type="primary" @click="submitUpload" :loading="uploading" :disabled="!selectedFile">
          开始上传并分片
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, type FormInstance, type FormRules, type UploadFile } from 'element-plus'
import { Plus, Search, Upload, UploadFilled } from '@element-plus/icons-vue'
import {
  getKnowledgeArticles,
  createKnowledgeArticle,
  updateKnowledgeArticle,
  deleteKnowledgeArticle,
  uploadManual,
  type KnowledgeArticle,
  type CreateKnowledgeArticleRequest
} from '@/api/knowledge'
import { equipmentTypeApi } from '@/api/equipment'
import type { EquipmentType } from '@/api/equipment'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const canManage = computed(() => authStore.hasRole('admin', 'engineer'))

const loading = ref(false)
const saving = ref(false)
const uploading = ref(false)
const articles = ref<KnowledgeArticle[]>([])
const equipmentTypes = ref<EquipmentType[]>([])

const searchKeyword = ref('')

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// Dialogs
const showCreateDialog = ref(false)
const showUploadDialog = ref(false)
const showDetailDialog = ref(false)
const editingArticle = ref<KnowledgeArticle | null>(null)
const currentArticle = ref<KnowledgeArticle | null>(null)

const articleFormRef = ref<FormInstance>()
const articleForm = reactive<CreateKnowledgeArticleRequest>({
  title: '',
  equipment_type_id: undefined,
  fault_phenomenon: '',
  cause_analysis: '',
  solution: '',
  source_type: 'manual',
  tags: []
})

const uploadForm = reactive({
  title: '',
  equipment_type_id: undefined as number | undefined
})
const selectedFile = ref<File | null>(null)

const articleRules: FormRules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  solution: [{ required: true, message: '请输入解决方案', trigger: 'blur' }]
}

// Methods
const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const getSourceTypeName = (type: string) => {
  const names: Record<string, string> = {
    repair: '维修单',
    manual: '手动录入',
    other: '其他'
  }
  return names[type] || '未知'
}

const getSourceTypeColor = (type: string) => {
  const colors: Record<string, string> = {
    repair: 'success',
    manual: 'info',
    other: 'warning'
  }
  return colors[type] || 'info'
}

const loadArticles = async () => {
  loading.value = true
  try {
    const res = await getKnowledgeArticles({
      keyword: searchKeyword.value || undefined,
      page: pagination.page,
      page_size: pagination.pageSize
    })
    articles.value = res.data.items
    pagination.total = res.data.total
  } catch (err) {
    ElMessage.error('加载知识库失败')
  } finally {
    loading.value = false
  }
}

const loadEquipmentTypes = async () => {
  try {
    const res = await equipmentTypeApi.getTypes()
    equipmentTypes.value = res.data
  } catch (err) {
    console.error('Failed to load equipment types:', err)
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadArticles()
}

const viewArticle = (article: KnowledgeArticle) => {
  currentArticle.value = article
  showDetailDialog.value = true
}

const handleEdit = (article: KnowledgeArticle) => {
  editingArticle.value = article
  Object.assign(articleForm, {
    title: article.title,
    equipment_type_id: article.equipment_type_id,
    fault_phenomenon: article.fault_phenomenon || '',
    cause_analysis: article.cause_analysis || '',
    solution: article.solution,
    tags: article.tags || []
  })
  showCreateDialog.value = true
}

const handleSave = async () => {
  if (!articleFormRef.value) return
  await articleFormRef.value.validate(async (valid) => {
    if (!valid) return
    saving.value = true
    try {
      if (editingArticle.value) {
        await updateKnowledgeArticle(editingArticle.value.id, articleForm)
        ElMessage.success('更新成功')
      } else {
        await createKnowledgeArticle(articleForm)
        ElMessage.success('创建成功')
      }
      showCreateDialog.value = false
      loadArticles()
    } catch (err: any) {
      ElMessage.error(err.response?.data?.error || '操作失败')
    } finally {
      saving.value = false
    }
  })
}

const handleDelete = async (id: number) => {
  try {
    await deleteKnowledgeArticle(id)
    ElMessage.success('删除成功')
    loadArticles()
  } catch (err: any) {
    ElMessage.error(err.response?.data?.error || '删除失败')
  }
}

const handleFileChange = (file: UploadFile) => {
  selectedFile.value = file.raw || null
}

const handleFileRemove = () => {
  selectedFile.value = null
}

const submitUpload = async () => {
  if (!selectedFile.value) return
  uploading.value = true
  try {
    const formData = new FormData()
    formData.append('file', selectedFile.value)
    if (uploadForm.title) {
      formData.append('title', uploadForm.title)
    }
    if (uploadForm.equipment_type_id) {
      formData.append('equipment_type_id', uploadForm.equipment_type_id.toString())
    }

    const res = await uploadManual(formData)
    ElMessage.success(res.data.message || '上传成功，正在后台解析分片')
    showUploadDialog.value = false
    uploadForm.title = ''
    uploadForm.equipment_type_id = undefined
    selectedFile.value = null
  } catch (err: any) {
    ElMessage.error(err.response?.data?.error || '上传失败')
  } finally {
    uploading.value = false
  }
}

const convertFromRepair = (orderId: number) => {
  // TODO: Implement convert from repair
  ElMessage.info('功能开发中')
}

onMounted(() => {
  loadArticles()
  loadEquipmentTypes()
})
</script>

<style scoped>
.knowledge-view {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header h2 {
  margin: 0;
  font-size: 20px;
  color: #303133;
}

.search-card {
  margin-bottom: 20px;
}

.list-card {
  margin-bottom: 20px;
}

.article-title {
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}

.article-tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.article-detail {
  padding: 10px 0;
}

.detail-header h3 {
  margin: 0 0 12px 0;
  font-size: 18px;
  color: #303133;
}

.detail-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  font-size: 12px;
  color: #909399;
}

.detail-creator,
.detail-time {
  color: #909399;
  font-size: 12px;
}

.detail-tags {
  margin-bottom: 16px;
}

.detail-desc {
  margin-bottom: 16px;
}

.solution-text {
  white-space: pre-wrap;
  font-family: inherit;
  line-height: 1.6;
}

.detail-actions {
  display: flex;
  gap: 8px;
}
</style>
