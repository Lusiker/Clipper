package repository

import (
	"database/sql"
	"time"

	"github.com/lusiker/clipper/internal/model"
)

type ClipRepository struct {
	db *sql.DB
}

func NewClipRepository(db *sql.DB) *ClipRepository {
	return &ClipRepository{db: db}
}

func (r *ClipRepository) Create(clip *model.Clip) error {
	query := `INSERT INTO clips (id, user_id, device_id, type, content, meta, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	clip.CreatedAt = now
	clip.UpdatedAt = now

	_, err := r.db.Exec(query, clip.ID, clip.UserID, clip.DeviceID, clip.Type, clip.Content, clip.Meta, clip.CreatedAt, clip.UpdatedAt)
	return err
}

func (r *ClipRepository) FindByUserID(userID string, limit, offset int) ([]model.Clip, error) {
	query := `SELECT id, user_id, device_id, type, content, meta, created_at, updated_at
		FROM clips WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clips []model.Clip
	for rows.Next() {
		var clip model.Clip
		var meta sql.NullString
		err := rows.Scan(
			&clip.ID, &clip.UserID, &clip.DeviceID, &clip.Type, &clip.Content, &meta,
			&clip.CreatedAt, &clip.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if meta.Valid {
			clip.Meta = meta.String
		}
		clips = append(clips, clip)
	}

	return clips, nil
}

func (r *ClipRepository) FindByID(id string) (*model.Clip, error) {
	query := `SELECT id, user_id, device_id, type, content, meta, created_at, updated_at FROM clips WHERE id = ?`

	clip := &model.Clip{}
	var meta sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&clip.ID, &clip.UserID, &clip.DeviceID, &clip.Type, &clip.Content, &meta,
		&clip.CreatedAt, &clip.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if meta.Valid {
		clip.Meta = meta.String
	}

	return clip, nil
}

func (r *ClipRepository) Delete(id string) error {
	query := `DELETE FROM clips WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *ClipRepository) DeleteByUserID(userID string) error {
	query := `DELETE FROM clips WHERE user_id = ?`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *ClipRepository) Count(userID string) (int, error) {
	query := `SELECT COUNT(*) FROM clips WHERE user_id = ?`
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}