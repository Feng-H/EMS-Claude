<template>
  <div class="spare-part-selector">
    <!-- 触发按钮 -->
    <van-field
      :value="selectedText"
      readonly
      is-link
      placeholder="选择备件"
      @click="showSelector = true"
    >
      <template #button>
        <van-tag v-if="selectedParts.length > 0" type="primary">
          已选 {{ selectedParts.length }} 项
        </van-tag>
      </template>
    </van-field>

    <!-- 已选备件列表 -->
    <div v-if="selectedParts.length > 0" class="selected-list">
      <div
        v-for="item in selectedParts"
        :key="item.id"
        class="selected-item"
      >
        <div class="item-info">
          <div class="item-name">{{ item.name }}</div>
          <div class="item-detail">
            <span class="item-code">{{ item.code }}</span>
            <span class="item-stock">库存: {{ item.stock }}</span>
          </div>
        </div>
        <div class="item-quantity">
          <van-stepper
            v-model="item.quantity"
            :min="1"
            :max="item.stock"
            integer
            @change="onQuantityChange(item)"
          />
          <van-icon
            name="cross"
            size="18"
            class="remove-icon"
            @click="removePart(item)"
          />
        </div>
      </div>
    </div>

    <!-- 选择器弹窗 -->
    <van-popup
      v-model:show="showSelector"
      position="bottom"
      :style="{ height: '80%' }"
      round
    >
      <div class="selector-popup">
        <van-nav-bar
          title="选择备件"
          left-text="取消"
          right-text="确定"
          @click-left="showSelector = false"
          @click-right="confirmSelection"
        />

        <!-- 搜索栏 -->
        <van-search
          v-model="searchKeyword"
          placeholder="搜索备件名称或编码"
          @update:model-value="onSearch"
        />

        <!-- 备件列表 -->
        <div class="parts-list">
          <van-loading v-if="loading" size="24">加载中...</van-loading>
          <van-empty v-else-if="filteredParts.length === 0" description="暂无备件" />
          <van-cell-group v-else inset>
            <van-cell
              v-for="part in filteredParts"
              :key="part.id"
              :title="part.name"
              :label="`${part.code} | 规格: ${part.specification || '无'} | 库存: ${part.stock}`"
              is-link
              @click="showPartDetail(part)"
            />
          </van-cell-group>
        </div>
      </div>
    </van-popup>

    <!-- 备件详情弹窗 -->
    <van-dialog
      v-model:show="showDetailDialog"
      :title="currentPart?.name || '备件详情'"
      :show-cancel-button="true"
      confirm-button-text="添加"
      @confirm="addPart"
    >
      <div class="part-detail" v-if="currentPart">
        <van-cell-group inset>
          <van-cell title="备件编码" :value="currentPart.code" />
          <van-cell title="规格" :value="currentPart.specification || '无'" />
          <van-cell title="当前库存" :value="String(currentPart.stock)" />
          <van-cell title="单位" :value="currentPart.unit" />
        </van-cell-group>

        <van-field
          v-model.number="tempQuantity"
          type="number"
          label="数量"
          placeholder="请输入数量"
          :min="1"
          :max="currentPart.stock"
          integer
          input-align="right"
        >
          <template #button>
            <span>{{ currentPart.unit }}</span>
          </template>
        </van-field>

        <div v-if="quantityError" class="error-message">
          {{ quantityError }}
        </div>
      </div>
    </van-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { showToast } from 'vant'
import { sparePartApi, type SparePart } from '@/api/sparepart'

interface SelectedPart extends SparePart {
  quantity: number
}

interface Props {
  factoryId?: number
  modelValue?: SelectedPart[]
}

const props = withDefaults(defineProps<Props>(), {
  factoryId: 0
})

const emit = defineEmits<{
  'update:modelValue': [parts: SelectedPart[]]
  change: [parts: SelectedPart[]]
}>()

const showSelector = ref(false)
const showDetailDialog = ref(false)
const loading = ref(false)
const searchKeyword = ref('')
const allParts = ref<SparePart[]>([])
const selectedParts = ref<SelectedPart[]>([])
const currentPart = ref<SparePart | null>(null)
const tempQuantity = ref(1)
const quantityError = ref('')

