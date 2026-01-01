package jobs

import (
	"context"

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

	c.cron.AddFunc("0 */20 * * * *", func() {
		c.service.ScheduleGetWaterLevel(ctx)
	})

	c.cron.Start()
}
