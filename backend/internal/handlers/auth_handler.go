package handlers

import (
	"net/http"
	"strings"

	"github.com/guatom999/self-boardcast/internal/models"
	"github.com/guatom999/self-boardcast/internal/services"
	"github.com/labstack/echo/v4"
)

type AuthHandlerInterface interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
	RefreshToken(c echo.Context) error
	Logout(c echo.Context) error
	GetMe(c echo.Context) error
}

type authHandler struct {
	authService services.AuthServiceInterface
}

func NewAuthHandler(authService services.AuthServiceInterface) AuthHandlerInterface {
	return &authHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Register request"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /auth/register [post]
func (h *authHandler) Register(c echo.Context) error {
	req := new(models.RegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Basic validation
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email, password and name are required"})
	}

	if len(req.Password) < 8 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password must be at least 8 characters"})
	}

	response, err := h.authService.Register(c.Request().Context(), req)
	if err != nil {
		if err == services.ErrUserAlreadyExists {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to register user"})
	}

	return c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *authHandler) Login(c echo.Context) error {
	req := new(models.LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email and password are required"})
	}

	response, err := h.authService.Login(c.Request().Context(), req)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to login"})
	}

	return c.JSON(http.StatusOK, response)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func (h *authHandler) RefreshToken(c echo.Context) error {
	req := new(models.RefreshTokenRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Refresh token is required"})
	}

	response, err := h.authService.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		if err == services.ErrInvalidToken {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to refresh token"})
	}

	return c.JSON(http.StatusOK, response)
}

// Logout godoc
// @Summary Logout user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LogoutRequest true "Logout request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/logout [post]
func (h *authHandler) Logout(c echo.Context) error {
	req := new(models.LogoutRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Refresh token is required"})
	}

	if err := h.authService.Logout(c.Request().Context(), req.RefreshToken); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to logout"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// GetMe godoc
// @Summary Get current user info
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} map[string]string
// @Router /auth/me [get]
func (h *authHandler) GetMe(c echo.Context) error {
	// Get claims from context (set by auth middleware)
	claims, ok := c.Get("user").(*models.TokenClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	return c.JSON(http.StatusOK, models.UserResponse{
		ID:    claims.UserID,
		Email: claims.Email,
		Role:  claims.Role,
	})
}

// Helper function to extract token from Authorization header
func ExtractTokenFromHeader(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}
