/*
Package service contains the business logic of the application.
It orchestrates the flow of data between the handler and repository layers.
Key responsibilities include:
- Dynamic Age Calculation logic.
- Pagination calculations.
- Additional business validation rules.
*/
package service

import (
	"context"
	"errors"
	"time"

	"github.com/rohanparmar/go-user-api/internal/repository"
	"github.com/rohanparmar/go-user-api/internal/models"
	db "github.com/rohanparmar/go-user-api/db/sqlc/generated"
)

type UserService interface {
	CreateUser(ctx context.Context, name string, dob string) (db.User, error)
	GetUserByID(ctx context.Context, id int32) (db.User, error)
	ListUsers(ctx context.Context, page, limit int) (models.UsersListResponse, error)
	UpdateUser(ctx context.Context, id int32, name string, dob string) (db.User, error)
	DeleteUser(ctx context.Context, id int32) error
	CalculateAge(dob time.Time) int
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, name string, dob string) (db.User, error) {
	if name == "" {
		return db.User{}, errors.New("name cannot be empty")
	}
	if _, err := time.Parse("2006-01-02", dob); err != nil {
		return db.User{}, errors.New("invalid date format, use YYYY-MM-DD")
	}
	return s.repo.Create(ctx, name, dob)
}

func (s *userService) GetUserByID(ctx context.Context, id int32) (db.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context, page, limit int) (models.UsersListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	
	// Get total count
	total, err := s.repo.Count(ctx)
	if err != nil {
		return models.UsersListResponse{}, err
	}
	
	// Get paginated users
	users, err := s.repo.List(ctx, int32(limit), int32(offset))
	if err != nil {
		return models.UsersListResponse{}, err
	}
	
	// Calculate total pages
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	// Map to response
	var responseData []models.UserResponse
	for _, user := range users {
		age := s.CalculateAge(user.Dob.Time)
		responseData = append(responseData, models.UserResponse{
			ID:   user.ID,
			Name: user.Name,
			DOB:  user.Dob.Time.Format("2006-01-02"),
			Age:  &age,
		})
	}
	
	return models.UsersListResponse{
		Data:       responseData,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int32, name string, dob string) (db.User, error) {
	if name == "" {
		return db.User{}, errors.New("name cannot be empty")
	}
	if _, err := time.Parse("2006-01-02", dob); err != nil {
		return db.User{}, errors.New("invalid date format, use YYYY-MM-DD")
	}
	return s.repo.Update(ctx, id, name, dob)
}

func (s *userService) DeleteUser(ctx context.Context, id int32) error {
	return s.repo.Delete(ctx, id)
}

// CalculateAge calculates the age from date of birth
func (s *userService) CalculateAge(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()
	
	// Adjust if birthday hasn't occurred this year
	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		age--
	}
	
	return age
}
