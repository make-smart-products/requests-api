package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/make-smart-products/requests-api/internal/auth"
	"github.com/make-smart-products/requests-api/internal/config"
	"github.com/make-smart-products/requests-api/internal/handler"
	"github.com/make-smart-products/requests-api/internal/middleware"
	"github.com/make-smart-products/requests-api/internal/model"
	"github.com/make-smart-products/requests-api/internal/notification"
	"github.com/make-smart-products/requests-api/internal/repository"
	"github.com/make-smart-products/requests-api/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := repository.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	store := repository.NewStore(db)
	if err := service.SeedAdmin(store, "admin@requests.local", "admin12345"); err != nil {
		log.Fatalf("seed admin: %v", err)
	}

	tokens := auth.NewTokenManager(cfg.JWTSecret, service.TokenTTL())
	notify := notification.NewSender(cfg, store)
	svc := service.New(store, tokens, notify)
	api := handler.NewAPI(svc)

	router := chi.NewRouter()
	router.Use(chimw.RequestID)
	router.Use(chimw.RealIP)
	router.Use(chimw.Logger)
	router.Use(chimw.Recoverer)
	router.Use(chimw.Timeout(60 * time.Second))

	router.Get("/health", api.Health)

	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", api.Register)
		r.Post("/auth/login", api.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticate(tokens))

			r.Get("/profile", api.GetProfile)
			r.Put("/profile", api.UpdateProfile)

			r.Route("/users", func(r chi.Router) {
				r.Use(middleware.RequireRoles(model.RoleAdmin, model.RoleManager))
				r.Get("/", api.ListUsers)
			})

			r.Route("/applications", func(r chi.Router) {
				r.Get("/", api.ListApplications)
				r.Post("/", api.CreateApplication)
				r.Get("/{id}", api.GetApplication)
				r.Patch("/{id}", api.UpdateApplication)
				r.Delete("/{id}", api.DeleteApplication)
			})

			r.Route("/notifications", func(r chi.Router) {
				r.Get("/", api.ListNotifications)
				r.Patch("/{id}/read", api.MarkNotificationRead)
			})
		})
	})

	server := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server started on %s", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
}
