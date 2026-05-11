package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/internal/agent/repository"
	"github.com/ledongthuc/pdf"
)

type ManualService struct {
	agentRepo repository.IAgentRepository
}

func NewManualService(agentRepo repository.IAgentRepository) *ManualService {
	return &ManualService{
		agentRepo: agentRepo,
	}
}

func (s *ManualService) UploadAndProcess(file io.Reader, filename string, title string, equipmentTypeID *uint) (*model.ManualDocument, error) {
	// 1. Save file
	saveDir := "./uploads/manuals"
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}
	
	filePath := filepath.Join(saveDir, fmt.Sprintf("%d_%s", os.Getpid(), filename))
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	// 2. Create Document record
	doc := &model.ManualDocument{
		Title:           title,
		EquipmentTypeID: equipmentTypeID,
		FilePath:        filePath,
	}
	if err := s.agentRepo.CreateManualDocument(doc); err != nil {
		return nil, fmt.Errorf("failed to create document record: %v", err)
	}

	// 3. Extract Text and Chunk (Async)
	go s.processPDF(doc)

	return doc, nil
}

func (s *ManualService) processPDF(doc *model.ManualDocument) {
	f, r, err := pdf.Open(doc.FilePath)
	if err != nil {
		fmt.Printf("Error opening PDF %s: %v\n", doc.FilePath, err)
		return
	}
	defer f.Close()

	var chunks []model.ManualChunk
	pageSize := r.NumPage()

	for i := 1; i <= pageSize; i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}
		
		text, err := p.GetPlainText(nil)
		if err != nil {
			fmt.Printf("Error extracting text from page %d: %v\n", i, err)
			continue
		}

		// Sliding window chunking
		pageChunks := s.splitIntoChunks(text, 800, 150)
		for _, content := range pageChunks {
			if len(strings.TrimSpace(content)) < 10 {
				continue
			}
			chunks = append(chunks, model.ManualChunk{
				DocumentID:   doc.ID,
				SectionTitle: fmt.Sprintf("%s - 第 %d 页", doc.Title, i),
				Content:      content,
				PageNumber:   i,
			})
		}
		
		// Batch save
		if len(chunks) >= 20 {
			if err := s.agentRepo.CreateManualChunks(chunks); err != nil {
				fmt.Printf("Error saving chunks: %v\n", err)
			}
			chunks = []model.ManualChunk{}
		}
	}

	if len(chunks) > 0 {
		if err := s.agentRepo.CreateManualChunks(chunks); err != nil {
			fmt.Printf("Error saving final chunks: %v\n", err)
		}
	}
}

func (s *ManualService) splitIntoChunks(text string, chunkSize int, overlap int) []string {
	runes := []rune(text)
	if len(runes) <= chunkSize {
		return []string{text}
	}

	var chunks []string
	step := chunkSize - overlap
	if step <= 0 {
		step = chunkSize / 2
	}

	for i := 0; i < len(runes); i += step {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
		if end == len(runes) {
			break
		}
	}
	return chunks
}
