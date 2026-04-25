package trace

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateTraceID generates a unique trace ID for agent requests
// Format: agt_YYYYMMDD_random
func GenerateTraceID() string {
	now := time.Now()
	datePart := now.Format("20060102")
	randomPart := rand.Intn(1000000)
	return fmt.Sprintf("agt_%s_%06d", datePart, randomPart)
}
