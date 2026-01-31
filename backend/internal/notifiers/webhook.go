package notifiers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// WebhookNotifier ส่ง notification ผ่าน Webhook (N8N, Zapier, etc.)
type WebhookNotifier struct {
	defaultURL string
}

// NewWebhookNotifier สร้าง Webhook notifier ใหม่
func NewWebhookNotifier(defaultURL string) *WebhookNotifier {
	return &WebhookNotifier{
		defaultURL: defaultURL,
	}
}

// GetChannel return channel name
func (n *WebhookNotifier) GetChannel() string {
	return "webhook"
}

// WebhookPayload ข้อมูลที่ส่งไป webhook
type WebhookPayload struct {
	Event   string              `json:"event"`
	Message NotificationMessage `json:"message"`
}

// Send ส่ง webhook notification
func (n *WebhookNotifier) Send(target string, message NotificationMessage) error {
	url := target
	if url == "" {
		url = n.defaultURL
	}
	if url == "" {
		return fmt.Errorf("webhook URL is required")
	}

	payload := WebhookPayload{
		Event:   "water_alert",
		Message: message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status: %d", resp.StatusCode)
	}

	return nil
}
