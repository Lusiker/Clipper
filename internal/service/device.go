package service

import (
	"github.com/google/uuid"
	"github.com/lusiker/clipper/internal/model"
	"github.com/lusiker/clipper/internal/repository"
)

type DeviceService struct {
	deviceRepo *repository.DeviceRepository
}

func NewDeviceService(deviceRepo *repository.DeviceRepository) *DeviceService {
	return &DeviceService{deviceRepo: deviceRepo}
}

func (s *DeviceService) Register(userID, deviceID, name, ip string) (*model.Device, error) {
	device, _ := s.deviceRepo.FindByID(deviceID)

	if device == nil {
		device = &model.Device{
			ID:       deviceID,
			UserID:   userID,
			Name:     name,
			IP:       ip,
			IsOnline: true,
		}
		if err := s.deviceRepo.Create(device); err != nil {
			return nil, err
		}
	} else {
		if device.UserID != userID {
			return nil, nil
		}
		if err := s.deviceRepo.UpdateOnlineStatus(deviceID, true); err != nil {
			return nil, err
		}
	}

	return device, nil
}

func (s *DeviceService) List(userID string) ([]model.Device, error) {
	return s.deviceRepo.FindByUserID(userID)
}

func (s *DeviceService) SetOffline(deviceID string) error {
	return s.deviceRepo.UpdateOnlineStatus(deviceID, false)
}

func (s *DeviceService) UpdateLastSeen(deviceID string) error {
	return s.deviceRepo.UpdateLastSeen(deviceID)
}

func (s *DeviceService) CreateIfNotExists(userID, name, ip string) (*model.Device, error) {
	device := &model.Device{
		ID:       uuid.New().String(),
		UserID:   userID,
		Name:     name,
		IP:       ip,
		IsOnline: true,
	}

	if err := s.deviceRepo.Create(device); err != nil {
		return nil, err
	}

	return device, nil
}