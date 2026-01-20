package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guatom999/self-boardcast/internal/config"
	"github.com/guatom999/self-boardcast/internal/handlers"
	"github.com/guatom999/self-boardcast/internal/repositories"
	"github.com/guatom999/self-boardcast/internal/services"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	db          *sqlx.DB
	echo        *echo.Echo
	cfg         *config.Config
	authService services.AuthServiceInterface
}

func NewServer(db *sqlx.DB, cfg *config.Config) *Server {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	// e.Use(middleware.SecurityHeaders())

	// Initialize auth service for use across modules
	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo, cfg)

	return &Server{
		db:          db,
		echo:        e,
		cfg:         cfg,
		authService: authService,
	}
}

func (s *Server) WaterModules() {
	repo := repositories.NewWaterLevelRepository(s.db)
	service := services.NewWaterLevelService(repo, s.cfg.App.BaseURL, s.cfg)
	handler := handlers.NewMapHandler(service)

	s.echo.GET("/heath", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})

	s.echo.GET("/markers", handler.GetMapMarkers)
	s.echo.GET("/markers/detail", handler.GetSectionDetail)
}

func (s *Server) ImageModules() {
	imageHandler := handlers.NewImageHandler(s.cfg)

	s.echo.GET("/images/:filename", imageHandler.ServeImage)
	s.echo.GET("/images/health", imageHandler.HealthCheck)
}

func (s *Server) AuthModules() {
	authHandler := handlers.NewAuthHandler(s.authService)

	// Public auth routes
	auth := s.echo.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.RefreshToken)
	auth.POST("/logout", authHandler.Logout)

	// Protected route example - requires valid access token
	// auth.GET("/me", authHandler.GetMe, customMiddleware.JWTMiddleware(s.authService))
}

func (s *Server) Start() error {
	go func() {
		if err := s.echo.Start(fmt.Sprintf(":%d", s.cfg.Server.Port)); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	s.AuthModules()
	s.WaterModules()
	s.ImageModules()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println(" Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.echo.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exited gracefully")
	return nil
}
