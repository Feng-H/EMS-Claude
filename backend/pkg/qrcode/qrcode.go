package qrcode

import (
	"fmt"
	"net/url"

	"github.com/ems/backend/pkg/config"
	"github.com/skip2/go-qrcode"
)

const (
	// QRCodeSize is the size of QR code in pixels
	DefaultQRCodeSize = 300
	// QRCodeRecoveryLevel is the error correction level
	QRCodeRecoveryLevel = qrcode.Medium
)

// Generate generates a QR code for the given equipment code
func Generate(equipmentCode string) ([]byte, error) {
	// Create the URL for this equipment
	qrURL := fmt.Sprintf("%s/equipment/%s", config.Cfg.App.QRCodeBaseURL, equipmentCode)

	// Generate QR code
	qr, err := qrcode.New(qrURL, QRCodeRecoveryLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Convert to PNG
	return qr PNG(DefaultQRCodeSize)
}

// GenerateWithSize generates a QR code with custom size
func GenerateWithSize(equipmentCode string, size int) ([]byte, error) {
	qrURL := fmt.Sprintf("%s/equipment/%s", config.Cfg.App.QRCodeBaseURL, equipmentCode)

	qr, err := qrcode.New(qrURL, QRCodeRecoveryLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	return qr.PNG(size)
}

// GenerateForQRCode generates a QR code for an equipment's QR code field
func GenerateForQRCode(qrCode string) ([]byte, error) {
	// The qr_code field in database is already a unique identifier
	// Use it directly for the QR code content
	qr, err := qrcode.New(qrCode, QRCodeRecoveryLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	return qr.PNG(DefaultQRCodeSize)
}

// GetEquipmentURL returns the full URL for an equipment
func GetEquipmentURL(equipmentCode string) string {
	return fmt.Sprintf("%s/equipment/%s", config.Cfg.App.QRCodeBaseURL, equipmentCode)
}

// ValidateQRCode validates if a QR code string is properly formatted
func ValidateQRCode(qrCode string) bool {
	if len(qrCode) < 4 || len(qrCode) > 255 {
		return false
	}

	// Basic URL validation
	_, err := url.ParseRequestURI(qrCode)
	return err == nil
}

// BatchGenerate generates QR codes for multiple equipment codes
func BatchGenerate(equipmentCodes []string) (map[string][]byte, error) {
	result := make(map[string][]byte)

	for _, code := range equipmentCodes {
		png, err := Generate(code)
		if err != nil {
			return nil, fmt.Errorf("failed to generate QR code for %s: %w", code, err)
		}
		result[code] = png
	}

	return result, nil
}
