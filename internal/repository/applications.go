package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/make-smart-products/requests-api/internal/model"
)

type ApplicationFilter struct {
	UserID              *int64
	AssignedManagerID   *int64
	IncludeUnassigned   bool
	Status              *model.ApplicationStatus
}

func (s *Store) CreateApplication(app *model.Application) (*model.Application, error) {
	if app.Title == "" || !app.Status.IsValid() {
		return nil, ErrInvalidInput
	}

	now := time.Now().UTC()
	result, err := s.db.Exec(`
		INSERT INTO applications (user_id, assigned_manager_id, title, description, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, app.UserID, app.AssignedManagerID, app.Title, app.Description, app.Status, now, now)
	if err != nil {
		return nil, fmt.Errorf("insert application: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("last insert id: %w", err)
	}

	return s.GetApplication(id)
}

func (s *Store) GetApplication(id int64) (*model.Application, error) {
	row := s.db.QueryRow(`
		SELECT id, user_id, assigned_manager_id, title, description, status, created_at, updated_at
		FROM applications WHERE id = ?
	`, id)
	return scanApplication(row)
}

func (s *Store) ListApplications(filter ApplicationFilter) ([]model.Application, error) {
	query := `
		SELECT id, user_id, assigned_manager_id, title, description, status, created_at, updated_at
		FROM applications WHERE 1=1
	`
	args := make([]any, 0)

	if filter.UserID != nil {
		query += ` AND user_id = ?`
		args = append(args, *filter.UserID)
	}
	if filter.AssignedManagerID != nil {
		if filter.IncludeUnassigned {
			query += ` AND (assigned_manager_id = ? OR assigned_manager_id IS NULL)`
			args = append(args, *filter.AssignedManagerID)
		} else {
			query += ` AND assigned_manager_id = ?`
			args = append(args, *filter.AssignedManagerID)
		}
	}
	if filter.Status != nil {
		query += ` AND status = ?`
		args = append(args, *filter.Status)
	}

	query += ` ORDER BY id DESC`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list applications: %w", err)
	}
	defer rows.Close()

	apps := make([]model.Application, 0)
	for rows.Next() {
		app, err := scanApplication(rows)
		if err != nil {
			return nil, err
		}
		apps = append(apps, *app)
	}

	return apps, rows.Err()
}

func (s *Store) UpdateApplication(app *model.Application) (*model.Application, error) {
	if app.ID == 0 || !app.Status.IsValid() {
		return nil, ErrInvalidInput
	}

	now := time.Now().UTC()
	result, err := s.db.Exec(`
		UPDATE applications
		SET assigned_manager_id = ?, title = ?, description = ?, status = ?, updated_at = ?
		WHERE id = ?
	`, app.AssignedManagerID, app.Title, app.Description, app.Status, now, app.ID)
	if err != nil {
		return nil, fmt.Errorf("update application: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return nil, ErrNotFound
	}

	return s.GetApplication(app.ID)
}

func (s *Store) DeleteApplication(id int64) error {
	result, err := s.db.Exec(`DELETE FROM applications WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete application: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func scanApplication(scanner interface {
	Scan(dest ...any) error
}) (*model.Application, error) {
	var app model.Application
	var assignedManagerID sql.NullInt64
	var createdAt, updatedAt string

	err := scanner.Scan(
		&app.ID, &app.UserID, &assignedManagerID, &app.Title, &app.Description, &app.Status, &createdAt, &updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan application: %w", err)
	}

	if assignedManagerID.Valid {
		value := assignedManagerID.Int64
		app.AssignedManagerID = &value
	}

	app.CreatedAt = parseTime(createdAt)
	app.UpdatedAt = parseTime(updatedAt)
	return &app, nil
}

func IsTransitionAllowed(current, next model.ApplicationStatus, role model.Role) bool {
	if current == next {
		return true
	}

	switch role {
	case model.RoleClient:
		if current == model.StatusDraft && next == model.StatusSubmitted {
			return true
		}
		return false
	case model.RoleManager:
		allowed := map[model.ApplicationStatus][]model.ApplicationStatus{
			model.StatusSubmitted: {model.StatusInReview, model.StatusRejected},
			model.StatusInReview:  {model.StatusApproved, model.StatusRejected},
		}
		for _, status := range allowed[current] {
			if status == next {
				return true
			}
		}
		return false
	case model.RoleAdmin:
		return next.IsValid()
	default:
		return false
	}
}

func NormalizeStatus(value string) (model.ApplicationStatus, error) {
	status := model.ApplicationStatus(strings.TrimSpace(value))
	if !status.IsValid() {
		return "", ErrInvalidInput
	}
	return status, nil
}
