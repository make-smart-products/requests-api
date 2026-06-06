package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/make-smart-products/requests-api/internal/middleware"
	"github.com/make-smart-products/requests-api/internal/model"
	"github.com/make-smart-products/requests-api/internal/repository"
	"github.com/make-smart-products/requests-api/internal/service"
)

type API struct {
	svc *service.Service
}

func NewAPI(svc *service.Service) *API {
	return &API{svc: svc}
}

func (a *API) Register(w http.ResponseWriter, r *http.Request) {
	var input service.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	resp, err := a.svc.Register(input)
	if err != nil {
		a.handleError(w, err)
		return
	}
	resp.User.PasswordHash = ""
	WriteJSON(w, http.StatusCreated, resp)
}

func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	var input service.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	resp, err := a.svc.Login(input)
	if err != nil {
		a.handleError(w, err)
		return
	}
	resp.User.PasswordHash = ""
	WriteJSON(w, http.StatusOK, resp)
}

func (a *API) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	profile, err := a.svc.GetProfile(claims, claims.UserID)
	if err != nil {
		a.handleError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, profile)
}

func (a *API) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	var profile model.Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	updated, err := a.svc.UpdateProfile(claims, claims.UserID, profile)
	if err != nil {
		a.handleError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, updated)
}

func (a *API) ListUsers(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	var role *model.Role
	if value := r.URL.Query().Get("role"); value != "" {
		parsed := model.Role(value)
		if !parsed.IsValid() {
			WriteError(w, http.StatusBadRequest, "invalid role")
			return
		}
		role = &parsed
	}

	users, err := a.svc.ListUsers(claims, role)
	if err != nil {
		a.handleError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, users)
}

func (a *API) ListApplications(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	var status *model.ApplicationStatus
	if value := r.URL.Query().Get("status"); value != "" {
		parsed, err := repository.NormalizeStatus(value)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "invalid status")
			return
		}
		status = &parsed
	}

	apps, err := a.svc.ListApplications(claims, status)
	if err != nil {
		a.handleError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, apps)
}

func (a *API) CreateApplication(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	var input service.ApplicationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	app, err := a.svc.CreateApplication(claims, input)
	if err != nil {
		a.handleError(w, err)
		return
	}
	WriteJSON(w, http.StatusCreated, app)
}

func (a *API) GetApplication(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	app, err := a.svc.GetApplication(claims, id)
	if err != nil {
		a.handleError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, app)
}

func (a *API) UpdateApplication(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input service.ApplicationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	app, err := a.svc.UpdateApplication(claims, id, input)
	if err != nil {
		a.handleError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, app)
}

func (a *API) DeleteApplication(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := a.svc.DeleteApplication(claims, id); err != nil {
		a.handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) ListNotifications(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	unreadOnly := r.URL.Query().Get("unread") == "true"

	items, err := a.svc.ListNotifications(claims, unreadOnly)
	if err != nil {
		a.handleError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, items)
}

func (a *API) MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	item, err := a.svc.MarkNotificationRead(claims, id)
	if err != nil {
		a.handleError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, item)
}

func (a *API) Health(w http.ResponseWriter, _ *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *API) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrUnauthorized):
		WriteError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, service.ErrForbidden):
		WriteError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, service.ErrBadRequest), errors.Is(err, repository.ErrInvalidInput):
		WriteError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, repository.ErrNotFound):
		WriteError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, repository.ErrConflict):
		WriteError(w, http.StatusConflict, err.Error())
	default:
		WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}

func parseID(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
