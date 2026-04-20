package storage

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/lusiker/clipper/internal/config"
)

const (
	MaxImageSize    = 20 * 1024 * 1024 // 20MB
	ThumbnailWidth  = 300
	ThumbnailQuality = 85
)

var allowedFormats = map[string]bool{
	"jpeg": true,
	"jpg":  true,
	"png":  true,
	"gif":  true,
	"webp": true,
}

type ImageMeta struct {
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Size      int64  `json:"size"`
	Format    string `json:"format"`
	ThumbPath string `json:"thumb_path"`
}

func getUploadDir() string {
	return config.GetUploadDir()
}

func EnsureUploadDir() error {
	uploadDir := getUploadDir()
	return os.MkdirAll(uploadDir, 0755)
}

func getUserDir(userID string) string {
	return filepath.Join(getUploadDir(), userID)
}

func ensureUserDir(userID string) error {
	return os.MkdirAll(getUserDir(userID), 0755)
}

func SaveImage(userID string, file *multipart.FileHeader) (*ImageMeta, string, error) {
	// Validate file size
	if file.Size > MaxImageSize {
		return nil, "", fmt.Errorf("image size exceeds 20MB limit")
	}

	// Extract format from filename
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file.Filename), "."))
	if !allowedFormats[ext] {
		return nil, "", fmt.Errorf("unsupported image format: %s", ext)
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Decode image to get dimensions and validate it's a real image
	img, format, err := image.Decode(src)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Reset reader position after decode
	src.Seek(0, 0)

	// Generate clip ID
	clipID := uuid.New().String()

	// Ensure user directory exists
	if err := ensureUserDir(userID); err != nil {
		return nil, "", fmt.Errorf("failed to create user directory: %w", err)
	}

	userDir := getUserDir(userID)

	// Save original image
	origPath := filepath.Join(userDir, fmt.Sprintf("%s_orig.%s", clipID, ext))
	origFile, err := os.Create(origPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create original file: %w", err)
	}

	if _, err := io.Copy(origFile, src); err != nil {
		origFile.Close()
		return nil, "", fmt.Errorf("failed to save original image: %w", err)
	}
	origFile.Close()

	// Generate thumbnail
	thumbImg := imaging.Resize(img, ThumbnailWidth, 0, imaging.Lanczos)
	thumbPath := filepath.Join(userDir, fmt.Sprintf("%s_thumb.jpg", clipID))
	thumbFile, err := os.Create(thumbPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create thumbnail file: %w", err)
	}

	if err := jpeg.Encode(thumbFile, thumbImg, &jpeg.Options{Quality: ThumbnailQuality}); err != nil {
		thumbFile.Close()
		return nil, "", fmt.Errorf("failed to encode thumbnail: %w", err)
	}
	thumbFile.Close()

	// Build metadata
	meta := &ImageMeta{
		Width:     img.Bounds().Dx(),
		Height:    img.Bounds().Dy(),
		Size:      file.Size,
		Format:    format,
		ThumbPath: fmt.Sprintf("%s/%s_thumb.jpg", userID, clipID),
	}

	// Relative path for content field
	contentPath := fmt.Sprintf("%s/%s_orig.%s", userID, clipID, ext)

	return meta, contentPath, nil
}

func DeleteClipFiles(userID, clipID string) error {
	userDir := getUserDir(userID)

	// Find and delete all files matching clipID pattern
	files, err := filepath.Glob(filepath.Join(userDir, clipID+"_*"))
	if err != nil {
		return fmt.Errorf("failed to list clip files: %w", err)
	}

	for _, file := range files {
		if err := os.Remove(file); err != nil {
			// Log but continue to try deleting other files
			fmt.Printf("Warning: failed to delete file %s: %v\n", file, err)
		}
	}

	return nil
}

func ParseMeta(metaJSON string) (*ImageMeta, error) {
	if metaJSON == "" {
		return nil, nil
	}

	var meta ImageMeta
	if err := json.Unmarshal([]byte(metaJSON), &meta); err != nil {
		return nil, fmt.Errorf("failed to parse image meta: %w", err)
	}

	return &meta, nil
}

func GetImagePaths(userID, clipID, metaJSON string) (origPath, thumbPath string, err error) {
	userDir := getUserDir(userID)

	meta, err := ParseMeta(metaJSON)
	if err != nil {
		return "", "", err
	}

	if meta != nil {
		thumbPath = filepath.Join(userDir, fmt.Sprintf("%s_thumb.jpg", clipID))
	}

	// Try to find original file with any extension
	files, err := filepath.Glob(filepath.Join(userDir, clipID+"_orig.*"))
	if err != nil {
		return "", "", err
	}

	if len(files) > 0 {
		origPath = files[0]
	}

	return origPath, thumbPath, nil
}