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

	"github.com/go-playground/validator/v10"
	"github.com/guatom999/self-boardcast/internal/config"
	"github.com/guatom999/self-boardcast/internal/handlers"
	"github.com/guatom999/self-boardcast/internal/repositories"
	"github.com/guatom999/self-boardcast/internal/services"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	customMiddleware "github.com/guatom999/self-boardcast/internal/middleware"
)

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

type Server struct {
	db   *sqlx.DB
	echo *echo.Echo
	cfg  *config.Config
	// authService services.AuthServiceInterface
}

func NewServer(db *sqlx.DB, cfg *config.Config) *Server {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Validator = &customValidator{validator: validator.New()}
	// e.Use(middleware.SecurityHeaders())

	return &Server{
		db:   db,
		echo: e,
		cfg:  cfg,
	}
}

func (s *Server) WaterModules() {
	repo := repositories.NewWaterLevelRepository(s.db)
	service := services.NewWaterLevelService(repo, s.cfg.App.BaseURL, s.cfg)
	handler := handlers.NewMapHandler(service)

	//schedule manual
	scheduleService := services.NewWaterLevelService(repo, s.cfg.App.BaseURL, s.cfg)

	s.echo.GET("/heath", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})

	s.echo.POST("/add_station_location", handler.AddStationLocation)

	s.echo.GET("/markers", handler.GetMapMarkers)
	s.echo.GET("/markers/detail", handler.GetSectionDetail)

	//test
	s.echo.GET("/test_manual_get_schedule", func(e echo.Context) error {

		scheduleService.ScheduleGetWaterLevel(context.Background())

		return nil
	})

}

func (s *Server) ImageModules() {
	imageHandler := handlers.NewImageHandler(s.cfg)

	s.echo.GET("/images/:filename", imageHandler.ServeImage)
	s.echo.GET("/images/health", imageHandler.HealthCheck)
}

func (s *Server) NotificaitonModules() {
	notificationRepo := repositories.NewNotificationRepository(s.db)
	authRepo := repositories.NewAuthRepository(s.db)
	notificationService := services.NewNotificationService(authRepo, notificationRepo)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	noti := s.echo.Group("/notification")

	noti.POST("/subscribe", notificationHandler.CreateSubscription)
}

func (s *Server) AuthModules() {

	authRepo := repositories.NewAuthRepository(s.db)
	notiRepo := repositories.NewNotificationRepository(s.db)
	authService := services.NewAuthService(notiRepo, authRepo, s.cfg)
	authHandler := handlers.NewAuthHandler(authService)

	auth := s.echo.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.RefreshToken)
	auth.POST("/logout", authHandler.Logout)

	// Protected route example - requires valid access token
	auth.GET("/me", authHandler.GetMe, customMiddleware.JWTMiddleware(authService))
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
	s.NotificaitonModules()

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
