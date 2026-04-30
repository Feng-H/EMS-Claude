<template>
  <div class="lark-integration">
    <div class="page-header">
      <div class="header-title">
        <h2>飞书机器人集成</h2>
        <p>配置您的专属飞书机器人，实现设备故障提醒与智能分析</p>
      </div>
    </div>

    <div class="content-layout">
      <!-- 左侧：配置表单 -->
      <div class="config-section">
        <el-card class="section-card">
          <template #header>
            <div class="card-header">
              <span>机器人参数配置</span>
            </div>
          </template>

          <el-form :model="configForm" label-position="top" @submit.prevent="handleSave">
            <el-form-item label="App ID" required>
              <el-input v-model="configForm.app_id" placeholder="例如：cli_a123456789" />
            </el-form-item>
            
            <el-form-item label="App Secret" required>
              <el-input 
                v-model="configForm.app_secret" 
                type="password" 
                show-password 
                :placeholder="hasAppSecret ? '******** (已保存)' : '请输入 App Secret'" 
              />
            </el-form-item>

            <el-divider />

            <el-form-item label="Verification Token (校验令牌)">
              <el-input v-model="configForm.verification_token" placeholder="飞书后台 -> 事件订阅 -> Verification Token" />
            </el-form-item>

            <el-form-item label="Encrypt Key (加密策略)">
              <el-input 
                v-model="configForm.encrypt_key" 
                type="password" 
                show-password 
                :placeholder="hasEncryptKey ? '******** (已保存)' : '飞书后台 -> 事件订阅 -> Encrypt Key'" 
              />
            </el-form-item>

            <el-form-item label="Webhook URL (数据接收地址)">
              <div class="webhook-url-container">
                <el-input v-model="webhookUrl" readonly>
                  <template #append>
                    <el-button @click="copyWebhookUrl">复制</el-button>
                  </template>
                </el-input>
                <p class="input-tip">将此地址粘贴到飞书开发者后台的“事件订阅” -> “请求地址”中</p>
              </div>
            </el-form-item>

            <div class="form-actions">
              <el-button type="primary" :loading="saving" @click="handleSave">保存配置</el-button>
              <el-button @click="fetchConfig">重置</el-button>
            </div>
          </el-form>
        </el-card>
      </div>

      <!-- 右侧：操作指引 -->
      <div class="guide-section">
        <el-card class="section-card">
          <template #header>
            <div class="card-header">
              <span>配置指引</span>
            </div>
          </template>
          
          <div class="guide-content">
            <div class="guide-step">
              <h4>1. 创建飞书应用</h4>
              <p>登录 <a href="https://open.feishu.cn/app" target="_blank">飞书开放平台</a>，点击“创建企业自建应用”。</p>
            </div>

            <div class="guide-step">
              <h4>2. 获取凭证</h4>
              <p>在应用详情页的“凭证与基础信息”中，找到 <strong>App ID</strong> 和 <strong>App Secret</strong>。</p>
            </div>

            <div class="guide-step">
              <h4>3. 配置事件订阅</h4>
              <p>在“事件订阅”栏目中：</p>
              <ul>
                <li>粘贴左侧显示的 <strong>Webhook URL</strong>。</li>
                <li>复制 <strong>Verification Token</strong> 和 <strong>Encrypt Key</strong> 到本系统。</li>
                <li>添加事件：<code>接收消息 v2.0 (im.message.receive_v1)</code>。</li>
              </ul>
            </div>

            <div class="guide-step">
              <h4>4. 开启机器人与权限</h4>
              <p>在“应用功能”中开启“机器人”功能。在“权限管理”中勾选：</p>
              <ul>
                <li>获取用户 ID (<code>user:id</code>)</li>
                <li>获取用户发送的消息 (<code>im:message</code>)</li>
                <li>给用户发送消息 (<code>im:message.p2p:readonly</code>)</li>
              </ul>
            </div>

            <div class="guide-step">
              <h4>5. 发布版本</h4>
              <p>在“版本管理与发布”中，创建一个新版本并申请上线（管理员审核通过后生效）。</p>
            </div>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { authApi, type LarkConfigReq } from '@/api/auth'
import { ElMessage } from 'element-plus'

const configForm = ref<LarkConfigReq>({
  app_id: '',
  app_secret: '',
  verification_token: '',
  encrypt_key: ''
})

const webhookUrl = ref('')
const hasAppSecret = ref(false)
const hasEncryptKey = ref(false)
const saving = ref(false)

const fetchConfig = async () => {
  try {
    const res = await authApi.getLarkConfig()
    const data = res.data
    configForm.value.app_id = data.app_id
    configForm.value.verification_token = data.verification_token
    webhookUrl.value = data.webhook_url
    hasAppSecret.value = data.has_app_secret
    hasEncryptKey.value = data.has_encrypt_key
    // Reset secrets in form if they exist
    configForm.value.app_secret = ''
    configForm.value.encrypt_key = ''
  } catch (err: any) {
    ElMessage.error('获取配置失败: ' + (err.message || '未知错误'))
  }
}

const handleSave = async () => {
  if (!configForm.value.app_id) {
    ElMessage.warning('请输入 App ID')
    return
  }
  
  saving.value = true
  try {
    await authApi.updateLarkConfig(configForm.value)
    ElMessage.success('保存成功')
    await fetchConfig()
  } catch (err: any) {
    ElMessage.error('保存失败: ' + (err.message || '未知错误'))
  } finally {
    saving.value = false
  }
}

const copyWebhookUrl = () => {
  navigator.clipboard.writeText(webhookUrl.value)
  ElMessage.success('已复制到剪贴板')
}

onMounted(() => {
  fetchConfig()
})
</script>

<style scoped>
.lark-integration {
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
}

.header-title h2 {
  margin: 0 0 8px 0;
  font-size: 24px;
  color: #1a1a1a;
}

.header-title p {
  margin: 0;
  color: #666;
  font-size: 14px;
}

.content-layout {
  display: grid;
  grid-template-columns: 1fr 350px;
  gap: 24px;
}

.section-card {
  height: 100%;
}

.card-header {
  font-weight: bold;
}

.webhook-url-container {
  background-color: #f8f9fa;
  padding: 8px;
  border-radius: 4px;
}

.input-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 8px;
}

.form-actions {
  margin-top: 32px;
  display: flex;
  gap: 12px;
}

.guide-content {
  font-size: 14px;
  line-height: 1.6;
}

.guide-step {
  margin-bottom: 24px;
}

.guide-step h4 {
  margin: 0 0 8px 0;
  color: #303133;
}

.guide-step p {
  margin: 0 0 8px 0;
  color: #606266;
}

.guide-step ul {
  padding-left: 20px;
  color: #606266;
  margin: 8px 0;
}

.guide-step code {
  background: #f0f2f5;
  padding: 2px 4px;
  border-radius: 3px;
  font-family: monospace;
}

@media (max-width: 1024px) {
  .content-layout {
    grid-template-columns: 1fr;
  }
}
</style>
