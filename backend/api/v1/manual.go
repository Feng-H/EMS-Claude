package v1

import (
	"net/http"
	"strconv"

	agentRepo "github.com/ems/backend/internal/agent/repository"
	"github.com/ems/backend/internal/service"
	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/pkg/database"
	"github.com/gin-gonic/gin"
)

var (
	manualService *service.ManualService
)

func InitManual() {
	var repo agentRepo.IAgentRepository
	if config.Cfg.Storage.Mode == "memory" {
		repo = agentRepo.NewMemoryAgentRepository()
	} else {
		repo = agentRepo.NewDBAgentRepository(database.GetDB())
	}
	manualService = service.NewManualService(repo)
}

// UploadManual handles PDF upload and chunking
// @Summary Upload technical manual (PDF)
// @Tags knowledge
// @Accept multipart/form-data
// @Param file formData file true "Technical manual PDF"
// @Param title formData string false "Manual title"
// @Param equipment_type_id formData uint false "Associated equipment type ID"
// @Success 200 {object} gin.H
// @Router /knowledge/manual/upload [post]
func UploadManual(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	title := c.PostForm("title")
	if title == "" {
		title = header.Filename
	}

	equipmentTypeIDStr := c.PostForm("equipment_type_id")
	var equipmentTypeID *uint
	if equipmentTypeIDStr != "" {
		id, err := strconv.ParseUint(equipmentTypeIDStr, 10, 32)
		if err == nil {
			uid := uint(id)
			equipmentTypeID = &uid
		}
	}

	doc, err := manualService.UploadAndProcess(file, header.Filename, title, equipmentTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "上传成功，后台处理中",
		"document_id": doc.ID,
		"title":       doc.Title,
	})
}
