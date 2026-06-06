package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/make-smart-products/requests-api/internal/model"
)

func (s *Store) CreateUser(email, passwordHash string, role model.Role) (*model.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || passwordHash == "" || !role.IsValid() {
		return nil, ErrInvalidInput
	}

	now := time.Now().UTC()
	result, err := s.db.Exec(
		`INSERT INTO users (email, password_hash, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		email, passwordHash, role, now, now,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return nil, ErrConflict
		}
		return nil, fmt.Errorf("insert user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("last insert id: %w", err)
	}

	return s.GetUserByID(id)
}

func (s *Store) GetUserByID(id int64) (*model.User, error) {
	row := s.db.QueryRow(`SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE id = ?`, id)
	return scanUser(row)
}

func (s *Store) GetUserByEmail(email string) (*model.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	row := s.db.QueryRow(`SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE email = ?`, email)
	return scanUser(row)
}

func (s *Store) ListUsers(role *model.Role) ([]model.User, error) {
	query := `SELECT id, email, password_hash, role, created_at, updated_at FROM users`
	args := []any{}

	if role != nil {
		query += ` WHERE role = ?`
		args = append(args, *role)
	}
	query += ` ORDER BY id`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	users := make([]model.User, 0)
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}

	return users, rows.Err()
}

func scanUser(scanner interface {
	Scan(dest ...any) error
}) (*model.User, error) {
	var user model.User
	var createdAt, updatedAt string

	err := scanner.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}

	user.CreatedAt = parseTime(createdAt)
	user.UpdatedAt = parseTime(updatedAt)
	return &user, nil
}
