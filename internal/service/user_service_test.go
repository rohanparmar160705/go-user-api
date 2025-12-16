package service

import (
	"context"
	"testing"
	"time"

	"github.com/rohanparmar/go-user-api/internal/repository"
	db "github.com/rohanparmar/go-user-api/db/sqlc/generated"
	"github.com/stretchr/testify/assert"
)

// Mock repository for testing
type mockRepo struct {
	repository.UserRepository
}

func (m *mockRepo) Create(ctx context.Context, name string, dob string) (db.User, error) {
	return db.User{}, nil
}

func TestCalculateAge(t *testing.T) {
	// Initialize service with mock repo (repo not used for age calc)
	userService := NewUserService(&mockRepo{})
	// We need to access the CalculateAge method which is exposed via interface
	// but mostly we want to test the logic directly attached to the struct implementation
	// or via the interface method if exposed.
	
	// Create a concrete instance to access CalculateAge directly since it's defined on the struct pointer
	// and exposed in interface.
	
	now := time.Now()
	currentYear := now.Year()
	
	tests := []struct {
		name     string
		dob      time.Time
		expected int
	}{
		{
			name:     "Birthday passed this year",
			dob:      time.Date(currentYear-20, now.Month()-1, 1, 0, 0, 0, 0, time.UTC),
			expected: 20,
		},
		{
			name:     "Birthday is today",
			dob:      time.Date(currentYear-20, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			expected: 20,
		},
		{
			name:     "Birthday not yet passed this year",
			dob:      time.Date(currentYear-20, now.Month()+1, 1, 0, 0, 0, 0, time.UTC),
			expected: 19,
		},
		{
			name:     "Born today",
			dob:      now,
			expected: 0,
		},
		{
			name:     "Leap year birthday (Feb 29) on non-leap year",
			dob:      time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC),
			// Age depends on current date relative to Feb 28/Mar 1
			// This test might be tricky without fixing "now", but generally:
			// If today is Mar 1 2023, age is 23. If Feb 28, age is 22.
			expected: calculateExpectedAge(time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			age := userService.CalculateAge(tt.dob)
			assert.Equal(t, tt.expected, age)
		})
	}
}

// Helper to double check logic for test Cases
func calculateExpectedAge(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()
	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		age--
	}
	return age
}
