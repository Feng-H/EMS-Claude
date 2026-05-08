<template>
  <div class="agent-integration-view">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span class="title">Agent 外部集成</span>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <!-- API Keys Tab -->
        <el-tab-pane label="API 密钥" name="apikeys">
          <div class="tab-content">
            <div class="actions">
              <el-button type="primary" @click="showCreateKeyDialog = true">
                <el-icon><Plus /></el-icon> 创建新密钥
              </el-button>
            </div>
            
            <el-table :data="apiKeys" border stripe v-loading="loadingKeys">
              <el-table-column prop="name" label="名称" width="150" />
              <el-table-column prop="description" label="描述" min-width="150" show-overflow-tooltip />
              <el-table-column prop="scopes" label="权限范围" width="180">
                <template #default="{ row }">
                  <el-tag v-for="s in (row.scopes || '').split(',')" :key="s" size="small" style="margin-right: 4px" v-show="s">
                    {{ s }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="rate_limit" label="频控 (req/m)" width="120">
                <template #default="{ row }">
                  {{ row.rate_limit > 0 ? row.rate_limit : '无限制' }}
                </template>
              </el-table-column>
              <el-table-column label="状态" width="80">
                <template #default="{ row }">
                  <el-tag :type="row.is_active ? 'success' : 'danger'" size="small">
                    {{ row.is_active ? '正常' : '禁用' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="last_used_at" label="最后使用" width="160">
                <template #default="{ row }">
                  {{ formatDate(row.last_used_at) }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="100" fixed="right">
                <template #default="{ row }">
                  <el-popconfirm title="确定要删除此密钥吗？" @confirm="handleDeleteKey(row.id)">
                    <template #reference>
                      <el-button type="danger" link>删除</el-button>
                    </template>
                  </el-popconfirm>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>

        <!-- Tools Tab -->
        <el-tab-pane label="工具集发现" name="tools">
          <div class="tab-content">
            <el-alert
              title="系统向 Agent 暴露的工具列表。Agent 可以调用这些工具来执行操作或获取数据。"
              type="info"
              show-icon
              style="margin-bottom: 20px"
            />
            <el-table :data="tools" border v-loading="loadingTools" row-key="name">
              <el-table-column type="expand">
                <template #default="{ row }">
                  <div class="schema-preview">
                    <h4>参数 Schema (JSON Schema):</h4>
                    <pre><code>{{ formatJson(row.input_schema) }}</code></pre>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="name" label="工具名称" width="200">
                <template #default="{ row }">
                  <el-tag>{{ row.name }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="description" label="描述" />
            </el-table>
          </div>
        </el-tab-pane>

        <!-- Proactive Push Tab -->
        <el-tab-pane label="主动推送配置" name="push">
          <div class="tab-content">
            <el-alert
              title="配置系统主动向外推送消息（例如发送到飞书、邮件或 Webhook）。每个推送类型可单独配置范围。"
              type="warning"
              show-icon
              style="margin-bottom: 20px"
            />
            
            <el-row :gutter="20">
              <el-col :span="8" v-for="(push, index) in pushConfigs" :key="index">
                <el-card class="push-card" shadow="hover">
                  <template #header>
                    <div class="push-header">
                      <span>{{ push.label }}</span>
                      <el-switch v-model="push.enabled" />
                    </div>
                  </template>
                  <div class="push-body">
                    <p class="push-desc">{{ push.description }}</p>
                    <div class="scope-config">
                      <span class="scope-label">Webhook URL</span>
                      <el-input
                        v-model="push.webhook_url"
                        placeholder="https://example.com/webhook"
                        style="margin-bottom: 10px"
                      />
                      <span class="scope-label">签名密钥 (Secret)</span>
                      <el-input
                        v-model="push.secret"
                        type="password"
                        show-password
                        placeholder="用于签名 payload"
                        style="margin-bottom: 10px"
                      />
                      <span class="scope-label">推送范围 (JSON)</span>
                      <el-input
                        v-model="push.scope"
                        type="textarea"
                        :rows="2"
                        placeholder='例如: {"workshop_id": 1}'
                      />
                    </div>
                    <el-button type="primary" class="save-btn" @click="handleSavePush(push)" :loading="push.saving">
                      保存配置
                    </el-button>
                  </div>
                </el-card>
              </el-col>
            </el-row>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Create API Key Dialog -->
    <el-dialog v-model="showCreateKeyDialog" title="创建 API 密钥" width="550px">
      <el-form ref="keyFormRef" :model="keyForm" :rules="keyRules" label-width="90px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="keyForm.name" placeholder="请输入应用名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="keyForm.description" type="textarea" placeholder="描述此密钥的用途" />
        </el-form-item>
        <el-form-item label="权限范围" prop="scopes">
          <el-select v-model="keyForm.scopes" multiple placeholder="请选择允许的操作范围" style="width: 100%">
            <el-option label="读取设备档案 (read:equipment)" value="read:equipment" />
            <el-option label="读取预测分析 (read:prediction)" value="read:prediction" />
            <el-option label="读取备件库存 (read:sparepart)" value="read:sparepart" />
            <el-option label="创建维修工单 (write:repair)" value="write:repair" />
          </el-select>
        </el-form-item>
        <el-form-item label="频率限制" prop="rate_limit">
          <el-input-number v-model="keyForm.rate_limit" :min="0" :max="1000" style="width: 100%" />
          <div class="form-tip">每分钟允许的请求数 (0 为不限制)</div>
        </el-form-item>
        <el-form-item label="过期时间" prop="expires_in">
          <el-select v-model="keyForm.expires_in" style="width: 100%">
            <el-option label="30 天" :value="30" />
            <el-option label="90 天" :value="90" />
            <el-option label="1 年" :value="365" />
            <el-option label="永不过期" :value="0" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCreateKeyDialog = false">取消</el-button>
          <el-button type="primary" @click="handleCreateKey" :loading="creatingKey">
            确认创建
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- Secret Display Dialog -->
    <el-dialog
      v-model="showSecretDialog"
      title="密钥已生成"
      width="500px"
      :close-on-click-modal="false"
      :show-close="false"
    >
      <el-alert
        title="请务必复制并妥善保存您的 API 密钥。此密钥仅显示一次，如果丢失，您将需要重新生成。"
        type="error"
        show-icon
        :closable="false"
        style="margin-bottom: 20px"
      />
      <div class="secret-box">
        <el-input v-model="newSecret" readonly>
          <template #append>
            <el-button @click="copySecret">
              <el-icon><DocumentCopy /></el-icon> 复制
            </el-button>
          </template>
        </el-input>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="showSecretDialog = false">我已保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Plus, DocumentCopy } from '@element-plus/icons-vue'
import { agentApi } from '@/api/agent'

// Tab state
const activeTab = ref('apikeys')

// Data formatters
const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString()
}

const formatJson = (obj: any) => {
  try {
    return JSON.stringify(obj, null, 2)
  } catch (e) {
    return String(obj)
  }
}

// --- API Keys Tab ---
const apiKeys = ref<any[]>([])
const loadingKeys = ref(false)

const loadApiKeys = async () => {
  loadingKeys.value = true
  try {
    const res = await agentApi.listAPIKeys()
    apiKeys.value = (res as any).data || res || [] 
  } catch (error: any) {
    ElMessage.error(error.message || '获取 API 密钥失败')
  } finally {
    loadingKeys.value = false
  }
}

const showCreateKeyDialog = ref(false)
const creatingKey = ref(false)
const keyFormRef = ref<FormInstance>()
const keyForm = reactive({
  name: '',
  description: '',
  scopes: ['read:equipment'],
  rate_limit: 0,
  expires_in: 0
})

const keyRules: FormRules = {
  name: [{ required: true, message: '请输入应用名称', trigger: 'blur' }]
}

const showSecretDialog = ref(false)
const newSecret = ref('')

const handleCreateKey = async () => {
  if (!keyFormRef.value) return
  await keyFormRef.value.validate(async (valid) => {
    if (valid) {
      creatingKey.value = true
      try {
        const res = await agentApi.createAPIKey(keyForm)
        const secretKey = (res as any).data?.key || (res as any).key
        if (secretKey) {
          newSecret.value = secretKey
          showSecretDialog.value = true
        }
        showCreateKeyDialog.value = false
        keyFormRef.value?.resetFields()
        loadApiKeys()
      } catch (error: any) {
        ElMessage.error(error.message || '创建失败')
      } finally {
        creatingKey.value = false
      }
    }
  })
}

const copySecret = () => {
  if (newSecret.value) {
    navigator.clipboard.writeText(newSecret.value).then(() => {
      ElMessage.success('已复制到剪贴板')
    }).catch(() => {
      ElMessage.error('复制失败，请手动选择复制')
    })
  }
}

const handleDeleteKey = async (id: number) => {
  try {
    await agentApi.deleteAPIKey(id)
    ElMessage.success('删除成功')
    loadApiKeys()
  } catch (error: any) {
    ElMessage.error(error.message || '删除失败')
  }
}

// --- Tools Tab ---
const tools = ref<any[]>([])
const loadingTools = ref(false)

const loadTools = async () => {
  loadingTools.value = true
  try {
    const res = await agentApi.listTools()
    tools.value = (res as any).data?.tools || (res as any).tools || []
  } catch (error: any) {
    ElMessage.error(error.message || '获取工具列表失败')
  } finally {
    loadingTools.value = false
  }
}

// --- Proactive Push Tab ---
const pushConfigs = ref([
  {
    type: 'predictive_maintenance',
    label: '预测性维护预警',
    description: '当分析出设备剩余寿命较低或停机风险较高时推送。',
    enabled: false,
    scope: '{}',
    webhook_url: '',
    secret: '',
    saving: false
  },
  {
    type: 'ng_inspection',
    label: '点检异常推送',
    description: '当点检发现异常状态时，主动推送通知。',
    enabled: false,
    scope: '{}',
    webhook_url: '',
    secret: '',
    saving: false
  },
  {
    type: 'repair_request',
    label: '新报修申请推送',
    description: '当有新的报修工单创建时，主动推送通知。',
    enabled: false,
    scope: '{}',
    webhook_url: '',
    secret: '',
    saving: false
  },
  {
    type: 'low_stock',
    label: '备件低库存推送',
    description: '当备件库存低于设定的安全阈值时，主动推送通知。',
    enabled: false,
    scope: '{}',
    webhook_url: '',
    secret: '',
    saving: false
  }
])

const loadSubscriptions = async () => {
  try {
    const res = await agentApi.listSubscriptions()
    const subs = (res as any).data || res || []
    subs.forEach((s: any) => {
      const config = pushConfigs.value.find(c => c.type === s.push_type)
      if (config) {
        config.enabled = s.enabled
        config.scope = s.scope || '{}'
        config.webhook_url = s.webhook_url || ''
        config.secret = s.secret || ''
      }
    })
  } catch (error) {
    console.error('Failed to load subscriptions')
  }
}

const handleSavePush = async (push: any) => {
  push.saving = true
  try {
    let parsedScope = {}
    if (push.scope && push.scope.trim()) {
      try {
        parsedScope = JSON.parse(push.scope)
      } catch (e) {
        ElMessage.warning('Scope 必须是有效的 JSON 格式')
        push.saving = false
        return
      }
    }

    await agentApi.subscribe({
      push_type: push.type,
      enabled: push.enabled,
      scope: parsedScope,
      webhook_url: push.webhook_url,
      secret: push.secret
    })
    ElMessage.success(`${push.label} 配置保存成功`)
  } catch (error: any) {
    ElMessage.error(error.message || '配置保存失败')
  } finally {
    push.saving = false
  }
}

onMounted(() => {
  loadApiKeys()
  loadTools()
  loadSubscriptions()
})
</script>

<style scoped>
.agent-integration-view {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header .title {
  font-size: 18px;
  font-weight: bold;
}

.tab-content {
  padding: 10px 0;
}

.actions {
  margin-bottom: 20px;
}

.schema-preview {
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 4px;
}

.schema-preview h4 {
  margin-top: 0;
  margin-bottom: 10px;
  color: #606266;
}

.schema-preview pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
  color: #333;
  font-family: 'Courier New', Courier, monospace;
}

.push-card {
  margin-bottom: 20px;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.push-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: bold;
}

.push-body {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
}

.push-desc {
  color: #606266;
  font-size: 14px;
  margin-bottom: 15px;
  min-height: 40px;
}

.scope-config {
  margin-bottom: 20px;
  flex-grow: 1;
}

.scope-label {
  display: block;
  font-size: 14px;
  margin-bottom: 8px;
  color: #606266;
}

.save-btn {
  width: 100%;
}

.secret-box {
  background-color: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
  border: 1px dashed #dcdfe6;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
  margin-top: 4px;
}
</style>
