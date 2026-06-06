package notification

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"

	"github.com/make-smart-products/requests-api/internal/config"
	"github.com/make-smart-products/requests-api/internal/model"
	"github.com/make-smart-products/requests-api/internal/repository"
)

type Sender struct {
	cfg   config.Config
	store *repository.Store
}

func NewSender(cfg config.Config, store *repository.Store) *Sender {
	return &Sender{cfg: cfg, store: store}
}

func (s *Sender) Dispatch(userID int64, channel model.NotificationChannel, notifType, title, body string) error {
	now := time.Now().UTC()
	item := &model.Notification{
		UserID:  userID,
		Channel: channel,
		Type:    notifType,
		Title:   title,
		Body:    body,
		IsRead:  channel == model.ChannelInApp,
		SentAt:  &now,
	}

	switch channel {
	case model.ChannelEmail:
		if err := s.sendEmail(userID, title, body); err != nil {
			return err
		}
	case model.ChannelSMS:
		if err := s.sendSMS(userID, body); err != nil {
			return err
		}
	case model.ChannelInApp:
		item.SentAt = nil
	default:
		return repository.ErrInvalidInput
	}

	_, err := s.store.CreateNotification(item)
	return err
}

func (s *Sender) NotifyApplicationCreated(app *model.Application) {
	title := "Заявка создана"
	body := fmt.Sprintf("Ваша заявка #%d «%s» сохранена со статусом %s.", app.ID, app.Title, app.Status)
	_ = s.Dispatch(app.UserID, model.ChannelInApp, "application_created", title, body)
	_ = s.Dispatch(app.UserID, model.ChannelEmail, "application_created", title, body)

	managers, err := s.store.ListUsers(ptrRole(model.RoleManager))
	if err != nil {
		return
	}
	for _, manager := range managers {
		managerBody := fmt.Sprintf("Новая заявка #%d от клиента %d: «%s».", app.ID, app.UserID, app.Title)
		_ = s.Dispatch(manager.ID, model.ChannelInApp, "application_created_manager", "Новая заявка", managerBody)
		_ = s.Dispatch(manager.ID, model.ChannelEmail, "application_created_manager", "Новая заявка", managerBody)
	}
}

func (s *Sender) NotifyApplicationStatusChanged(app *model.Application, previous model.ApplicationStatus) {
	title := "Статус заявки изменён"
	body := fmt.Sprintf("Заявка #%d «%s»: %s → %s.", app.ID, app.Title, previous, app.Status)
	_ = s.Dispatch(app.UserID, model.ChannelInApp, "application_status_changed", title, body)
	_ = s.Dispatch(app.UserID, model.ChannelEmail, "application_status_changed", title, body)
	_ = s.Dispatch(app.UserID, model.ChannelSMS, "application_status_changed", title, body)
}

func (s *Sender) sendEmail(userID int64, subject, body string) error {
	user, err := s.store.GetUserByID(userID)
	if err != nil {
		return err
	}

	profile, err := s.store.GetProfile(userID)
	if err != nil && err != repository.ErrNotFound {
		return err
	}

	to := user.Email
	if profile != nil && strings.TrimSpace(profile.Phone) != "" {
		_ = profile
	}

	message := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s", to, subject, body))
	addr := fmt.Sprintf("%s:%d", s.cfg.SMTP.Host, s.cfg.SMTP.Port)

	if s.cfg.SMTP.Host == "" || s.cfg.Env == "development" && s.cfg.SMTP.Host == "localhost" {
		log.Printf("[email] to=%s subject=%q body=%q", to, subject, body)
		return nil
	}

	var auth smtp.Auth
	if s.cfg.SMTP.User != "" {
		auth = smtp.PlainAuth("", s.cfg.SMTP.User, s.cfg.SMTP.Password, s.cfg.SMTP.Host)
	}

	return smtp.SendMail(addr, auth, s.cfg.SMTP.From, []string{to}, message)
}

func (s *Sender) sendSMS(userID int64, body string) error {
	profile, err := s.store.GetProfile(userID)
	if err != nil {
		return err
	}
	if strings.TrimSpace(profile.Phone) == "" {
		return repository.ErrInvalidInput
	}

	switch s.cfg.SMS.Provider {
	case "log", "":
		log.Printf("[sms] to=%s body=%q", profile.Phone, body)
		return nil
	default:
		log.Printf("[sms:%s] to=%s body=%q", s.cfg.SMS.Provider, profile.Phone, body)
		return nil
	}
}

func ptrRole(role model.Role) *model.Role {
	return &role
}
