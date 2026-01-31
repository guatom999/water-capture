package main

import (
	"log"
	"os"

	"github.com/guatom999/self-boardcast/internal/config"
	"github.com/guatom999/self-boardcast/internal/notifiers"
	"github.com/guatom999/self-boardcast/internal/tasks"
	"github.com/hibiken/asynq"
)

func main() {
	// Load config
	config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			return "../../.env"
		}
		return os.Args[1]
	}())

	// Initialize notifiers
	notifierRegistry := notifiers.NewNotifierRegistry()

	// Register Webhook notifier with default URL from ENV
	// webhookURL := os.Getenv("WEBHOOK_URL")
	// if webhookURL == "" {
	// 	log.Fatal("[WORKER] WEBHOOK_URL is required")
	// }

	// Init Line Notifier
	lineToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	if lineToken == "" {
		log.Fatal("[WORKER] LINE_CHANNEL_ACCESS_TOKEN is required")
	}

	// notifierRegistry.Register(notifiers.NewWebhookNotifier(webhookURL))
	notifierRegistry.Register(notifiers.NewLineNotifier(lineToken))

	notiHandler := tasks.NewNotificationHandler(notifierRegistry)

	// Redis address
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"notifications": 10,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeWaterAlert, notiHandler.HandleWaterAlert)

	log.Println("[WORKER] Starting worker server (Broadcast Mode)...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
