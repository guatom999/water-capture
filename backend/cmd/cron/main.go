package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/guatom999/self-boardcast/internal/config"
	"github.com/guatom999/self-boardcast/internal/database"
	"github.com/guatom999/self-boardcast/internal/jobs"
	"github.com/guatom999/self-boardcast/internal/repositories"
	"github.com/guatom999/self-boardcast/internal/services"
)

func main() {

	cfg := config.LoadConfig("../../.env")

	db := database.DatabaseConnect(cfg)

	repo := repositories.NewWaterLevelRepository(db)
	service := services.NewWaterLevelService(repo, cfg.App.BaseURL, cfg)

	log.Println("Starting cron job scheduler...")
	jobs.NewWaterJob(service).ScheduleGetWaterLevel(context.Background())
	log.Println("Cron job started successfully. Press Ctrl+C to stop.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down cron job...")
}