const filteredParts = computed(() => {
  if (!searchKeyword.value) {
    return allParts.value
  }
  const keyword = searchKeyword.value.toLowerCase()
  return allParts.value.filter(part =>
    part.name.toLowerCase().includes(keyword) ||
    part.code.toLowerCase().includes(keyword)
  )
})

const selectedText = computed(() => {
  if (selectedParts.value.length === 0) {
    return ''
  }
  if (selectedParts.value.length === 1) {
    return `${selectedParts.value[0].name} x${selectedParts.value[0].quantity}`
  }
  return `${selectedParts.value[0].name} 等 ${selectedParts.value.length} 项`
})

const loadSpareParts = async () => {
  if (!props.factoryId) {
    showToast('请先选择工厂')
    return
  }

  loading.value = true
  try {
    const data = await sparePartApi.getList({ factory_id: props.factoryId })
    allParts.value = data.filter(p => p.stock > 0)
  } catch (error: any) {
    showToast(error.message || '加载备件失败')
  } finally {
    loading.value = false
  }
}

const onSearch = () => {
  // 搜索逻辑在computed中处理
}

const showPartDetail = (part: SparePart) => {
  currentPart.value = part
  const existing = selectedParts.value.find(p => p.id === part.id)
  tempQuantity.value = existing ? existing.quantity : 1
  quantityError.value = ''
  showDetailDialog.value = true
}

const validateQuantity = (): boolean => {
  if (!currentPart.value) return false

  if (!tempQuantity.value || tempQuantity.value < 1) {
    quantityError.value = '数量必须大于0'
    return false
  }

  if (tempQuantity.value > currentPart.value.stock) {
    quantityError.value = `库存不足，当前库存: ${currentPart.value.stock}`
    return false
  }

  quantityError.value = ''
  return true
}

const addPart = () => {
  if (!currentPart.value || !validateQuantity()) {
    return
  }

  const existingIndex = selectedParts.value.findIndex(p => p.id === currentPart.value!.id)
  if (existingIndex >= 0) {
    selectedParts.value[existingIndex].quantity = tempQuantity.value
  } else {
    selectedParts.value.push({
      ...currentPart.value,
      quantity: tempQuantity.value
    })
  }

  showDetailDialog.value = false
  showToast('已添加')
}

const removePart = (part: SelectedPart) => {
  const index = selectedParts.value.findIndex(p => p.id === part.id)
  if (index >= 0) {
    selectedParts.value.splice(index, 1)
  }
}

const onQuantityChange = (part: SelectedPart) => {
  emitChange()
}

const confirmSelection = () => {
  emitChange()
  showSelector.value = false
}

const emitChange = () => {
  emit('update:modelValue', selectedParts.value)
  emit('change', selectedParts.value)
}

onMounted(() => {
  loadSpareParts()
})

// 暴露方法供父组件调用
defineExpose({
  loadSpareParts,
  clear: () => {
    selectedParts.value = []
    emitChange()
  }
})
</script>

<style scoped>
.spare-part-selector {
  margin-bottom: 16px;
}

.selected-list {
  margin-top: 12px;
  background: #f5f5f5;
  border-radius: 8px;
  padding: 8px;
}

.selected-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: #fff;
  border-radius: 8px;
  margin-bottom: 8px;
}

.selected-item:last-child {
  margin-bottom: 0;
}

.item-info {
  flex: 1;
}

.item-name {
  font-weight: 500;
  margin-bottom: 4px;
}

.item-detail {
  font-size: 12px;
  color: #999;
}

.item-code {
  margin-right: 12px;
}

.item-stock {
  color: #07c160;
}

.item-quantity {
  display: flex;
  align-items: center;
  gap: 12px;
}

.remove-icon {
  color: #999;
  cursor: pointer;
}

.selector-popup {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.parts-list {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.part-detail {
  padding: 16px;
}

.error-message {
  padding: 8px 16px;
  color: #ee0a24;
  font-size: 14px;
}
</style>
