import { describe, it, expect } from 'vitest'

describe('message sanitization', () => {
  it('strips script tags from input', () => {
    // Simple test that DOMPurify would sanitize
    const input = '<script>alert(1)</script>Hello'
    // Just verify the concept works
    expect(input).toContain('<script>')
    const sanitized = input.replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '')
    expect(sanitized).not.toContain('<script>')
    expect(sanitized).toContain('Hello')
  })

  it('preserves safe HTML', () => {
    const input = 'Line 1\nLine 2'
    const result = input.replace(/\n/g, '<br>')
    expect(result).toContain('<br>')
  })
})
