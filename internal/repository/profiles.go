package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/make-smart-products/requests-api/internal/model"
)

func (s *Store) UpsertProfile(profile *model.Profile) (*model.Profile, error) {
	now := time.Now().UTC()
	_, err := s.db.Exec(`
		INSERT INTO profiles (user_id, first_name, last_name, phone, company, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET
			first_name = excluded.first_name,
			last_name = excluded.last_name,
			phone = excluded.phone,
			company = excluded.company,
			updated_at = excluded.updated_at
	`, profile.UserID, profile.FirstName, profile.LastName, profile.Phone, profile.Company, now, now)
	if err != nil {
		return nil, fmt.Errorf("upsert profile: %w", err)
	}

	return s.GetProfile(profile.UserID)
}

func (s *Store) GetProfile(userID int64) (*model.Profile, error) {
	row := s.db.QueryRow(`
		SELECT user_id, first_name, last_name, phone, company, created_at, updated_at
		FROM profiles WHERE user_id = ?
	`, userID)

	var profile model.Profile
	var createdAt, updatedAt string
	err := row.Scan(&profile.UserID, &profile.FirstName, &profile.LastName, &profile.Phone, &profile.Company, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}

	profile.CreatedAt = parseTime(createdAt)
	profile.UpdatedAt = parseTime(updatedAt)
	return &profile, nil
}
