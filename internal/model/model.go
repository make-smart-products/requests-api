package model

import "time"

type Role string

const (
	RoleClient  Role = "client"
	RoleManager Role = "manager"
	RoleAdmin   Role = "admin"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleClient, RoleManager, RoleAdmin:
		return true
	default:
		return false
	}
}

type ApplicationStatus string

const (
	StatusDraft      ApplicationStatus = "draft"
	StatusSubmitted  ApplicationStatus = "submitted"
	StatusInReview   ApplicationStatus = "in_review"
	StatusApproved   ApplicationStatus = "approved"
	StatusRejected   ApplicationStatus = "rejected"
)

func (s ApplicationStatus) IsValid() bool {
	switch s {
	case StatusDraft, StatusSubmitted, StatusInReview, StatusApproved, StatusRejected:
		return true
	default:
		return false
	}
}

type NotificationChannel string

const (
	ChannelEmail NotificationChannel = "email"
	ChannelSMS   NotificationChannel = "sms"
	ChannelInApp NotificationChannel = "in_app"
)

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Profile struct {
	UserID    int64     `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	Company   string    `json:"company"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Application struct {
	ID                int64             `json:"id"`
	UserID            int64             `json:"user_id"`
	AssignedManagerID *int64            `json:"assigned_manager_id,omitempty"`
	Title             string            `json:"title"`
	Description       string            `json:"description"`
	Status            ApplicationStatus `json:"status"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

type Notification struct {
	ID        int64               `json:"id"`
	UserID    int64               `json:"user_id"`
	Channel   NotificationChannel `json:"channel"`
	Type      string              `json:"type"`
	Title     string              `json:"title"`
	Body      string              `json:"body"`
	IsRead    bool                `json:"is_read"`
	SentAt    *time.Time          `json:"sent_at,omitempty"`
	CreatedAt time.Time           `json:"created_at"`
}
