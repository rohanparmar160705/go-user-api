/*
Package handler manages HTTP requests and responses.
The UserHandler struct methods correspond to API endpoints.
Responsibilities:
- Parsing request body and query parameters.
- Validating input using `go-playground/validator`.
- Calling the Service layer for business logic.
- Formatting and sending JSON responses with appropriate HTTP status codes.
*/
package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
	"github.com/rohanparmar/go-user-api/internal/logger"
	"github.com/rohanparmar/go-user-api/internal/models"
	"github.com/rohanparmar/go-user-api/internal/service"
	"go.uber.org/zap"
)


type UserHandler struct {
	service  service.UserService
	validate *validator.Validate
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if err := h.validate.Struct(req); err != nil {
		logger.Log.Error("Validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	// Create user
	user, err := h.service.CreateUser(c.Context(), req.Name, req.DOB)
	if err != nil {
		logger.Log.Error("Failed to create user", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	logger.Log.Info("User created successfully", zap.Int32("user_id", user.ID))

	// Return response without age (as per task requirement)
	response := models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Time.Format("2006-01-02"),
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	// Get ID from params
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log.Error("Invalid user ID", zap.String("id", idStr))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Get user from service
	user, err := h.service.GetUserByID(c.Context(), int32(id))
	if err != nil {
		logger.Log.Error("User not found", zap.Int("id", id), zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Calculate age
	age := h.service.CalculateAge(user.Dob.Time)

	logger.Log.Info("User retrieved successfully", zap.Int32("user_id", user.ID))

	// Return response with age
	response := models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Time.Format("2006-01-02"),
		Age:  &age,
	}

	return c.JSON(response)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	// Parse page and limit from query params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	// Get paginated users
	response, err := h.service.ListUsers(c.Context(), page, limit)
	if err != nil {
		logger.Log.Error("Failed to list users", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}

	logger.Log.Info("Users listed successfully", 
		zap.Int("page", page),
		zap.Int("limit", limit),
		zap.Int64("total", response.Total),
	)

	return c.JSON(response)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	// Get ID from params
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log.Error("Invalid user ID", zap.String("id", idStr))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req models.UpdateUserRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if err := h.validate.Struct(req); err != nil {
		logger.Log.Error("Validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	// Update user
	user, err := h.service.UpdateUser(c.Context(), int32(id), req.Name, req.DOB)
	if err != nil {
		logger.Log.Error("Failed to update user", zap.Int("id", id), zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	logger.Log.Info("User updated successfully", zap.Int32("user_id", user.ID))

	// Return response without age (as per task requirement)
	response := models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Time.Format("2006-01-02"),
	}

	return c.JSON(response)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	// Get ID from params
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log.Error("Invalid user ID", zap.String("id", idStr))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Delete user
	if err := h.service.DeleteUser(c.Context(), int32(id)); err != nil {
		logger.Log.Error("Failed to delete user", zap.Int("id", id), zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	logger.Log.Info("User deleted successfully", zap.Int("id", id))

	return c.SendStatus(fiber.StatusNoContent)
}

