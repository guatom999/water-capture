package tasks

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

type NotificationProducer struct {
	client *asynq.Client
}

func NewNotificationProducer(redisAddr string) *NotificationProducer {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	return &NotificationProducer{client: client}
}

func (p *NotificationProducer) Close() error {
	return p.client.Close()
}

func (p *NotificationProducer) EnqueueWaterAlert(payload WaterAlertPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeWaterAlert, data,
		asynq.MaxRetry(3),
		asynq.Queue("notifications"),
		asynq.Timeout(30*time.Second),
	)

	_, err = p.client.Enqueue(task)
	return err
}
