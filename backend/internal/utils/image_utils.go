package utils

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// BuildImageURL constructs a full URL for an image
// Returns empty string if filename is empty or nil
func BuildImageURL(baseURL, filename string) string {
	if filename == "" {
		return ""
	}

	// Remove trailing slash from baseURL if present
	baseURL = strings.TrimSuffix(baseURL, "/")

	// Ensure filename has extension (default to .jpg if missing)
	if !strings.Contains(filename, ".") {
		filename = filename + ".png"
	}

	return fmt.Sprintf("%s/images/%s", baseURL, filename)
}

// ValidateImagePath validates that a filename is safe to use
// Prevents path traversal attacks
func ValidateImagePath(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	// Check for path traversal attempts
	if strings.Contains(filename, "..") {
		return fmt.Errorf("invalid filename: path traversal detected")
	}

	if strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return fmt.Errorf("invalid filename: path separators not allowed")
	}

	// Check for null bytes
	if strings.Contains(filename, "\x00") {
		return fmt.Errorf("invalid filename: null bytes not allowed")
	}

	return nil
}

// ValidateImageType checks if a file is a valid image type
// Only allows JPEG, PNG, and GIF
func ValidateImageType(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	contentType := http.DetectContentType(buffer)

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}

	if !allowedTypes[contentType] {
		return fmt.Errorf("invalid file type: %s (only JPEG, PNG, GIF allowed)", contentType)
	}

	return nil
}

// SanitizeFilename removes potentially dangerous characters from a filename
func SanitizeFilename(filename string) string {
	// Remove path separators
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")

	// Remove parent directory references
	filename = strings.ReplaceAll(filename, "..", "")

	// Remove null bytes
	filename = strings.ReplaceAll(filename, "\x00", "")

	return filename
}

// GetSafeFilePath constructs a safe file path within the upload directory
// Returns error if the resulting path is outside the upload directory
func GetSafeFilePath(uploadDir, filename string) (string, error) {
	// Validate filename first
	if err := ValidateImagePath(filename); err != nil {
		return "", err
	}

	// Clean the upload directory path
	uploadDir = filepath.Clean(uploadDir)

	// Construct the full path
	fullPath := filepath.Join(uploadDir, filename)

	// Clean the full path
	fullPath = filepath.Clean(fullPath)

	// Verify the path is still within the upload directory
	if !strings.HasPrefix(fullPath, uploadDir) {
		return "", fmt.Errorf("invalid path: outside upload directory")
	}

	return fullPath, nil
}
