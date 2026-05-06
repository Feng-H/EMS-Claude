<template>
  <div class="qr-scanner">
    <div class="video-container">
      <video ref="videoRef" class="video" playsinline></video>
      <canvas ref="canvasRef" class="canvas" v-show="false"></canvas>
      
      <!-- 扫描框 UI -->
      <div class="scan-overlay">
        <div class="scan-frame">
          <div class="corner top-left"></div>
          <div class="corner top-right"></div>
          <div class="corner bottom-left"></div>
          <div class="corner bottom-right"></div>
          <div class="scan-line"></div>
        </div>
        <p class="scan-hint">请将设备二维码置于框内</p>
      </div>
    </div>
    
    <div class="scanner-actions">
      <van-button icon="photograph" plain round type="primary" @click="chooseImage">
        从相册选择
      </van-button>
      <input
        type="file"
        ref="fileInputRef"
        accept="image/*"
        class="hidden-input"
        @change="onFileChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import jsQR from 'jsqr'
import { showToast, showLoadingToast, closeToast } from 'vant'

const props = defineProps<{
  active: boolean
}>()

const emit = defineEmits<{
  (e: 'success', code: string): void
  (e: 'error', message: string): void
}>()

const videoRef = ref<HTMLVideoElement | null>(null)
const canvasRef = ref<HTMLCanvasElement | null>(null)
const fileInputRef = ref<HTMLInputElement | null>(null)
let stream: MediaStream | null = null
let animationFrame: number | null = null

const startCamera = async () => {
  try {
    if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
      throw new Error('当前浏览器不支持摄像头访问')
    }

    stream = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: 'environment' }
    })

    if (videoRef.value) {
      videoRef.value.srcObject = stream
      videoRef.value.setAttribute('playsinline', 'true') // iOS fix
      videoRef.value.play()
      requestAnimationFrame(tick)
    }
  } catch (err: any) {
    console.error('Camera error:', err)
    emit('error', err.message || '无法访问摄像头')
  }
}

const stopCamera = () => {
  if (animationFrame) {
    cancelAnimationFrame(animationFrame)
    animationFrame = null
  }
  if (stream) {
    stream.getTracks().forEach(track => track.stop())
    stream = null
  }
}

const tick = () => {
  if (!props.active || !videoRef.value || !canvasRef.value) {
    animationFrame = requestAnimationFrame(tick)
    return
  }

  if (videoRef.value.readyState === videoRef.value.HAVE_ENOUGH_DATA) {
    const canvas = canvasRef.value
    const video = videoRef.value
    
    // 降低采样频率以减轻 CPU 负担，每 4 帧处理一次
    frameCount++
    if (frameCount % 4 === 0) {
      canvas.height = video.videoHeight
      canvas.width = video.videoWidth
      const ctx = canvas.getContext('2d', { willReadFrequently: true })
      
      if (ctx) {
        ctx.drawImage(video, 0, 0, canvas.width, canvas.height)
        const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height)
        const code = jsQR(imageData.data, imageData.width, imageData.height, {
          inversionAttempts: 'attemptBoth', // 尝试识别反色二维码
        })

        if (code && code.data) {
          console.log('QR Code detected:', code.data)
          // 震动反馈 (如果支持)
          if (navigator.vibrate) {
            navigator.vibrate(200)
          }
          emit('success', code.data)
          return 
        }
      }
    }
  }
  animationFrame = requestAnimationFrame(tick)
}

let frameCount = 0

const chooseImage = () => {
  fileInputRef.value?.click()
}

const onFileChange = (event: Event) => {
  const input = event.target as HTMLInputElement
  if (!input.files?.length) return

  const file = input.files[0]
  const reader = new FileReader()
  
  showLoadingToast({
    message: '正在识别...',
    forbidClick: true,
  })

  reader.onload = (e) => {
    const img = new Image()
    img.onload = () => {
      if (!canvasRef.value) return
      const canvas = canvasRef.value
      const ctx = canvas.getContext('2d')
      if (!ctx) return

      canvas.width = img.width
      canvas.height = img.height
      ctx.drawImage(img, 0, 0)
      
      const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height)
      const code = jsQR(imageData.data, imageData.width, imageData.height)
      
      closeToast()
      if (code && code.data) {
        emit('success', code.data)
      } else {
        showToast('未能识别二维码，请重试')
      }
    }
    img.src = e.target?.result as string
  }
  reader.readAsDataURL(file)
}

onMounted(() => {
  if (props.active) {
    startCamera()
  }
})

onUnmounted(() => {
  stopCamera()
})

defineExpose({
  stopCamera,
  startCamera
})
</script>

<style scoped>
.qr-scanner {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  background: #000;
}

.video-container {
  position: relative;
  width: 100%;
  aspect-ratio: 3/4;
  overflow: hidden;
  background: #000;
}

.video {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.scan-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}

.scan-frame {
  position: relative;
  width: 200px;
  height: 200px;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.corner {
  position: absolute;
  width: 20px;
  height: 20px;
  border-color: #1989fa;
  border-style: solid;
}

.top-left { top: -2px; left: -2px; border-width: 4px 0 0 4px; }
.top-right { top: -2px; right: -2px; border-width: 4px 4px 0 0; }
.bottom-left { bottom: -2px; left: -2px; border-width: 0 0 4px 4px; }
.bottom-right { bottom: -2px; right: -2px; border-width: 0 4px 4px 0; }

.scan-line {
  position: absolute;
  top: 0;
  left: 5%;
  width: 90%;
  height: 2px;
  background: linear-gradient(to right, transparent, #1989fa, transparent);
  animation: scan 2s linear infinite;
}

@keyframes scan {
  from { top: 0; }
  to { top: 100%; }
}

.scan-hint {
  color: #fff;
  font-size: 14px;
  margin-top: 24px;
  text-shadow: 0 1px 2px rgba(0,0,0,0.5);
}

.scanner-actions {
  padding: 20px;
  width: 100%;
  display: flex;
  justify-content: center;
  background: #fff;
}

.hidden-input {
  display: none;
}
</style>
