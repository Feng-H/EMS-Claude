/**
 * Device detection utilities
 * Used to determine if the current device is mobile, tablet, or desktop
 */

/**
 * Check if the current device is a mobile device
 * Mobile: width < 768px
 */
export function isMobile(): boolean {
  if (typeof window !== 'undefined') {
    return window.innerWidth < 768
  }
  return false
}

/**
 * Check if the current device is a tablet
 * Tablet: 768px <= width < 1024px
 */
export function isTablet(): boolean {
  if (typeof window !== 'undefined') {
    const width = window.innerWidth
    return width >= 768 && width < 1024
  }
  return false
}

/**
 * Check if the current device is a desktop
 * Desktop: width >= 1024px
 */
export function isDesktop(): boolean {
  if (typeof window !== 'undefined') {
    return window.innerWidth >= 1024
  }
  return false
}

/**
 * Get the current device type
 * @returns 'mobile' | 'tablet' | 'desktop'
 */
export function getDeviceType(): 'mobile' | 'tablet' | 'desktop' {
  if (isMobile()) return 'mobile'
  if (isTablet()) return 'tablet'
  return 'desktop'
}
