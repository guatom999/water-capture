package jobs

import (
	"context"
	"log"
	"time"

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

	c.cron.AddFunc("0 */10 * * * *", func() {
		time.Sleep(time.Second * 10)
		waterLevels, err := c.service.ScheduleGetWaterLevel(ctx)
		if err != nil {
			log.Println("failed to schedule get water level", err)
			return
		}

		// fileName = waterLevel.Image
		// locationID = int(waterLevel.StationID)

		for _, waterLevel := range waterLevels {
			if waterLevel.Danger == "DANGER" || waterLevel.Danger == "WATCH" {

				payload := tasks.WaterAlertPayload{
					LocationID:   int(waterLevel.StationID),
					LocationName: waterLevel.Source.String,
					// ShoreLevel:   1.29,
					Status:      waterLevel.Danger,
					WaterLevel:  waterLevel.LevelCm,
					Description: waterLevel.Note,
					MeasuredAt:  utils.ParseTimeToString(waterLevel.MeasuredAt),
				}

				if err := c.producer.EnqueueWaterAlert(payload); err != nil {
					log.Printf("[CRON] Failed to enqueue alert: %v", err)
				}
			}
		}

	})

	_ = fileName
	_ = locationID

	// c.cron.AddFunc("0 0 0 * * *", func() {
	// 	if err := c.service.ScheduleDeleteWaterLevel(ctx, fileName, locationID); err != nil {
	// 		log.Println("failed to schedule delete data of water level", err)
	// 	}
	// })

	c.cron.Start()
}
