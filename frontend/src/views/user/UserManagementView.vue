<template>
  <div class="user-management">
    <div class="page-header">
      <div class="header-title">
        <h2>人员管理</h2>
        <p>管理系统用户账号及权限</p>
      </div>
      <el-button type="primary" @click="openCreateDialog">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 16px; height: 16px; margin-right: 8px;">
          <line x1="12" y1="5" x2="12" y2="19"/>
          <line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
        新增用户
      </el-button>
    </div>

    <!-- 待审核申请 -->
    <div v-if="pendingUsers.length > 0" class="section-card pending-section">
      <div class="section-header">
        <h3>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 20px; height: 20px; margin-right: 8px;">
            <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/>
            <path d="M13.73 21a2 2 0 0 1-3.46 0"/>
          </svg>
          待审核申请
          <el-badge :value="pendingUsers.length" style="margin-left: 8px;" />
        </h3>
      </div>
      <el-table :data="pendingUsers" style="width: 100%">
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="name" label="姓名" width="120" />
        <el-table-column prop="role" label="角色" width="120">
          <template #default="{ row }">
            <el-tag :type="getRoleTagType(row.role)">{{ getRoleText(row.role) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="phone" label="电话" width="140" />
        <el-table-column prop="created_at" label="申请时间" width="180" />
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button type="success" size="small" @click="handleApprove(row, true)">通过</el-button>
            <el-button type="danger" size="small" @click="handleReject(row)">拒绝</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 用户列表 -->
    <div class="section-card">
      <div class="section-header">
        <h3>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 20px; height: 20px; margin-right: 8px;">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
          </svg>
          用户列表
        </h3>
        <el-input
          v-model="searchKeyword"
          placeholder="搜索用户名或姓名"
          style="width: 250px;"
          clearable
        >
          <template #prefix>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 16px; height: 16px;">
              <circle cx="11" cy="11" r="8"/>
              <path d="m21 21-4.35-4.35"/>
            </svg>
          </template>
        </el-input>
      </div>

      <el-table :data="filteredUsers" style="width: 100%">
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="name" label="姓名" width="120" />
        <el-table-column prop="role" label="角色" width="120">
          <template #default="{ row }">
            <el-tag :type="getRoleTagType(row.role)">{{ getRoleText(row.role) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="phone" label="电话" width="140" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_active ? 'success' : 'info'">
              {{ row.is_active ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="审核状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getApprovalStatusType(row.approval_status)">
              {{ getApprovalStatusText(row.approval_status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="openEditDialog(row)">编辑</el-button>
            <el-button link type="danger" size="small" @click="handleToggleActive(row)">
              {{ row.is_active ? '禁用' : '启用' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 创建/编辑用户对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑用户' : '新增用户'"
      width="500px"
    >
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="formData.username" :disabled="isEdit" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!isEdit">
          <el-input v-model="formData.password" type="password" placeholder="默认密码：password123" show-password />
        </el-form-item>
        <el-form-item label="姓名" prop="name">
          <el-input v-model="formData.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="formData.role" placeholder="请选择角色" style="width: 100%;">
            <el-option label="系统管理员" value="admin" />
            <el-option label="设备主管" value="supervisor" />
            <el-option label="设备工程师" value="engineer" />
            <el-option label="维修工" value="maintenance" />
            <el-option label="操作工" value="operator" />
          </el-select>
        </el-form-item>
        <el-form-item label="电话" prop="phone">
          <el-input v-model="formData.phone" placeholder="请输入电话" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          {{ isEdit ? '保存' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { userApi, type User, type CreateUserRequest, type UpdateUserRequest } from '@/api/user'

const users = ref<User[]>([])
const pendingUsers = ref<User[]>([])
const searchKeyword = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()

const currentUser = ref<User | null>(null)

const formData = reactive<CreateUserRequest & { phone?: string }>({
  username: '',
  password: '',
  name: '',
  role: '',
  phone: '',
})

const formRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '用户名长度在 3 到 50 个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 个字符', trigger: 'blur' },
  ],
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' },
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' },
  ],
}

const filteredUsers = computed(() => {
  if (!searchKeyword.value) {
    return users.value.filter(u => u.approval_status !== 'pending')
  }
  const keyword = searchKeyword.value.toLowerCase()
  return users.value.filter(u =>
    u.approval_status !== 'pending' &&
    (u.username.toLowerCase().includes(keyword) || u.name.toLowerCase().includes(keyword))
  )
})

function getRoleText(role: string): string {
  const roleMap: Record<string, string> = {
    admin: '系统管理员',
    supervisor: '设备主管',
    engineer: '设备工程师',
    maintenance: '维修工',
    operator: '操作工',
  }
  return roleMap[role] || role
}

function getRoleTagType(role: string): string {
  const typeMap: Record<string, string> = {
    admin: 'danger',
    supervisor: 'warning',
    engineer: 'primary',
    maintenance: 'info',
    operator: 'success',
  }
  return typeMap[role] || 'info'
}

function getApprovalStatusText(status: string): string {
  const statusMap: Record<string, string> = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝',
  }
  return statusMap[status] || status
}

function getApprovalStatusType(status: string): string {
  const typeMap: Record<string, string> = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger',
  }
  return typeMap[status] || 'info'
}

async function loadUsers() {
  try {
    const [allUsers, pendingApps] = await Promise.all([
      userApi.getUsers(),
      userApi.getPendingApplications(),
    ])
    users.value = allUsers.data
    pendingUsers.value = pendingApps.data
  } catch (error) {
    // Error handled by request interceptor
  }
}

function openCreateDialog() {
  isEdit.value = false
  currentUser.value = null
  Object.assign(formData, {
    username: '',
    password: 'password123',
    name: '',
    role: '',
    phone: '',
  })
  dialogVisible.value = true
}

function openEditDialog(user: User) {
  isEdit.value = true
  currentUser.value = user
  Object.assign(formData, {
    username: user.username,
    password: '',
    name: user.name,
    role: user.role,
    phone: user.phone,
  })
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      if (isEdit.value && currentUser.value) {
        const updateData: UpdateUserRequest = {
          name: formData.name,
          role: formData.role,
          phone: formData.phone,
        }
        await userApi.updateUser(currentUser.value.id, updateData)
        ElMessage.success('用户更新成功')
      } else {
        await userApi.createUser(formData)
        ElMessage.success('用户创建成功')
      }
      dialogVisible.value = false
      loadUsers()
    } catch (error) {
      // Error handled by request interceptor
    } finally {
      submitting.value = false
    }
  })
}

async function handleApprove(user: User, approve: boolean) {
  if (approve) {
    try {
      await userApi.approveApplication(user.id, { approve: true })
      ElMessage.success('已通过申请')
      loadUsers()
    } catch (error) {
      // Error handled by request interceptor
    }
  } else {
    ElMessageBox.prompt('请输入拒绝原因', '拒绝申请', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPattern: /.+/,
      inputErrorMessage: '请输入拒绝原因',
    }).then(async ({ value }) => {
      try {
        await userApi.approveApplication(user.id, { approve: false, reason: value })
        ElMessage.success('已拒绝申请')
        loadUsers()
      } catch (error) {
        // Error handled by request interceptor
      }
    }).catch(() => {
      // User cancelled
    })
  }
}

function handleReject(user: User) {
  handleApprove(user, false)
}

async function handleToggleActive(user: User) {
  const action = user.is_active ? '禁用' : '启用'
  try {
    await ElMessageBox.confirm(`确定要${action}用户 ${user.name} 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await userApi.updateUser(user.id, { is_active: !user.is_active })
    ElMessage.success(`${action}成功`)
    loadUsers()
  } catch (error) {
    // User cancelled or error
  }
}

onMounted(() => {
  loadUsers()
})
</script>

<style scoped>
.user-management {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header-title h2 {
  font-size: 24px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 4px 0;
}

.header-title p {
  font-size: 14px;
  color: var(--color-text-tertiary);
  margin: 0;
}

.section-card {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 24px;
  margin-bottom: 24px;
}

.pending-section {
  border-left: 4px solid var(--color-warning);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
  display: flex;
  align-items: center;
}

:deep(.el-table) {
  background: transparent;
}

:deep(.el-table__inner-wrapper) {
  background: transparent;
}

:deep(.el-table__body-wrapper) {
  background: transparent;
}

:deep(.el-table th) {
  background: var(--color-bg-secondary) !important;
  color: var(--color-text-secondary);
  font-weight: 600;
  border-color: var(--color-border);
}

:deep(.el-table td) {
  border-color: var(--color-border);
}

:deep(.el-table__row) {
  background: transparent;
}

:deep(.el-table__row:hover) {
  background: var(--color-bg-secondary) !important;
}

:deep(.el-dialog) {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
}

:deep(.el-dialog__header) {
  border-bottom: 1px solid var(--color-border);
}

:deep(.el-dialog__title) {
  color: var(--color-text-primary);
}

:deep(.el-dialog__body) {
  padding: 24px;
}

:deep(.el-form-item__label) {
  color: var(--color-text-secondary);
}

:deep(.el-input__wrapper) {
  background: var(--color-bg-secondary);
  border-color: var(--color-border);
}

:deep(.el-input__wrapper:hover) {
  border-color: var(--color-primary-dim);
}

:deep(.el-input__wrapper.is-focus) {
  border-color: var(--color-primary);
}
</style>
