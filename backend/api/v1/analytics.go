package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ems/backend/internal/dto"
	"github.com/ems/backend/internal/service"
	"github.com/gin-gonic/gin"
)

var (
	analyticsService *service.AnalyticsService
)

func InitAnalytics() {
	analyticsService = service.NewAnalyticsService()
}

// GetDashboardOverview returns the dashboard overview
// @Summary Get dashboard overview
// @Tags analytics
// @Router /analytics/dashboard [get]
func GetDashboardOverview(c *gin.Context) {
	overview, err := analyticsService.GetDashboardOverview()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, overview)
}

// GetMTTRMTBF returns MTTR and MTBF statistics
// @Summary Get MTTR/MTBF
// @Tags analytics
// @Router /analytics/mttr-mtbf [get]
func GetMTTRMTBF(c *gin.Context) {
	var query dto.AnalyticsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mttrMtbf, err := analyticsService.GetMTTRMTBF(query.FactoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mttrMtbf)
}

// GetTrendData returns trend data
// @Summary Get trend data
// @Tags analytics
// @Router /analytics/trends [get]
func GetTrendData(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	trends, err := analyticsService.GetTrendData(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trends)
}

// GetFailureAnalysis returns failure analysis
// @Summary Get failure analysis
// @Tags analytics
// @Router /analytics/failures [get]
func GetFailureAnalysis(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.ParseUint(l, 10, 32); err == nil && parsed > 0 {
			limit = int(parsed)
		}
	}

	analysis, err := analyticsService.GetFailureAnalysis(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// GetTopFailureEquipment returns equipment with most failures
// @Summary Get top failure equipment
// @Tags analytics
// @Router /analytics/top-failures [get]
func GetTopFailureEquipment(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.ParseUint(l, 10, 32); err == nil && parsed > 0 {
			limit = int(parsed)
		}
	}

	equipment, err := analyticsService.GetTopFailureEquipment(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, equipment)
}
