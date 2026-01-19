package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/guatom999/self-boardcast/internal/utils"
	"github.com/hibiken/asynq"
)

func HandleWaterAlert(ctx context.Context, t *asynq.Task) error {
	var payload WaterAlertPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("[WORKER] Processing alert for %s (LocationID: %d)", payload.LocationName, payload.LocationID)

	if err := utils.HttpPostJSON("http://badzboss-n8n.duckdns.org:5678/webhook-test/da1f7e4e-9927-4b87-b2bb-8295604937b8", payload); err != nil {
		return fmt.Errorf("failed to post JSON: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("[WORKER] Notification sent successfully for %s", payload.LocationName)
	return nil
}
