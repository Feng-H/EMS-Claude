import { ref, onMounted, onUnmounted } from 'vue'
import { isMobile, isTablet, isDesktop, getDeviceType } from '@/utils/device'

export interface DeviceInfo {
  isMobile: boolean
  isTablet: boolean
  isDesktop: boolean
  deviceType: 'mobile' | 'tablet' | 'desktop'
  width: number
}

/**
 * Reactive device detection composable
 * Updates device information when window is resized
 *
 * @example
 * ```ts
 * const device = useDeviceDetection()
 * console.log(device.isMobile) // true/false
 * console.log(device.deviceType) // 'mobile' | 'tablet' | 'desktop'
 * ```
 */
export function useDeviceDetection() {
  const device = ref<DeviceInfo>({
    isMobile: isMobile(),
    isTablet: isTablet(),
    isDesktop: isDesktop(),
    deviceType: getDeviceType(),
    width: typeof window !== 'undefined' ? window.innerWidth : 1024
  })

  const updateDevice = () => {
    if (typeof window !== 'undefined') {
      device.value = {
        isMobile: isMobile(),
        isTablet: isTablet(),
        isDesktop: isDesktop(),
        deviceType: getDeviceType(),
        width: window.innerWidth
      }
    }
  }

  onMounted(() => {
    window.addEventListener('resize', updateDevice)
  })

  onUnmounted(() => {
    window.removeEventListener('resize', updateDevice)
  })

  return {
    device
  }
}
