package repository

import (
	"database/sql"
	"time"

	"github.com/lusiker/clipper/internal/model"
)

type DeviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) Create(device *model.Device) error {
	query := `INSERT INTO devices (id, user_id, name, ip, last_seen, is_online)
		VALUES (?, ?, ?, ?, ?, ?)`

	device.LastSeen = time.Now()

	_, err := r.db.Exec(query, device.ID, device.UserID, device.Name, device.IP, device.LastSeen, device.IsOnline)
	return err
}

func (r *DeviceRepository) FindByUserID(userID string) ([]model.Device, error) {
	query := `SELECT id, user_id, name, ip, last_seen, is_online FROM devices WHERE user_id = ? ORDER BY last_seen DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []model.Device
	for rows.Next() {
		var device model.Device
		var ip sql.NullString
		err := rows.Scan(
			&device.ID, &device.UserID, &device.Name, &ip, &device.LastSeen, &device.IsOnline,
		)
		if err != nil {
			return nil, err
		}
		if ip.Valid {
			device.IP = ip.String
		}
		devices = append(devices, device)
	}

	return devices, nil
}

func (r *DeviceRepository) FindByID(id string) (*model.Device, error) {
	query := `SELECT id, user_id, name, ip, last_seen, is_online FROM devices WHERE id = ?`

	device := &model.Device{}
	var ip sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&device.ID, &device.UserID, &device.Name, &ip, &device.LastSeen, &device.IsOnline,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if ip.Valid {
		device.IP = ip.String
	}

	return device, nil
}

func (r *DeviceRepository) UpdateOnlineStatus(id string, isOnline bool) error {
	query := `UPDATE devices SET is_online = ?, last_seen = ? WHERE id = ?`
	_, err := r.db.Exec(query, isOnline, time.Now(), id)
	return err
}

func (r *DeviceRepository) UpdateLastSeen(id string) error {
	query := `UPDATE devices SET last_seen = ? WHERE id = ?`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}

func (r *DeviceRepository) Delete(id string) error {
	query := `DELETE FROM devices WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *DeviceRepository) SetAllOffline(userID string) error {
	query := `UPDATE devices SET is_online = FALSE WHERE user_id = ?`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *DeviceRepository) Count(userID string) (int, error) {
	query := `SELECT COUNT(*) FROM devices WHERE user_id = ?`
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}