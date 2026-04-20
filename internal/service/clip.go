package service

import (
	"encoding/json"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/lusiker/clipper/internal/model"
	"github.com/lusiker/clipper/internal/pkg/storage"
	"github.com/lusiker/clipper/internal/repository"
)

type ClipService struct {
	clipRepo   *repository.ClipRepository
	deviceRepo *repository.DeviceRepository
}

func NewClipService(clipRepo *repository.ClipRepository, deviceRepo *repository.DeviceRepository) *ClipService {
	return &ClipService{clipRepo: clipRepo, deviceRepo: deviceRepo}
}

func (s *ClipService) Create(userID, deviceID string, data *model.ClipCreate) (*model.Clip, error) {
	clip := &model.Clip{
		ID:      uuid.New().String(),
		UserID:  userID,
		DeviceID: deviceID,
		Type:    data.Type,
		Content: data.Content,
		Meta:    data.Meta,
	}

	if err := s.clipRepo.Create(clip); err != nil {
		return nil, err
	}

	return clip, nil
}

func (s *ClipService) List(userID string, limit, offset int) ([]model.Clip, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.clipRepo.FindByUserID(userID, limit, offset)
}

func (s *ClipService) Get(userID, clipID string) (*model.Clip, error) {
	clip, err := s.clipRepo.FindByID(clipID)
	if err != nil {
		return nil, err
	}

	if clip == nil || clip.UserID != userID {
		return nil, nil
	}

	return clip, nil
}

func (s *ClipService) Delete(userID, clipID string) error {
	clip, err := s.clipRepo.FindByID(clipID)
	if err != nil {
		return err
	}

	if clip == nil || clip.UserID != userID {
		return nil
	}

	// Delete image files if clip type is image
	if clip.Type == model.ClipTypeImage {
		if err := storage.DeleteClipFiles(userID, clipID); err != nil {
			// Log but continue to delete database record
		}
	}

	return s.clipRepo.Delete(clipID)
}

func (s *ClipService) UploadImage(userID, deviceID string, file *multipart.FileHeader) (*model.Clip, error) {
	// Save image and generate thumbnail
	meta, contentPath, err := storage.SaveImage(userID, file)
	if err != nil {
		return nil, err
	}

	// Serialize meta to JSON
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}

	clip := &model.Clip{
		ID:       uuid.New().String(),
		UserID:   userID,
		DeviceID: deviceID,
		Type:     model.ClipTypeImage,
		Content:  contentPath,
		Meta:     string(metaJSON),
	}

	if err := s.clipRepo.Create(clip); err != nil {
		// Clean up uploaded files if database insert fails
		storage.DeleteClipFiles(userID, clip.ID)
		return nil, err
	}

	return clip, nil
}