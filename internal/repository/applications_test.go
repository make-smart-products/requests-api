package repository

import (
	"testing"

	"github.com/make-smart-products/requests-api/internal/model"
)

func TestIsTransitionAllowed(t *testing.T) {
	tests := []struct {
		name    string
		current model.ApplicationStatus
		next    model.ApplicationStatus
		role    model.Role
		want    bool
	}{
		{"client draft to submitted", model.StatusDraft, model.StatusSubmitted, model.RoleClient, true},
		{"client draft to approved", model.StatusDraft, model.StatusApproved, model.RoleClient, false},
		{"manager submitted to in_review", model.StatusSubmitted, model.StatusInReview, model.RoleManager, true},
		{"manager in_review to approved", model.StatusInReview, model.StatusApproved, model.RoleManager, true},
		{"admin any valid", model.StatusDraft, model.StatusRejected, model.RoleAdmin, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTransitionAllowed(tt.current, tt.next, tt.role); got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
