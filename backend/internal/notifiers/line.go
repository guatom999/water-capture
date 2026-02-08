package notifiers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const lineNotifyURL = "https://api.line.me/v2/bot/message/broadcast"

type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Request struct {
	Messages []Message `json:"messages"`
}

// LineNotifier ส่ง notification ผ่าน LINE Notify
type LineNotifier struct {
	defaultToken string
}

// NewLineNotifier สร้าง LINE notifier ใหม่
func NewLineNotifier(defaultToken string) *LineNotifier {
	return &LineNotifier{
		defaultToken: defaultToken,
	}
}

// GetChannel return channel name
func (n *LineNotifier) GetChannel() string {
	return "line"
}

// Send ส่ง LINE notification
func (n *LineNotifier) Send(target string, message NotificationMessage) error {
	token := target
	if token == "" {
		token = n.defaultToken
	}

	msg := Request{
		Messages: []Message{
			{
				Type: "text",
				Text: fmt.Sprintf(
					"🌊 แจ้งเตือนระดับน้ำ\n"+
						"📍 สถานี: %s\n"+
						"💧 ระดับน้ำ: %.2f cm\n"+
						"🏔️ ระดับตลิ่ง: %.2f cm\n"+
						"⚠️ สถานะ: %s\n"+
						"🕐 เวลา: %s",
					message.LocationName,
					message.WaterLevel,
					message.ShoreLevel,
					message.Status,
					message.MeasuredAt,
				),
				// Text: "สวัสดี",
			},
		},
	}

	data, _ := json.Marshal(msg)

	req, err := http.NewRequest("POST", lineNotifyURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send LINE notify: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("LINE notify returned status: %d", resp.StatusCode)
	}

	return nil
}

// SendWithImage ส่ง LINE notification พร้อมรูปภาพ
func (n *LineNotifier) SendWithImage(target string, message NotificationMessage, imageURL string) error {
	token := target
	if token == "" {
		token = n.defaultToken
	}

	msg := fmt.Sprintf(
		"🌊 แจ้งเตือนระดับน้ำ\n"+
			"📍 สถานี: %s\n"+
			"💧 ระดับน้ำ: %.2f cm\n"+
			"⚠️ สถานะ: %s",
		message.LocationName,
		message.WaterLevel,
		message.Status,
	)

	data := url.Values{}
	data.Set("message", msg)
	if imageURL != "" {
		data.Set("imageThumbnail", imageURL)
		data.Set("imageFullsize", imageURL)
	}

	req, err := http.NewRequest("POST", lineNotifyURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("LINE notify returned status: %d", resp.StatusCode)
	}

	return nil
}
