package tasks

const (
	TypeWaterAlert = "notification:water_alert"
)

type WaterAlertPayload struct {
	LocationID   int     `json:"location_id"`
	LocationName string  `json:"location_name"`
	ShoreLevel   float64 `json:"shore_level"`
	WaterLevel   float64 `json:"water_level"`
	Description  string  `json:"description"`
	MeasuredAt   string  `json:"measured_at"`
}
