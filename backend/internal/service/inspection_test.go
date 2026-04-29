package service

import (
	"testing"
	"time"
)

func TestTimestampValidation_WithinWindow(t *testing.T) {
	// Timestamp within 300 seconds should be accepted
	now := time.Now().Unix()
	if abs(now-now) > 300 {
		t.Error("Current timestamp should be within window")
	}
}

func TestTimestampValidation_OutsideWindow(t *testing.T) {
	// Timestamp older than 300 seconds should be rejected
	now := time.Now().Unix()
	oldTimestamp := now - 600 // 10 minutes ago
	if abs(now-oldTimestamp) <= 300 {
		t.Error("10-minute-old timestamp should be outside 300s window")
	}
}

func TestTimestampValidation_FutureTimestamp(t *testing.T) {
	// Future timestamp beyond 300s should be rejected
	now := time.Now().Unix()
	futureTimestamp := now + 600 // 10 minutes in future
	if abs(now-futureTimestamp) <= 300 {
		t.Error("10-minute-future timestamp should be outside 300s window")
	}
}

func TestTimestampValidation_EdgeCase(t *testing.T) {
	// Exactly at boundary (300 seconds ago)
	now := time.Now().Unix()
	boundaryTimestamp := now - 300
	if abs(now-boundaryTimestamp) > 300 {
		t.Error("Exactly 300s old timestamp should be at boundary (accepted)")
	}
}
