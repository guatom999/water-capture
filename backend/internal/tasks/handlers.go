package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/guatom999/self-boardcast/internal/notifiers"
	"github.com/hibiken/asynq"
)

// NotificationHandler จัดการ notification tasks
type NotificationHandler struct {
	notifiers *notifiers.NotifierRegistry
}

// NewNotificationHandler สร้าง handler ใหม่
func NewNotificationHandler(
	notifierRegistry *notifiers.NotifierRegistry,
) *NotificationHandler {
	return &NotificationHandler{
		notifiers: notifierRegistry,
	}
}

// HandleWaterAlert จัดการ water alert task (Broadcast Mode)
func (h *NotificationHandler) HandleWaterAlert(ctx context.Context, t *asynq.Task) error {
	var payload WaterAlertPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("[WORKER] Broadcasting alert for %s (LocationID: %d, Status: %s)",
		payload.LocationName, payload.LocationID, payload.Status)

	// สร้าง notification message
	message := notifiers.NotificationMessage{
		LocationID:   payload.LocationID,
		LocationName: payload.LocationName,
		WaterLevel:   payload.WaterLevel,
		ShoreLevel:   payload.ShoreLevel,
		Status:       payload.Status,
		MeasuredAt:   payload.MeasuredAt,
	}

	// Broadcast ไปที่ webhook โดยตรง (ไม่ต้องเช็ค subscription)
	// notifier, ok := h.notifiers.Get("webhook")
	// if !ok {
	// 	return fmt.Errorf("webhook notifier not registered")
	// }

	notifier, ok := h.notifiers.Get("line")
	if !ok {
		return fmt.Errorf("line notifier not registered")
	}

	// ส่ง notification (target = "" จะใช้ default URL จาก ENV)
	if err := notifier.Send("", message); err != nil {
		log.Printf("[WORKER] Failed to broadcast: %v", err)
		return fmt.Errorf("failed to broadcast: %w", err)
	}

	log.Printf("[WORKER] Broadcast successful for %s", payload.LocationName)
	return nil
}
