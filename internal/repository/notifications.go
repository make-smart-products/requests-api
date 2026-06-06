package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/make-smart-products/requests-api/internal/model"
)

func (s *Store) CreateNotification(notification *model.Notification) (*model.Notification, error) {
	now := time.Now().UTC()
	result, err := s.db.Exec(`
		INSERT INTO notifications (user_id, channel, type, title, body, is_read, sent_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, notification.UserID, notification.Channel, notification.Type, notification.Title, notification.Body, boolToInt(notification.IsRead), nullableTime(notification.SentAt), now)
	if err != nil {
		return nil, fmt.Errorf("insert notification: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("last insert id: %w", err)
	}

	return s.GetNotification(id)
}

func (s *Store) GetNotification(id int64) (*model.Notification, error) {
	row := s.db.QueryRow(`
		SELECT id, user_id, channel, type, title, body, is_read, sent_at, created_at
		FROM notifications WHERE id = ?
	`, id)
	return scanNotification(row)
}

func (s *Store) ListNotifications(userID int64, unreadOnly bool) ([]model.Notification, error) {
	query := `
		SELECT id, user_id, channel, type, title, body, is_read, sent_at, created_at
		FROM notifications WHERE user_id = ?
	`
	args := []any{userID}

	if unreadOnly {
		query += ` AND is_read = 0`
	}
	query += ` ORDER BY id DESC`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()

	items := make([]model.Notification, 0)
	for rows.Next() {
		item, err := scanNotification(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, rows.Err()
}

func (s *Store) MarkNotificationRead(id, userID int64) (*model.Notification, error) {
	result, err := s.db.Exec(`UPDATE notifications SET is_read = 1 WHERE id = ? AND user_id = ?`, id, userID)
	if err != nil {
		return nil, fmt.Errorf("mark notification read: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return nil, ErrNotFound
	}

	return s.GetNotification(id)
}

func scanNotification(scanner interface {
	Scan(dest ...any) error
}) (*model.Notification, error) {
	var item model.Notification
	var isRead int
	var sentAt sql.NullString
	var createdAt string

	err := scanner.Scan(&item.ID, &item.UserID, &item.Channel, &item.Type, &item.Title, &item.Body, &isRead, &sentAt, &createdAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan notification: %w", err)
	}

	item.IsRead = isRead == 1
	if sentAt.Valid {
		parsed := parseTime(sentAt.String)
		item.SentAt = &parsed
	}
	item.CreatedAt = parseTime(createdAt)
	return &item, nil
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func nullableTime(value *time.Time) any {
	if value == nil {
		return nil
	}
	return value.UTC()
}
