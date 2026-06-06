package service

import (
	"errors"
	"strings"
	"time"

	"github.com/make-smart-products/requests-api/internal/auth"
	"github.com/make-smart-products/requests-api/internal/model"
	"github.com/make-smart-products/requests-api/internal/notification"
	"github.com/make-smart-products/requests-api/internal/repository"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrBadRequest   = errors.New("bad request")
)

type Service struct {
	store   *repository.Store
	tokens  *auth.TokenManager
	notify  *notification.Sender
}

func New(store *repository.Store, tokens *auth.TokenManager, notify *notification.Sender) *Service {
	return &Service{store: store, tokens: tokens, notify: notify}
}

type RegisterInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Company   string `json:"company"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

func (s *Service) Register(input RegisterInput) (*AuthResponse, error) {
	if strings.TrimSpace(input.Email) == "" || len(input.Password) < 8 {
		return nil, ErrBadRequest
	}

	hash, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user, err := s.store.CreateUser(input.Email, hash, model.RoleClient)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return nil, ErrBadRequest
		}
		return nil, err
	}

	_, err = s.store.UpsertProfile(&model.Profile{
		UserID:    user.ID,
		FirstName: strings.TrimSpace(input.FirstName),
		LastName:  strings.TrimSpace(input.LastName),
		Phone:     strings.TrimSpace(input.Phone),
		Company:   strings.TrimSpace(input.Company),
	})
	if err != nil {
		return nil, err
	}

	token, err := s.tokens.Generate(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: user}, nil
}

func (s *Service) Login(input LoginInput) (*AuthResponse, error) {
	user, err := s.store.GetUserByEmail(input.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}

	if !auth.CheckPassword(user.PasswordHash, input.Password) {
		return nil, ErrUnauthorized
	}

	token, err := s.tokens.Generate(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: user}, nil
}

func (s *Service) GetProfile(actor *auth.Claims, userID int64) (*model.Profile, error) {
	if err := s.ensureUserAccess(actor, userID); err != nil {
		return nil, err
	}
	return s.store.GetProfile(userID)
}

func (s *Service) UpdateProfile(actor *auth.Claims, userID int64, profile model.Profile) (*model.Profile, error) {
	if err := s.ensureUserAccess(actor, userID); err != nil {
		return nil, err
	}
	if actor.Role == model.RoleClient && actor.UserID != userID {
		return nil, ErrForbidden
	}
	if actor.Role != model.RoleAdmin && actor.UserID != userID {
		return nil, ErrForbidden
	}

	profile.UserID = userID
	return s.store.UpsertProfile(&profile)
}

func (s *Service) ListUsers(actor *auth.Claims, role *model.Role) ([]model.User, error) {
	if actor.Role != model.RoleAdmin && actor.Role != model.RoleManager {
		return nil, ErrForbidden
	}
	users, err := s.store.ListUsers(role)
	if err != nil {
		return nil, err
	}
	for i := range users {
		users[i].PasswordHash = ""
	}
	return users, nil
}

type ApplicationInput struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Status      *model.ApplicationStatus `json:"status,omitempty"`
	ManagerID   *int64                 `json:"assigned_manager_id,omitempty"`
}

func (s *Service) CreateApplication(actor *auth.Claims, input ApplicationInput) (*model.Application, error) {
	if actor.Role != model.RoleClient && actor.Role != model.RoleAdmin {
		return nil, ErrForbidden
	}
	if strings.TrimSpace(input.Title) == "" {
		return nil, ErrBadRequest
	}

	status := model.StatusDraft
	if input.Status != nil {
		status = *input.Status
	}
	if actor.Role == model.RoleClient && status != model.StatusDraft && status != model.StatusSubmitted {
		return nil, ErrForbidden
	}

	app := &model.Application{
		UserID:            actor.UserID,
		AssignedManagerID: input.ManagerID,
		Title:             strings.TrimSpace(input.Title),
		Description:       strings.TrimSpace(input.Description),
		Status:            status,
	}
	if actor.Role == model.RoleAdmin && input.ManagerID == nil {
		app.UserID = actor.UserID
	}

	created, err := s.store.CreateApplication(app)
	if err != nil {
		return nil, err
	}
	s.notify.NotifyApplicationCreated(created)
	return created, nil
}

func (s *Service) GetApplication(actor *auth.Claims, id int64) (*model.Application, error) {
	app, err := s.store.GetApplication(id)
	if err != nil {
		return nil, err
	}
	if err := s.ensureApplicationAccess(actor, app); err != nil {
		return nil, err
	}
	return app, nil
}

func (s *Service) ListApplications(actor *auth.Claims, status *model.ApplicationStatus) ([]model.Application, error) {
	filter := repository.ApplicationFilter{Status: status}
	switch actor.Role {
	case model.RoleClient:
		filter.UserID = &actor.UserID
	case model.RoleManager:
		filter.AssignedManagerID = &actor.UserID
		filter.IncludeUnassigned = true
	case model.RoleAdmin:
	default:
		return nil, ErrForbidden
	}
	return s.store.ListApplications(filter)
}

func (s *Service) UpdateApplication(actor *auth.Claims, id int64, input ApplicationInput) (*model.Application, error) {
	current, err := s.store.GetApplication(id)
	if err != nil {
		return nil, err
	}
	if err := s.ensureApplicationAccess(actor, current); err != nil {
		return nil, err
	}

	updated := *current
	if strings.TrimSpace(input.Title) != "" {
		updated.Title = strings.TrimSpace(input.Title)
	}
	if input.Description != "" || actor.Role != model.RoleClient {
		updated.Description = strings.TrimSpace(input.Description)
	}

	previousStatus := current.Status
	if input.Status != nil {
		if !repository.IsTransitionAllowed(current.Status, *input.Status, actor.Role) {
			return nil, ErrForbidden
		}
		updated.Status = *input.Status
	}

	if actor.Role == model.RoleManager || actor.Role == model.RoleAdmin {
		if input.ManagerID != nil {
			updated.AssignedManagerID = input.ManagerID
		}
		if actor.Role == model.RoleManager && updated.AssignedManagerID == nil {
			updated.AssignedManagerID = &actor.UserID
		}
	}

	result, err := s.store.UpdateApplication(&updated)
	if err != nil {
		return nil, err
	}

	if previousStatus != result.Status {
		s.notify.NotifyApplicationStatusChanged(result, previousStatus)
	}
	return result, nil
}

func (s *Service) DeleteApplication(actor *auth.Claims, id int64) error {
	app, err := s.store.GetApplication(id)
	if err != nil {
		return err
	}
	if actor.Role != model.RoleClient || app.UserID != actor.UserID {
		return ErrForbidden
	}
	if app.Status != model.StatusDraft {
		return ErrForbidden
	}
	return s.store.DeleteApplication(id)
}

func (s *Service) ListNotifications(actor *auth.Claims, unreadOnly bool) ([]model.Notification, error) {
	return s.store.ListNotifications(actor.UserID, unreadOnly)
}

func (s *Service) MarkNotificationRead(actor *auth.Claims, id int64) (*model.Notification, error) {
	return s.store.MarkNotificationRead(id, actor.UserID)
}

func (s *Service) ensureUserAccess(actor *auth.Claims, userID int64) error {
	switch actor.Role {
	case model.RoleAdmin, model.RoleManager:
		return nil
	case model.RoleClient:
		if actor.UserID == userID {
			return nil
		}
	}
	return ErrForbidden
}

func (s *Service) ensureApplicationAccess(actor *auth.Claims, app *model.Application) error {
	switch actor.Role {
	case model.RoleAdmin:
		return nil
	case model.RoleClient:
		if app.UserID == actor.UserID {
			return nil
		}
	case model.RoleManager:
		if app.AssignedManagerID == nil || *app.AssignedManagerID == actor.UserID {
			return nil
		}
	}
	return ErrForbidden
}

func SeedAdmin(store *repository.Store, email, password string) error {
	existing, err := store.GetUserByEmail(email)
	if err == nil {
		existing.PasswordHash = ""
		return nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return err
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	user, err := store.CreateUser(email, hash, model.RoleAdmin)
	if err != nil {
		return err
	}

	_, err = store.UpsertProfile(&model.Profile{
		UserID:    user.ID,
		FirstName: "System",
		LastName:  "Admin",
		Phone:     "+70000000000",
		Company:   "Requests API",
	})
	return err
}

func TokenTTL() time.Duration {
	return 24 * time.Hour
}
