package jobs

import (
	"context"
	"log"
	"time"

	"github.com/guatom999/self-boardcast/internal/services"
	"github.com/robfig/cron"
)

type WaterJob struct {
	cron    *cron.Cron
	service services.WaterLevelServiceInterface
}

type WaterJobInterface interface {
	ScheduleGetWaterLevel(ctx context.Context)
}

func NewWaterJob(service services.WaterLevelServiceInterface) *WaterJob {
	return &WaterJob{
		cron:    cron.New(),
		service: service,
	}
}

func (c *WaterJob) ScheduleGetWaterLevel(ctx context.Context) {

	var fileName string
	var locationID int

	c.cron.AddFunc("0 0 */1 * * *", func() {
		time.Sleep(time.Second * 20)
		jobsFileName, jobLocationID, err := c.service.ScheduleGetWaterLevel(ctx)
		fileName, locationID = jobsFileName, jobLocationID
		if err != nil {
			log.Println("failed to schedule get water level", err)
		}
	})

	c.cron.AddFunc("0 0 0 * * *", func() {
		if err := c.service.ScheduleDeleteWaterLevel(ctx, fileName, locationID); err != nil {
			log.Println("failed to schedule delete data of water level", err)
		}
	})

	c.cron.Start()
}
