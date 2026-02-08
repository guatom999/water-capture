package notifiers

// Notifier interface สำหรับ notification channels ต่างๆ
type Notifier interface {
	Send(target string, message NotificationMessage) error
	GetChannel() string
}

// NotificationMessage ข้อมูลที่จะส่ง notification
type NotificationMessage struct {
	LocationID   int     `json:"location_id"`
	LocationName string  `json:"location_name"`
	WaterLevel   float64 `json:"water_level"`
	ShoreLevel   float64 `json:"shore_level"`
	Status       string  `json:"status"`
	MeasuredAt   string  `json:"measured_at"`
}

// NotifierRegistry เก็บ notifiers ทั้งหมด
type NotifierRegistry struct {
	notifiers map[string]Notifier
}

// NewNotifierRegistry สร้าง registry ใหม่
func NewNotifierRegistry() *NotifierRegistry {
	return &NotifierRegistry{
		notifiers: make(map[string]Notifier),
	}
}

// Register ลงทะเบียน notifier
func (r *NotifierRegistry) Register(notifier Notifier) {
	r.notifiers[notifier.GetChannel()] = notifier
}

// Get ดึง notifier ตาม channel
func (r *NotifierRegistry) Get(channel string) (Notifier, bool) {
	n, ok := r.notifiers[channel]
	return n, ok
}

// GetAll ดึง notifiers ทั้งหมด
func (r *NotifierRegistry) GetAll() map[string]Notifier {
	return r.notifiers
}
