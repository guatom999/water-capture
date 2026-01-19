package jobs

import (
	"context"
	"log"

	"github.com/guatom999/self-boardcast/internal/services"
	"github.com/guatom999/self-boardcast/internal/tasks"
	"github.com/guatom999/self-boardcast/internal/utils"
	"github.com/robfig/cron"
)

type WaterJob struct {
	cron     *cron.Cron
	service  services.WaterLevelServiceInterface
	producer *tasks.NotificationProducer
}

type WaterJobInterface interface {
	ScheduleGetWaterLevel(ctx context.Context)
}

func NewWaterJob(service services.WaterLevelServiceInterface) *WaterJob {
	producer := tasks.NewNotificationProducer("localhost:6379")
	return &WaterJob{
		cron:     cron.New(),
		service:  service,
		producer: producer,
	}
}

func (c *WaterJob) ScheduleGetWaterLevel(ctx context.Context) {

	var fileName string
	var locationID int

	c.cron.AddFunc("0 */1 * * * *", func() {
		waterLevel, err := c.service.ScheduleGetWaterLevel(ctx)
		if err != nil {
			log.Println("failed to schedule get water level", err)
			return
		}

		fileName = waterLevel.Image
		locationID = int(waterLevel.LocationID)

		if waterLevel.Danger == "DANGER" || waterLevel.Danger == "WATCH" || waterLevel.Danger == "SAFE" {

			payload := tasks.WaterAlertPayload{
				LocationID:   int(waterLevel.LocationID),
				LocationName: waterLevel.Source.String,
				ShoreLevel:   1.29,
				WaterLevel:   waterLevel.LevelCm,
				Description:  waterLevel.Note,
				MeasuredAt:   utils.ParseTimeToString(waterLevel.MeasuredAt),
			}

			if err := c.producer.EnqueueWaterAlert(payload); err != nil {
				log.Printf("[CRON] Failed to enqueue alert: %v", err)
			}
		}
	})

	c.cron.AddFunc("0 */20 * * * *", func() {
		if err := c.service.ScheduleDeleteWaterLevel(ctx, fileName, locationID); err != nil {
			log.Println("failed to schedule delete data of water level", err)
		}
	})

	c.cron.Start()
}
