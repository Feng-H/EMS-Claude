package response

import (
	"github.com/ems/backend/pkg/trace"
	"github.com/gin-gonic/gin"
)

type ErrorEnvelope struct {
	Success bool   `json:"success"`
	TraceID string `json:"trace_id"`
	Error   string `json:"error"`
}

func Error(c *gin.Context, status int, err error) {
	c.JSON(status, ErrorEnvelope{
		Success: false,
		TraceID: trace.GenerateTraceID(),
		Error:   err.Error(),
	})
}
