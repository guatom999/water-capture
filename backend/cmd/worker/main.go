package main

import (
	"log"

	"github.com/guatom999/self-boardcast/internal/tasks"
	"github.com/hibiken/asynq"
)

func main() {
	redisAddr := "localhost:6379"

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
	mux.HandleFunc(tasks.TypeWaterAlert, tasks.HandleWaterAlert)

	log.Println("[WORKER] Starting worker server...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
